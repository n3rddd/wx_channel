package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
}

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
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("[WebSocket] 连接升级失败: %v\n", err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	client := NewClient(conn, h.hub)
	h.hub.RegisterClient(client)

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump()
}
