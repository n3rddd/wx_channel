package websocket

import (
	"fmt"
	"net/http"

	"github.com/coder/websocket"
)

// Handler WebSocket HTTP 处理器
type Handler struct {
	hub *Hub
}

// NewHandler 创建新的处理器
func NewHandler(hub *Hub) *Handler {
	return &Handler{hub: hub}
}

// ServeHTTP 处理 WebSocket 连接请求
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // 允许所有来源
		CompressionMode:    websocket.CompressionContextTakeover,
	})
	if err != nil {
		fmt.Printf("[WebSocket] 连接升级失败: %v\n", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	// 获取远程地址
	remoteAddr := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		remoteAddr = realIP
	} else if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		remoteAddr = forwarded
	}

	client := NewClientWithAddr(conn, h.hub, remoteAddr)
	h.hub.RegisterClient(client)

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump()
}
