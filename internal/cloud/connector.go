package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"wx_channel/internal/config"
	"wx_channel/internal/utils"
	hubws "wx_channel/internal/websocket"

	"github.com/gorilla/websocket"
)

// Connector 云端连接器
type Connector struct {
	cfg   *config.Config
	local *hubws.Hub
	conn  *websocket.Conn
	mu    sync.Mutex

	clientID string
	ctx      context.Context
	cancel   context.CancelFunc
}

// NewConnector 创建云端连接器
func NewConnector(cfg *config.Config, localHub *hubws.Hub) *Connector {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Connector{
		cfg:      cfg,
		local:    localHub,
		clientID: cfg.MachineID,
		ctx:      ctx,
		cancel:   cancel,
	}

	if c.clientID == "" {
		hostname, _ := os.Hostname()
		if hostname == "" {
			hostname = "unknown"
		}
		c.clientID = fmt.Sprintf("%s-%d", hostname, time.Now().Unix()%10000)
	}

	return c
}

// Start 启动连接器
func (c *Connector) Start() {
	if c.cfg.CloudHubURL == "" {
		utils.LogInfo("云端管理未启用 (未配置 cloud_hub_url)")
		return
	}

	utils.LogInfo("正在启动云端连接器 (ID: %s, URL: %s)", c.clientID, c.cfg.CloudHubURL)

	go c.connectLoop()
}

// Stop 停止连接器
func (c *Connector) Stop() {
	c.cancel()
	c.mu.Lock()
	if c.conn != nil {
		c.conn.Close()
	}
	c.mu.Unlock()
}

func (c *Connector) connectLoop() {
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			err := c.connect()
			if err != nil {
				utils.LogError("云端 Hub 连接失败: %v, 5秒后重试...", err)
				time.Sleep(5 * time.Second)
				continue
			}

			utils.LogInfo("✓ 已连接到云端 Hub")
			c.handleConnection()
			utils.LogWarn("云端 Hub 连接已断开，尝试重新连接...")
			time.Sleep(3 * time.Second)
		}
	}
}

func (c *Connector) connect() error {
	header := http.Header{}
	if c.cfg.CloudSecret != "" {
		header.Add("X-Cloud-Secret", c.cfg.CloudSecret)
	}
	header.Add("X-Client-ID", c.clientID)

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(c.cfg.CloudHubURL, header)
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.conn = conn
	c.mu.Unlock()
	return nil
}

func (c *Connector) handleConnection() {
	// 启动心跳
	go c.heartbeatLoop()

	// 监听消息
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		var msg CloudMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			utils.LogError("云端消息解析失败: %v", err)
			continue
		}

		go c.processMessage(msg)
	}
}

func (c *Connector) heartbeatLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			hostname, _ := os.Hostname()
			payload := HeartbeatPayload{
				Hostname: hostname,
				Version:  c.cfg.Version,
				Status:   "running",
			}
			payloadData, _ := json.Marshal(payload)

			msg := CloudMessage{
				ID:        fmt.Sprintf("hb-%d", time.Now().Unix()),
				Type:      MsgTypeHeartbeat,
				ClientID:  c.clientID,
				Payload:   payloadData,
				Timestamp: time.Now().Unix(),
			}

			if err := c.send(msg); err != nil {
				utils.LogError("心跳发送失败: %v", err)
				return // 触发重连
			}
		}
	}
}

func (c *Connector) send(msg CloudMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return fmt.Errorf("connection closed")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *Connector) processMessage(msg CloudMessage) {
	if msg.Type != MsgTypeCommand {
		return
	}

	var cmd CommandPayload
	if err := json.Unmarshal(msg.Payload, &cmd); err != nil {
		c.sendError(msg.ID, "Invalid command payload")
		return
	}

	utils.LogInfo("收到云端指令: %s", cmd.Action)

	switch cmd.Action {
	case "api_call":
		c.handleAPICall(msg.ID, cmd.Data)
	default:
		c.sendError(msg.ID, fmt.Sprintf("Unknown action: %s", cmd.Action))
	}
}

func (c *Connector) handleAPICall(reqID string, data json.RawMessage) {
	var call struct {
		Key  string          `json:"key"`
		Body json.RawMessage `json:"body"`
	}

	if err := json.Unmarshal(data, &call); err != nil {
		c.sendError(reqID, "Invalid API call parameters")
		return
	}

	// 调用本地 Hub
	respData, err := c.local.CallAPI(call.Key, call.Body, 30*time.Second)
	if err != nil {
		c.sendError(reqID, err.Error())
		return
	}

	// 返回结果
	c.sendResponse(reqID, true, respData, "")
}

func (c *Connector) sendResponse(reqID string, success bool, data json.RawMessage, errMsg string) {
	resp := ResponsePayload{
		RequestID: reqID,
		Success:   success,
		Data:      data,
		Error:     errMsg,
	}
	respData, _ := json.Marshal(resp)

	msg := CloudMessage{
		ID:        fmt.Sprintf("resp-%s", reqID),
		Type:      MsgTypeResponse,
		ClientID:  c.clientID,
		Payload:   respData,
		Timestamp: time.Now().Unix(),
	}

	_ = c.send(msg)
}

func (c *Connector) sendError(reqID string, errMsg string) {
	c.sendResponse(reqID, false, nil, errMsg)
}
