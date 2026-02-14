package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"wx_channel/internal/utils"

	"github.com/qtgolang/SunnyNet/SunnyNet"
)

// HandleInitUpload 处理分片上传初始化请求
func (h *UploadHandler) HandleInitUpload(Conn *SunnyNet.HttpConn) bool {
	path := Conn.Request.URL.Path
	if path != "/__wx_channels_api/upload/init" {
		return false
	}

	// 权限验证
	if !h.authenticate(Conn) {
		return true
	}

	var req struct {
		Filename  string `json:"filename"`
		TotalSize int64  `json:"totalSize"`
		ChunkSize int64  `json:"chunkSize"`
		Total     int    `json:"total"`
		Type      string `json:"type"` // "video" or "cover"
	}

	body, err := io.ReadAll(Conn.Request.Body)
	if err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("read body failed: %v", err))
		return true
	}
	defer Conn.Request.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("invalid json: %v", err))
		return true
	}

	if req.Filename == "" || req.TotalSize <= 0 || req.Total <= 0 {
		h.sendErrorResponse(Conn, fmt.Errorf("invalid parameters"))
		return true
	}

	// 生成上传ID
	uploadID := fmt.Sprintf("%d_%s", time.Now().UnixNano(), utils.RandomString(8))

	// 创建临时目录
	tempDir := filepath.Join(os.TempDir(), "wx_channel_upload", uploadID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("create temp dir failed: %v", err))
		return true
	}

	// 保存元数据
	metaFile := filepath.Join(tempDir, "meta.json")
	metaData, _ := json.Marshal(req)
	os.WriteFile(metaFile, metaData, 0644)

	utils.LogInfo("[Upload] Init: %s (size: %d, chunks: %d)", req.Filename, req.TotalSize, req.Total)

	h.sendJSONResponse(Conn, 200, map[string]interface{}{
		"code":     0,
		"message":  "success",
		"uploadId": uploadID,
	})
	return true
}

// HandleUploadChunk 处理分片上传请求
func (h *UploadHandler) HandleUploadChunk(Conn *SunnyNet.HttpConn) bool {
	path := Conn.Request.URL.Path
	if path != "/__wx_channels_api/upload/chunk" {
		return false
	}

	// 权限验证
	if !h.authenticate(Conn) {
		return true
	}

	// 限制并发上传数
	h.chunkSem <- struct{}{}
	defer func() { <-h.chunkSem }()

	// 解析 multipart form
	// 增加最大内存限制到 10MB
	if err := Conn.Request.ParseMultipartForm(10 << 20); err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("parse multipart failed: %v", err))
		return true
	}

	uploadID := Conn.Request.FormValue("uploadId")
	chunkIndex := Conn.Request.FormValue("chunkIndex")
	file, _, err := Conn.Request.FormFile("chunk")
	if err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("get chunk file failed: %v", err))
		return true
	}
	defer file.Close()

	if uploadID == "" || chunkIndex == "" {
		h.sendErrorResponse(Conn, fmt.Errorf("missing uploadId or chunkIndex"))
		return true
	}

	idx, err := strconv.Atoi(chunkIndex)
	if err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("invalid chunkIndex"))
		return true
	}

	// 验证临时目录
	tempDir := filepath.Join(os.TempDir(), "wx_channel_upload", uploadID)
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		h.sendErrorResponse(Conn, fmt.Errorf("upload session not found"))
		return true
	}

	// 保存分片
	chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", idx))
	out, err := os.Create(chunkPath)
	if err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("create chunk file failed: %v", err))
		return true
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("write chunk failed: %v", err))
		return true
	}

	// utils.LogDebug("[Upload] Chunk %d saved", idx)

	h.sendJSONResponse(Conn, 200, map[string]interface{}{
		"code":    0,
		"message": "success",
	})
	return true
}

// HandleCompleteUpload 处理分片上传完成请求
func (h *UploadHandler) HandleCompleteUpload(Conn *SunnyNet.HttpConn) bool {
	path := Conn.Request.URL.Path
	if path != "/__wx_channels_api/upload/complete" {
		return false
	}

	// 权限验证
	if !h.authenticate(Conn) {
		return true
	}

	// 限制并发合并数
	h.mergeSem <- struct{}{}
	defer func() { <-h.mergeSem }()

	var req struct {
		UploadID string `json:"uploadId"`
	}

	body, err := io.ReadAll(Conn.Request.Body)
	if err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("read body failed: %v", err))
		return true
	}
	defer Conn.Request.Body.Close()

	if err := json.Unmarshal(body, &req); err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("invalid json: %v", err))
		return true
	}

	tempDir := filepath.Join(os.TempDir(), "wx_channel_upload", req.UploadID)
	metaFile := filepath.Join(tempDir, "meta.json")

	// 读取元数据
	metaBytes, err := os.ReadFile(metaFile)
	if err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("read meta failed: %v", err))
		return true
	}

	var meta struct {
		Filename  string `json:"filename"`
		TotalSize int64  `json:"totalSize"`
		Total     int    `json:"total"`
		Type      string `json:"type"`
	}
	json.Unmarshal(metaBytes, &meta)

	// 检查所有分片是否存在
	for i := 0; i < meta.Total; i++ {
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", i))
		if _, err := os.Stat(chunkPath); os.IsNotExist(err) {
			h.sendErrorResponse(Conn, fmt.Errorf("chunk %d missing", i))
			return true
		}
	}

	// 获取下载目录
	downloadsDir, err := h.getDownloadsDir()
	if err != nil {
		h.sendErrorResponse(Conn, err)
		return true
	}

	// 确保目录存在
	if err := utils.EnsureDir(downloadsDir); err != nil {
		h.sendErrorResponse(Conn, err)
		return true
	}

	// 合并文件
	targetPath := filepath.Join(downloadsDir, utils.CleanFilename(meta.Filename))
	// 如果文件已存在，添加随机后缀
	if _, err := os.Stat(targetPath); err == nil {
		ext := filepath.Ext(meta.Filename)
		name := strings.TrimSuffix(meta.Filename, ext)
		targetPath = filepath.Join(downloadsDir, fmt.Sprintf("%s_%s%s", name, utils.RandomString(4), ext))
	}

	outFile, err := os.Create(targetPath)
	if err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("create target file failed: %v", err))
		return true
	}
	defer outFile.Close()

	for i := 0; i < meta.Total; i++ {
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", i))
		chunkFile, err := os.Open(chunkPath)
		if err != nil {
			h.sendErrorResponse(Conn, fmt.Errorf("read chunk %d failed: %v", i, err))
			return true
		}
		io.Copy(outFile, chunkFile)
		chunkFile.Close()
	}

	utils.LogInfo("[Upload] Completed: %s -> %s", meta.Filename, targetPath)

	// 清理临时文件
	go os.RemoveAll(tempDir)

	h.sendJSONResponse(Conn, 200, map[string]interface{}{
		"code":     0,
		"message":  "success",
		"filePath": targetPath,
	})
	return true
}

// HandleUploadStatus 查询已上传的分片列表
func (h *UploadHandler) HandleUploadStatus(Conn *SunnyNet.HttpConn) bool {
	path := Conn.Request.URL.Path
	if path != "/__wx_channels_api/upload/status" {
		return false
	}

	// 权限验证
	if !h.authenticate(Conn) {
		return true
	}

	uploadID := Conn.Request.URL.Query().Get("uploadId")
	if uploadID == "" {
		h.sendErrorResponse(Conn, fmt.Errorf("missing uploadId"))
		return true
	}

	tempDir := filepath.Join(os.TempDir(), "wx_channel_upload", uploadID)
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		// 会话不存在
		h.sendJSONResponse(Conn, 200, map[string]interface{}{
			"code":    0,
			"message": "session not found",
			"chunks":  []int{},
		})
		return true
	}

	// 扫描已存在的分片
	var chunks []int
	files, err := os.ReadDir(tempDir)
	if err == nil {
		for _, file := range files {
			if strings.HasPrefix(file.Name(), "chunk_") {
				idxStr := strings.TrimPrefix(file.Name(), "chunk_")
				if idx, err := strconv.Atoi(idxStr); err == nil {
					chunks = append(chunks, idx)
				}
			}
		}
	}

	h.sendJSONResponse(Conn, 200, map[string]interface{}{
		"code":    0,
		"message": "success",
		"chunks":  chunks,
	})
	return true
}

// HandleCancelDownload 处理取消下载请求
func (h *UploadHandler) HandleCancelDownload(Conn *SunnyNet.HttpConn) bool {
	path := Conn.Request.URL.Path
	if path != "/__wx_channels_api/download/cancel" {
		return false
	}

	// 权限验证
	if !h.authenticate(Conn) {
		return true
	}

	taskID := Conn.Request.URL.Query().Get("taskId")
	if taskID == "" {
		// 尝试从 body 读取
		var req struct {
			TaskID string `json:"taskId"`
		}
		body, _ := io.ReadAll(Conn.Request.Body)
		Conn.Request.Body.Close()
		json.Unmarshal(body, &req)
		taskID = req.TaskID
	}

	if taskID != "" && h.gopeedService != nil {
		utils.LogInfo("[Download] Request cancel task: %s", taskID)
		err := h.gopeedService.DeleteTask(taskID)
		if err != nil {
			utils.Warn("[Download] Cancel task failed: %v", err)
		} else {
			utils.LogInfo("[Download] Task cancelled success")
		}
	}

	h.sendSuccessResponse(Conn)
	return true
}
