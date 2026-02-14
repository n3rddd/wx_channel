package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"wx_channel/internal/utils"

	"github.com/qtgolang/SunnyNet/SunnyNet"
)

// HandleSaveVideo 处理直接保存视频文件请求
func (h *UploadHandler) HandleSaveVideo(Conn *SunnyNet.HttpConn) bool {
	path := Conn.Request.URL.Path
	if path != "/__wx_channels_api/save/video" {
		return false
	}

	// 权限验证
	if !h.authenticate(Conn) {
		return true
	}

	// 解析 multipart form
	// 500MB 大小限制
	if err := Conn.Request.ParseMultipartForm(500 << 20); err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("parse multipart failed: %v", err))
		return true
	}

	file, header, err := Conn.Request.FormFile("video")
	if err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("get video file failed: %v", err))
		return true
	}
	defer file.Close()

	// 获取其他表单字段
	filename := Conn.Request.FormValue("filename")
	if filename == "" {
		filename = header.Filename
	}
	title := Conn.Request.FormValue("title")
	author := Conn.Request.FormValue("author")
	videoID := Conn.Request.FormValue("videoId")

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

	// 生成目标路径
	targetPath := filepath.Join(downloadsDir, utils.CleanFilename(filename))

	// 如果文件已存在，添加随机后缀
	if _, err := os.Stat(targetPath); err == nil {
		ext := filepath.Ext(filename)
		name := strings.TrimSuffix(filename, ext)
		targetPath = filepath.Join(downloadsDir, fmt.Sprintf("%s_%s%s", name, utils.RandomString(4), ext))
	}

	// 保存文件
	out, err := os.Create(targetPath)
	if err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("create file failed: %v", err))
		return true
	}
	defer out.Close()

	written, err := io.Copy(out, file)
	if err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("save file failed: %v", err))
		return true
	}

	utils.LogInfo("[Save] Video saved: %s (size: %.2f MB)", targetPath, float64(written)/(1024*1024))

	// 记录到数据库 (如果有 downloadService)
	if h.downloadService != nil {
		h.downloadService.AddRecord(videoID, title, author, targetPath, written, "completed")
	}

	h.sendJSONResponse(Conn, 200, map[string]interface{}{
		"code":     0,
		"message":  "success",
		"filePath": targetPath,
	})
	return true
}

// HandleSaveCover 处理保存封面图片请求
func (h *UploadHandler) HandleSaveCover(Conn *SunnyNet.HttpConn) bool {
	path := Conn.Request.URL.Path
	if path != "/__wx_channels_api/save/cover" {
		return false
	}

	// 权限验证
	if !h.authenticate(Conn) {
		return true
	}

	// 解析 multipart form
	if err := Conn.Request.ParseMultipartForm(10 << 20); err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("parse multipart failed: %v", err))
		return true
	}

	file, header, err := Conn.Request.FormFile("cover")
	if err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("get cover file failed: %v", err))
		return true
	}
	defer file.Close()

	filename := Conn.Request.FormValue("filename")
	if filename == "" {
		filename = header.Filename
	}

	// 获取下载目录
	downloadsDir, err := h.getDownloadsDir()
	if err != nil {
		h.sendErrorResponse(Conn, err)
		return true
	}

	coversDir := filepath.Join(downloadsDir, "covers")

	// 确保目录存在
	if err := utils.EnsureDir(coversDir); err != nil {
		h.sendErrorResponse(Conn, err)
		return true
	}

	targetPath := filepath.Join(coversDir, utils.CleanFilename(filename))

	out, err := os.Create(targetPath)
	if err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("create file failed: %v", err))
		return true
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		h.sendErrorResponse(Conn, fmt.Errorf("save file failed: %v", err))
		return true
	}

	utils.LogInfo("[Save] Cover saved: %s", targetPath)

	h.sendJSONResponse(Conn, 200, map[string]interface{}{
		"code":     0,
		"message":  "success",
		"filePath": targetPath,
	})
	return true
}
