package websocket

import (
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client 表示一个 WebSocket 客户端连接
type Client struct {
	conn     *websocket.Conn
	send     chan []byte
	hub      *Hub
	mu       sync.Mutex
	closed   bool
	lastPing time.Time
}

// NewClient 创建新的客户端
func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		conn:     conn,
		send:     make(chan []byte, 256),
		hub:      hub,
		lastPing: time.Now(),
	}
}

// ReadPump 从 WebSocket 连接读取消息
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		c.lastPing = time.Now()
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 记录非预期的关闭错误
			}
			break
		}

		// 解析消息
		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		// 处理 API 响应
		if msg.Type == WSMessageTypeAPIResponse {
			var resp APICallResponse
			if err := json.Unmarshal(msg.Data, &resp); err != nil {
				continue
			}
			c.hub.handleAPIResponse(resp)
		}
	}
}

// WritePump 向 WebSocket 连接写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Send 发送消息到客户端
func (c *Client) Send(data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return errors.New("client is closed")
	}

	select {
	case c.send <- data:
		return nil
	default:
		return errors.New("send buffer is full")
	}
}

// Close 关闭客户端连接
func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.closed {
		c.closed = true
		close(c.send)
	}
}
