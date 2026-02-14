package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"wx_channel/internal/utils"

	"github.com/qtgolang/SunnyNet/SunnyNet"
)

type VideoDownloadRequest struct {
	VideoURL  string `json:"videoUrl"`
	VideoID   string `json:"videoId"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Key       string `json:"key"` // 解密KEY
	ForceSave bool   `json:"forceSave"`
}

// HandleDownloadVideo 处理从URL下载视频请求
func (h *UploadHandler) HandleDownloadVideo(Conn *SunnyNet.HttpConn) bool {
	path := Conn.Request.URL.Path
	if path != "/__wx_channels_api/download/video" {
		return false
	}

	// 权限验证
	if !h.authenticate(Conn) {
		return true
	}

	var req VideoDownloadRequest
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

	if req.VideoURL == "" {
		h.sendErrorResponse(Conn, fmt.Errorf("missing videoUrl"))
		return true
	}

	// 异步处理下载
	go func() {
		// 创建上下文
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
		defer cancel()

		// 获取下载目录
		downloadsDir, err := h.getDownloadsDir()
		if err != nil {
			utils.LogError("[Download] Get downloads dir failed: %v", err)
			return
		}

		if err := utils.EnsureDir(downloadsDir); err != nil {
			utils.LogError("[Download] Create downloads dir failed: %v", err)
			return
		}

		// 生成文件名
		filename := utils.CleanFilename(req.Title)
		if filename == "" {
			filename = req.VideoID
		}
		if filename == "" {
			// 使用URL哈希作为文件名
			sum := sha256.Sum256([]byte(req.VideoURL))
			filename = hex.EncodeToString(sum[:])[:16]
		}
		if !strings.HasSuffix(filename, ".mp4") {
			filename += ".mp4"
		}

		targetPath := filepath.Join(downloadsDir, filename)
		// 如果文件已存在，添加随机后缀
		if _, err := os.Stat(targetPath); err == nil {
			ext := filepath.Ext(filename)
			name := strings.TrimSuffix(filename, ext)
			targetPath = filepath.Join(downloadsDir, fmt.Sprintf("%s_%s%s", name, utils.RandomString(4), ext))
		}

		utils.LogInfo("[Download] Start downloading: %s -> %s", req.VideoURL, targetPath)

		// 记录开始下载 (DownloadService)
		if h.downloadService != nil {
			// h.downloadService.AddRecord(req.VideoID, req.Title, req.Author, targetPath, 0, "downloading")
		}

		// 创建 HTTP 请求
		httpReq, err := http.NewRequestWithContext(ctx, "GET", req.VideoURL, nil)
		if err != nil {
			utils.LogError("[Download] Create request failed: %v", err)
			return
		}

		// 设置 Header
		httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
		httpReq.Header.Set("Referer", "https://channels.weixin.qq.com/")

		client := &http.Client{
			Timeout: 2 * time.Hour,
			Transport: &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    30 * time.Second,
				DisableCompression: true, // 我们自己处理解密，不需要自动解压
			},
		}

		needDecrypt := req.Key != ""
		var written int64 = 0
		var expectedTotalSize int64 = 0

		// 重试逻辑
		err = h.downloadVideoWithRetry(ctx, client, httpReq, req, targetPath, needDecrypt, 0, &written, &expectedTotalSize)

		if err != nil {
			utils.LogError("[Download] Failed: %v", err)
			if h.downloadService != nil {
				// h.downloadService.UpdateStatus(req.VideoID, "failed", err.Error())
			}
		} else {
			// duration := 0 // 持续时间未知
			sizeMB := float64(written) / (1024 * 1024)
			utils.LogInfo("[Download] Success: %s (size: %.2f MB)", targetPath, sizeMB)
			if h.downloadService != nil {
				h.downloadService.AddRecord(req.VideoID, req.Title, req.Author, targetPath, written, "completed")
			}
		}
	}()

	h.sendJSONResponse(Conn, 200, map[string]interface{}{
		"code":    0,
		"message": "download started",
	})
	return true
}

// downloadVideoWithRetry 执行一次视频下载尝试
func (h *UploadHandler) downloadVideoWithRetry(ctx context.Context, client *http.Client, httpReq *http.Request, req VideoDownloadRequest, videoPath string, needDecrypt bool, resumeOffset int64, written *int64, expectedTotalSize *int64) error {

	// 如果是从断点续传
	if resumeOffset > 0 {
		httpReq.Header.Set("Range", fmt.Sprintf("bytes=%d-", resumeOffset))
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 获取总大小
	contentLength := resp.ContentLength
	if contentLength > 0 {
		*expectedTotalSize = contentLength + resumeOffset
	}

	// 打开文件（追加模式或创建模式）
	flags := os.O_CREATE | os.O_WRONLY
	if resumeOffset > 0 {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	outFile, err := os.OpenFile(videoPath, flags, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()

	if resumeOffset > 0 {
		outFile.Seek(resumeOffset, io.SeekStart)
	}

	// 如果需要解密
	// 注意：服务端解密需要更完整的实现，参考 decrypt_reader.go
	var reader io.Reader = resp.Body
	if needDecrypt {
		// 临时占位，实际需要解密流
	}

	buffer := make([]byte, 32*1024)
	lastReportTime := time.Now()

	for {
		// 检查取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, readErr := reader.Read(buffer)
		if n > 0 {
			w, writeErr := outFile.Write(buffer[:n])
			if writeErr != nil {
				return writeErr
			}
			*written += int64(w)

			// 报告进度
			if time.Since(lastReportTime) > 5*time.Second {
				lastReportTime = time.Now()
			}
		}

		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			return readErr
		}
	}

	return nil
}
