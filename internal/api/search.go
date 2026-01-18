package api

import (
	"encoding/json"
	"net/http"
	"time"

	"wx_channel/internal/websocket"
)

// SearchService 搜索服务
type SearchService struct {
	hub *websocket.Hub
}

// NewSearchService 创建搜索服务
func NewSearchService(hub *websocket.Hub) *SearchService {
	return &SearchService{hub: hub}
}

// SearchContact 搜索账号
func (s *SearchService) SearchContact(w http.ResponseWriter, r *http.Request) {
	keyword := r.URL.Query().Get("keyword")
	if keyword == "" {
		http.Error(w, "keyword is required", http.StatusBadRequest)
		return
	}

	// 调用前端 API
	body := websocket.SearchContactBody{
		Keyword: keyword,
	}

	data, err := s.hub.CallAPI("key:channels:contact_list", body, 20*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GetFeedList 获取账号的视频列表
func (s *SearchService) GetFeedList(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	nextMarker := r.URL.Query().Get("next_marker")

	// 调用前端 API
	body := websocket.FeedListBody{
		Username:   username,
		NextMarker: nextMarker,
	}

	data, err := s.hub.CallAPI("key:channels:feed_list", body, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GetFeedProfile 获取视频详情
func (s *SearchService) GetFeedProfile(w http.ResponseWriter, r *http.Request) {
	objectID := r.URL.Query().Get("object_id")
	nonceID := r.URL.Query().Get("nonce_id")
	url := r.URL.Query().Get("url")

	if objectID == "" && url == "" {
		http.Error(w, "object_id or url is required", http.StatusBadRequest)
		return
	}

	// 调用前端 API
	body := websocket.FeedProfileBody{
		ObjectID: objectID,
		NonceID:  nonceID,
		URL:      url,
	}

	data, err := s.hub.CallAPI("key:channels:feed_profile", body, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 返回结果
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

// GetStatus 获取 WebSocket 连接状态
func (s *SearchService) GetStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"connected": s.hub.ClientCount() > 0,
		"clients":   s.hub.ClientCount(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
