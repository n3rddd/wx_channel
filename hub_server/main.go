package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	"wx_channel/hub_server/database"

	"github.com/gorilla/websocket"
)

// --- Models ---

type MessageType string

const (
	MsgTypeHeartbeat MessageType = "heartbeat"
	MsgTypeCommand   MessageType = "command"
	MsgTypeResponse  MessageType = "response"
)

type CloudMessage struct {
	ID        string          `json:"id"`
	Type      MessageType     `json:"type"`
	ClientID  string          `json:"client_id"`
	Payload   json.RawMessage `json:"payload"`
	Timestamp int64           `json:"timestamp"`
}

type HeartbeatPayload struct {
	Hostname string `json:"hostname"`
	Version  string `json:"version"`
	Status   string `json:"status"`
}

type CommandPayload struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

type ResponsePayload struct {
	RequestID string          `json:"request_id"`
	Success   bool            `json:"success"`
	Data      json.RawMessage `json:"data"`
	Error     string          `json:"error"`
}

// --- Hub Server Logic ---

type Client struct {
	ID       string    `json:"id"`
	Hostname string    `json:"hostname"`
	Version  string    `json:"version"`
	LastSeen time.Time `json:"last_seen"`
	Conn     *websocket.Conn
	mu       sync.Mutex

	respChannels map[string]chan ResponsePayload
	respMu       sync.RWMutex
}

type Hub struct {
	clients  map[string]*Client
	mu       sync.RWMutex
	upgrader websocket.Upgrader
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	clientID := r.Header.Get("X-Client-ID")
	if clientID == "" {
		clientID = r.URL.Query().Get("client_id")
	}
	if clientID == "" {
		http.Error(w, "X-Client-ID required", 400)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:           clientID,
		LastSeen:     time.Now(),
		Conn:         conn,
		respChannels: make(map[string]chan ResponsePayload),
	}

	h.mu.Lock()
	if old, ok := h.clients[clientID]; ok {
		old.Conn.Close()
	}
	h.clients[clientID] = client
	h.mu.Unlock()

	log.Printf("Client connected: %s", clientID)
	// DB: Mark as online
	database.UpsertNode(&database.Node{
		ID:       clientID,
		IP:       r.RemoteAddr,
		Status:   "online",
		LastSeen: time.Now(),
	})

	defer h.cleanup(clientID)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var msg CloudMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		h.handleMessage(client, msg)
	}
}

func (h *Hub) cleanup(id string) {
	h.mu.Lock()
	if c, ok := h.clients[id]; ok {
		c.Conn.Close()
		delete(h.clients, id)
	}
	h.mu.Unlock()
	log.Printf("Client disconnected: %s", id)

	// DB: Mark as offline
	database.UpdateNodeStatus(id, "offline")
}

func (h *Hub) handleMessage(c *Client, msg CloudMessage) {
	c.mu.Lock()
	c.LastSeen = time.Now()
	c.mu.Unlock()

	switch msg.Type {
	case MsgTypeHeartbeat:
		var p HeartbeatPayload
		json.Unmarshal(msg.Payload, &p)
		c.mu.Lock()
		c.Hostname = p.Hostname
		c.Version = p.Version
		c.mu.Unlock()

		// DB: Update heartbeat info
		database.UpsertNode(&database.Node{
			ID:       c.ID,
			Hostname: p.Hostname,
			Version:  p.Version,
			Status:   "online",
			LastSeen: time.Now(),
		})

	case MsgTypeResponse:
		var resp ResponsePayload
		if err := json.Unmarshal(msg.Payload, &resp); err == nil {
			c.respMu.RLock()
			ch, ok := c.respChannels[resp.RequestID]
			c.respMu.RUnlock()
			if ok {
				ch <- resp
			}
		}
	}
}

func (h *Hub) Call(clientID string, action string, data interface{}, timeout time.Duration) (ResponsePayload, error) {
	h.mu.RLock()
	c, ok := h.clients[clientID]
	h.mu.RUnlock()

	if !ok {
		return ResponsePayload{}, fmt.Errorf("client offline")
	}

	reqID := fmt.Sprintf("hub-%d", time.Now().UnixNano())
	payloadData, _ := json.Marshal(data)
	cmd := CommandPayload{Action: action, Data: payloadData}
	cmdData, _ := json.Marshal(cmd)

	// DB: Create Task
	task := &database.Task{
		Type:    action,
		NodeID:  clientID,
		Payload: string(payloadData),
		Status:  "pending",
	}
	database.CreateTask(task)

	msg := CloudMessage{
		ID:        reqID,
		Type:      MsgTypeCommand,
		ClientID:  "hub-server",
		Payload:   cmdData,
		Timestamp: time.Now().Unix(),
	}

	respChan := make(chan ResponsePayload, 1)
	c.respMu.Lock()
	c.respChannels[reqID] = respChan
	c.respMu.Unlock()

	defer func() {
		c.respMu.Lock()
		delete(c.respChannels, reqID)
		c.respMu.Unlock()
	}()

	msgData, _ := json.Marshal(msg)
	c.mu.Lock()
	err := c.Conn.WriteMessage(websocket.TextMessage, msgData)
	c.mu.Unlock()

	if err != nil {
		database.UpdateTaskResult(task.ID, "failed", "", err.Error())
		return ResponsePayload{}, err
	}

	select {
	case resp := <-respChan:
		resBytes, _ := json.Marshal(resp.Data)
		status := "success"
		if !resp.Success {
			status = "failed"
		}
		database.UpdateTaskResult(task.ID, status, string(resBytes), resp.Error)
		return resp, nil
	case <-time.After(timeout):
		database.UpdateTaskResult(task.ID, "timeout", "", "request timeout")
		return ResponsePayload{}, fmt.Errorf("timeout")
	}
}

func main() {
	// 初始化数据库
	if err := database.InitDB("hub_server.db"); err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}

	hub := NewHub()

	// Middleware: Panic Recovery
	withRecovery := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Printf("PANIC: %v\nStack: %s", err, string(debug.Stack()))
					http.Error(w, "Internal Server Error", 500)
				}
			}()
			next(w, r)
		}
	}

	// WebSocket 接入点
	http.HandleFunc("/ws/client", withRecovery(hub.ServeWs))

	// API: 获取节点列表 (读数据库)
	http.HandleFunc("/api/clients", withRecovery(func(w http.ResponseWriter, r *http.Request) {
		nodes, err := database.GetNodes()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(nodes)
	}))

	// API: 获取任务列表
	http.HandleFunc("/api/tasks", withRecovery(func(w http.ResponseWriter, r *http.Request) {
		offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		nodeID := r.URL.Query().Get("node_id")

		if limit <= 0 {
			limit = 20
		}

		tasks, count, err := database.GetTasks(nodeID, offset, limit)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total": count,
			"list":  tasks,
		})
	}))

	http.HandleFunc("/api/call", withRecovery(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ClientID string          `json:"client_id"`
			Action   string          `json:"action"`
			Data     json.RawMessage `json:"data"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		resp, err := hub.Call(req.ClientID, req.Action, req.Data, 30*time.Second)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))

	// Video play API - 支持加密视频的流式解密播放
	http.HandleFunc("/api/video/play", withRecovery(func(w http.ResponseWriter, r *http.Request) {
		targetURL := r.URL.Query().Get("url")
		if targetURL == "" {
			http.Error(w, "url parameter required", http.StatusBadRequest)
			return
		}

		// 获取可选的解密密钥
		decryptKeyStr := r.URL.Query().Get("key")
		var decryptKey uint64
		var needsDecryption bool

		if decryptKeyStr != "" {
			var err error
			decryptKey, err = strconv.ParseUint(decryptKeyStr, 10, 64)
			if err != nil {
				http.Error(w, "invalid decryption key", http.StatusBadRequest)
				return
			}
			needsDecryption = true
		}

		// 创建上游请求
		req, err := http.NewRequest(r.Method, targetURL, nil)
		if err != nil {
			http.Error(w, "invalid URL", http.StatusBadRequest)
			return
		}

		// 复制 Range 头（支持视频拖动）
		if rangeHeader := r.Header.Get("Range"); rangeHeader != "" {
			req.Header.Set("Range", rangeHeader)
		}

		// 发送请求
		client := &http.Client{Timeout: 0}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "failed to fetch video", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// 设置CORS头
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, HEAD, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Range")

		// 复制响应头
		for k, v := range resp.Header {
			w.Header()[k] = v
		}

		// 确保设置 Accept-Ranges
		if w.Header().Get("Accept-Ranges") == "" {
			w.Header().Set("Accept-Ranges", "bytes")
		}

		// 如果需要解密
		if needsDecryption {
			// 解析 Content-Range 头以获取起始偏移
			var startOffset uint64 = 0
			if cr := resp.Header.Get("Content-Range"); cr != "" {
				// Content-Range 格式: "bytes start-end/total"
				parts := strings.Split(cr, " ")
				if len(parts) == 2 {
					rangePart := parts[1]
					dashIdx := strings.Index(rangePart, "-")
					if dashIdx > 0 {
						if v, err := strconv.ParseUint(rangePart[:dashIdx], 10, 64); err == nil {
							startOffset = v
						}
					}
				}
			}

			// 创建解密读取器
			// 加密区域大小为 131072 字节（128KB）
			decryptReader := newDecryptReader(resp.Body, decryptKey, startOffset, 131072)

			// 写入状态码
			w.WriteHeader(resp.StatusCode)

			// 如果是 HEAD 请求，不传输内容
			if r.Method == "HEAD" {
				return
			}

			// 流式复制解密后的数据到客户端
			io.Copy(w, decryptReader)
		} else {
			// 无需解密，直接代理
			w.WriteHeader(resp.StatusCode)

			// 如果是 HEAD 请求，不传输内容
			if r.Method == "HEAD" {
				return
			}

			// 流式复制数据到客户端
			io.Copy(w, resp.Body)
		}
	}))

	// 静态文件服务 - Vue SPA 支持
	// 优先服务 frontend/dist 目录下的静态资源
	fs := http.FileServer(http.Dir("frontend/dist"))
	http.HandleFunc("/", withRecovery(func(w http.ResponseWriter, r *http.Request) {
		// 如果是 API 调用或 WebSocket，不处理
		if strings.HasPrefix(r.URL.Path, "/api/") || strings.HasPrefix(r.URL.Path, "/ws/") {
			return
		}

		path := r.URL.Path
		// 检查文件是否存在于 dist 目录
		if _, err := os.Stat("frontend/dist" + path); os.IsNotExist(err) {
			// 文件不存在，返回 index.html (SPA History Mode)
			http.ServeFile(w, r, "frontend/dist/index.html")
			return
		}

		// 文件存在，直接服务
		fs.ServeHTTP(w, r)
	}))

	log.Println("Hub Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// --- 解密相关代码 (ISAAC64 + DecryptReader) ---

// decryptReader 是一个支持流式解密的 io.Reader 包装器
type decryptReader struct {
	reader   io.Reader
	ctx      *isaac64Ctx
	limit    uint64
	consumed uint64
	ks       [8]byte
	ksPos    int
}

type isaac64Ctx struct {
	randrsl [256]uint64
	randcnt uint64
	mm      [256]uint64
	aa      uint64
	bb      uint64
	cc      uint64
}

func newDecryptReader(reader io.Reader, key uint64, offset uint64, limit uint64) *decryptReader {
	ctx := newIsaac64Context(key)
	dr := &decryptReader{
		reader:   reader,
		ctx:      ctx,
		limit:    limit,
		consumed: 0,
		ksPos:    8,
	}

	if limit > 0 {
		if offset >= limit {
			dr.consumed = limit
		} else {
			dr.consumed = offset
			skipBlocks := offset / 8
			for i := uint64(0); i < skipBlocks; i++ {
				_ = dr.ctx.isaac64Random()
			}
			randNumber := dr.ctx.isaac64Random()
			binary.BigEndian.PutUint64(dr.ks[:], randNumber)
			dr.ksPos = int(offset % 8)
		}
	}
	return dr
}

func (dr *decryptReader) Read(p []byte) (int, error) {
	n, err := dr.reader.Read(p)
	if n <= 0 {
		return n, err
	}

	if dr.limit == 0 || dr.consumed >= dr.limit {
		return n, err
	}

	toDecrypt := uint64(n)
	remaining := dr.limit - dr.consumed
	if toDecrypt > remaining {
		toDecrypt = remaining
	}

	for i := uint64(0); i < toDecrypt; i++ {
		if dr.ksPos >= 8 {
			randNumber := dr.ctx.isaac64Random()
			binary.BigEndian.PutUint64(dr.ks[:], randNumber)
			dr.ksPos = 0
		}
		p[i] ^= dr.ks[dr.ksPos]
		dr.ksPos++
	}
	dr.consumed += toDecrypt
	return n, err
}

func newIsaac64Context(seed uint64) *isaac64Ctx {
	ctx := &isaac64Ctx{}
	ctx.randrsl[0] = seed
	ctx.randinit(true)
	return ctx
}

func (ctx *isaac64Ctx) randinit(flag bool) {
	var a, b, c, d, e, f, g, h uint64
	a = 0x9e3779b97f4a7c13
	b, c, d, e, f, g, h = a, a, a, a, a, a, a

	for j := 0; j < 4; j++ {
		a, b, c, d, e, f, g, h = ctx.mix(a, b, c, d, e, f, g, h)
	}

	for j := 0; j < 256; j += 8 {
		if flag {
			a += ctx.randrsl[j]
			b += ctx.randrsl[j+1]
			c += ctx.randrsl[j+2]
			d += ctx.randrsl[j+3]
			e += ctx.randrsl[j+4]
			f += ctx.randrsl[j+5]
			g += ctx.randrsl[j+6]
			h += ctx.randrsl[j+7]
		}
		a, b, c, d, e, f, g, h = ctx.mix(a, b, c, d, e, f, g, h)
		ctx.mm[j] = a
		ctx.mm[j+1] = b
		ctx.mm[j+2] = c
		ctx.mm[j+3] = d
		ctx.mm[j+4] = e
		ctx.mm[j+5] = f
		ctx.mm[j+6] = g
		ctx.mm[j+7] = h
	}

	if flag {
		for j := 0; j < 256; j += 8 {
			a += ctx.mm[j]
			b += ctx.mm[j+1]
			c += ctx.mm[j+2]
			d += ctx.mm[j+3]
			e += ctx.mm[j+4]
			f += ctx.mm[j+5]
			g += ctx.mm[j+6]
			h += ctx.mm[j+7]
			a, b, c, d, e, f, g, h = ctx.mix(a, b, c, d, e, f, g, h)
			ctx.mm[j] = a
			ctx.mm[j+1] = b
			ctx.mm[j+2] = c
			ctx.mm[j+3] = d
			ctx.mm[j+4] = e
			ctx.mm[j+5] = f
			ctx.mm[j+6] = g
			ctx.mm[j+7] = h
		}
	}

	ctx.isaac64()
	ctx.randcnt = 256
}

func (ctx *isaac64Ctx) mix(a, b, c, d, e, f, g, h uint64) (uint64, uint64, uint64, uint64, uint64, uint64, uint64, uint64) {
	a -= e
	f ^= h >> 9
	h += a
	b -= f
	g ^= a << 9
	a += b
	c -= g
	h ^= b >> 23
	b += c
	d -= h
	a ^= c << 15
	c += d
	e -= a
	b ^= d >> 14
	d += e
	f -= b
	c ^= e << 20
	e += f
	g -= c
	d ^= f >> 17
	f += g
	h -= d
	e ^= g << 14
	g += h
	return a, b, c, d, e, f, g, h
}

func (ctx *isaac64Ctx) isaac64() {
	ctx.cc++
	ctx.bb += ctx.cc

	for j := 0; j < 256; j++ {
		x := ctx.mm[j]
		switch j % 4 {
		case 0:
			ctx.aa = ^(ctx.aa ^ (ctx.aa << 21))
		case 1:
			ctx.aa = ctx.aa ^ (ctx.aa >> 5)
		case 2:
			ctx.aa = ctx.aa ^ (ctx.aa << 12)
		case 3:
			ctx.aa = ctx.aa ^ (ctx.aa >> 33)
		}
		ctx.aa += ctx.mm[(j+128)%256]
		y := ctx.mm[(x>>3)%256] + ctx.aa + ctx.bb
		ctx.mm[j] = y
		ctx.bb = ctx.mm[(y>>11)%256] + x
		ctx.randrsl[j] = ctx.bb
	}
}

func (ctx *isaac64Ctx) isaac64Random() uint64 {
	if ctx.randcnt == 0 {
		ctx.isaac64()
		ctx.randcnt = 256
	}
	ctx.randcnt--
	return ctx.randrsl[ctx.randcnt]
}
