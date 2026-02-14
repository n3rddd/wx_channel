package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"wx_channel/internal/config"
	"wx_channel/internal/response"
	"wx_channel/internal/services"
	"wx_channel/internal/utils"
	internalws "wx_channel/internal/websocket"

	"github.com/qtgolang/SunnyNet/SunnyNet"
)

// UploadHandler 文件上传处理器
type UploadHandler struct {
	cfg             *config.Config
	downloadService *services.DownloadRecordService
	gopeedService   *services.GopeedService
	chunkSem        chan struct{}
	mergeSem        chan struct{}
	wsHub           *internalws.Hub
}

// NewUploadHandler 创建上传处理器
func NewUploadHandler(cfg *config.Config, wsHub *internalws.Hub, gopeedService *services.GopeedService) *UploadHandler {
	return &UploadHandler{
		cfg:             cfg,
		downloadService: services.NewDownloadRecordService(),
		gopeedService:   gopeedService,
		chunkSem:        make(chan struct{}, 5), // 限制并发分片上传数
		mergeSem:        make(chan struct{}, 2), // 限制并发合并数
		wsHub:           wsHub,
	}
}

// getConfig 获取当前配置（动态获取最新配置）
func (h *UploadHandler) getConfig() *config.Config {
	if h.cfg != nil {
		return h.cfg
	}
	return config.Get()
}

// getDownloadsDir 获取解析后的下载目录
func (h *UploadHandler) getDownloadsDir() (string, error) {
	cfg := h.getConfig()
	if cfg == nil {
		return "", fmt.Errorf("config not initialized")
	}
	return utils.ResolveDownloadDir(cfg.DownloadsDir)
}

// Handle implements router.Interceptor
func (h *UploadHandler) Handle(Conn *SunnyNet.HttpConn) bool {
	if Conn.Request == nil || Conn.Request.URL == nil {
		return false
	}

	// Add local panic recovery
	defer func() {
		if r := recover(); r != nil {
			utils.Error("UploadHandler.Handle panic: %v", r)
		}
	}()

	path := Conn.Request.URL.Path

	// 路由分发
	if strings.HasPrefix(path, "/__wx_channels_api/upload/init") {
		return h.HandleInitUpload(Conn)
	}
	if strings.HasPrefix(path, "/__wx_channels_api/upload/chunk") {
		return h.HandleUploadChunk(Conn)
	}
	if strings.HasPrefix(path, "/__wx_channels_api/upload/complete") {
		return h.HandleCompleteUpload(Conn)
	}
	if strings.HasPrefix(path, "/__wx_channels_api/upload/status") {
		return h.HandleUploadStatus(Conn)
	}

	if strings.HasPrefix(path, "/__wx_channels_api/save/video") {
		return h.HandleSaveVideo(Conn)
	}
	if strings.HasPrefix(path, "/__wx_channels_api/save/cover") {
		return h.HandleSaveCover(Conn)
	}

	if strings.HasPrefix(path, "/__wx_channels_api/download/video") {
		return h.HandleDownloadVideo(Conn)
	}
	if strings.HasPrefix(path, "/__wx_channels_api/download/cancel") {
		return h.HandleCancelDownload(Conn)
	}

	return false
}

// sendSuccessResponse 发送成功响应
func (h *UploadHandler) sendSuccessResponse(Conn *SunnyNet.HttpConn) {
	h.sendJSONResponse(Conn, 200, map[string]interface{}{
		"code":    0,
		"message": "success",
	})
}

// sendJSONResponse 发送JSON响应
func (h *UploadHandler) sendJSONResponse(Conn *SunnyNet.HttpConn, statusCode int, data interface{}) {
	headers := http.Header{}
	headers.Set("Content-Type", "application/json")
	body, err := json.Marshal(data)
	if err != nil {
		utils.LogError("JSON marshal failed: %v", err)
		Conn.StopRequest(500, `{"code":-1,"message":"internal server error"}`, headers)
		return
	}
	Conn.StopRequest(statusCode, string(body), headers)
}

// sendErrorResponse 发送错误响应
func (h *UploadHandler) sendErrorResponse(Conn *SunnyNet.HttpConn, err error) {
	utils.LogError("UploadHandler Error: %v", err)
	h.sendJSONResponse(Conn, 200, map[string]interface{}{
		"code":    -1,
		"message": err.Error(),
	})
}

// authenticate 验证请求
func (h *UploadHandler) authenticate(Conn *SunnyNet.HttpConn) bool {
	cfg := h.getConfig()
	if cfg != nil && cfg.SecretToken != "" {
		clientToken := Conn.Request.Header.Get("X-Local-Auth")
		if clientToken != cfg.SecretToken {
			headers := http.Header{}
			headers.Set("Content-Type", "application/json")
			Conn.StopRequest(401, string(response.ErrorJSON(401, "unauthorized")), headers)
			return false
		}
	}
	return true
}
