package websocket

import "encoding/json"

// WebSocket 消息类型
type WSMessageType string

const (
	WSMessageTypeAPICall     WSMessageType = "api_call"
	WSMessageTypeAPIResponse WSMessageType = "api_response"
	WSMessageTypePing        WSMessageType = "ping"
	WSMessageTypePong        WSMessageType = "pong"
)

// WebSocket 消息
type WSMessage struct {
	Type WSMessageType   `json:"type"`
	Data json.RawMessage `json:"data,omitempty"`
}

// API 调用请求
type APICallRequest struct {
	ID   string      `json:"id"`
	Key  string      `json:"key"`
	Body interface{} `json:"body"`
}

// API 调用响应
type APICallResponse struct {
	ID      string          `json:"id"`
	Data    json.RawMessage `json:"data"`
	ErrCode int             `json:"errCode,omitempty"`
	ErrMsg  string          `json:"errMsg,omitempty"`
}

// 搜索账号请求体
type SearchContactBody struct {
	Keyword string `json:"keyword"`
}

// 获取账号视频列表请求体
type FeedListBody struct {
	Username   string `json:"username"`
	NextMarker string `json:"next_marker"`
}

// 获取视频详情请求体
type FeedProfileBody struct {
	ObjectID string `json:"objectId"`
	NonceID  string `json:"nonceId"`
	URL      string `json:"url"`
}
