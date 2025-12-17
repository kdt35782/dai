package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"sm-medical/pkg/config"
	"sm-medical/pkg/utils"
	"time"
)

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

// Upload 上传文件
func (h *FileHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, "请选择文件")
		return
	}

	fileType := c.PostForm("fileType")
	if fileType == "" {
		fileType = "general"
	}

	// 检查文件大小
	if file.Size > config.AppConfig.Upload.MaxSize {
		utils.BadRequest(c, "文件大小超过限制")
		return
	}

	// 生成文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", fileType, time.Now().UnixNano(), ext)
	savePath := filepath.Join(config.AppConfig.Upload.UploadPath, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		utils.InternalError(c, "文件保存失败")
		return
	}

	// 返回文件URL（实际项目需要配置文件服务器地址）
	fileURL := fmt.Sprintf("/uploads/%s", filename)

	utils.SuccessWithMessage(c, "上传成功", gin.H{
		"fileUrl":      fileURL,
		"fileName":     file.Filename,
		"fileSize":     file.Size,
		"uploadTime":   time.Now().Format("2006-01-02 15:04:05"),
	})
}
