package websocket

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Hub 管理所有 WebSocket 客户端连接
type Hub struct {
	clients      map[*Client]bool
	register     chan *Client
	unregister   chan *Client
	mu           sync.RWMutex
	lastClient   *Client // 最后注册的客户端

	// API 调用管理
	requests   map[string]chan APICallResponse
	requestsMu sync.RWMutex
	reqSeq     uint64
}

// NewHub 创建新的 Hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		requests:   make(map[string]chan APICallResponse),
	}
}

// Run 启动 Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.lastClient = client // 记录最后注册的客户端
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				client.Close()
				// 如果注销的是最后一个客户端，清除引用
				if h.lastClient == client {
					h.lastClient = nil
					// 尝试找到另一个活跃的客户端
					for c := range h.clients {
						h.lastClient = c
						break
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

// RegisterClient 注册新客户端
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// GetClient 获取一个可用的客户端（优先返回最后注册的客户端）
func (h *Hub) GetClient() (*Client, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// 优先使用最后注册的客户端
	if h.lastClient != nil {
		if _, ok := h.clients[h.lastClient]; ok {
			return h.lastClient, nil
		}
	}

	// 如果最后注册的客户端不可用，使用任意一个
	for client := range h.clients {
		return client, nil
	}
	
	return nil, errors.New("no available client")
}

// ClientCount 返回当前连接的客户端数量
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// CallAPI 调用前端 API
func (h *Hub) CallAPI(key string, body interface{}, timeout time.Duration) (json.RawMessage, error) {
	client, err := h.GetClient()
	if err != nil {
		return nil, err
	}

	// 生成请求 ID
	id := atomic.AddUint64(&h.reqSeq, 1)
	reqID := fmt.Sprintf("%d", id)

	// 创建响应通道
	respChan := make(chan APICallResponse, 1)
	h.requestsMu.Lock()
	h.requests[reqID] = respChan
	h.requestsMu.Unlock()

	defer func() {
		h.requestsMu.Lock()
		delete(h.requests, reqID)
		h.requestsMu.Unlock()
	}()

	// 构建请求消息
	req := APICallRequest{
		ID:   reqID,
		Key:  key,
		Body: body,
	}

	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	msg := WSMessage{
		Type: WSMessageTypeAPICall,
		Data: reqData,
	}

	msgData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// 发送请求
	if err := client.Send(msgData); err != nil {
		return nil, err
	}

	// 等待响应
	select {
	case resp := <-respChan:
		if resp.ErrCode != 0 {
			return nil, errors.New(resp.ErrMsg)
		}
		return resp.Data, nil
	case <-time.After(timeout):
		return nil, errors.New("request timeout")
	}
}

// handleAPIResponse 处理 API 响应
func (h *Hub) handleAPIResponse(resp APICallResponse) {
	h.requestsMu.RLock()
	respChan, ok := h.requests[resp.ID]
	h.requestsMu.RUnlock()

	if ok {
		select {
		case respChan <- resp:
			// 响应已发送
		default:
			// 响应通道已满（不应该发生）
		}
	}
}
