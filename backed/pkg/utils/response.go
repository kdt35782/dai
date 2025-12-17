package utils

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:      200,
		Message:   "操作成功",
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Response{
		Code:      200,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(200, Response{
		Code:      code,
		Message:   message,
		Data:      nil,
		Timestamp: time.Now().UnixMilli(),
	})
}

// BadRequest 400错误
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

// Unauthorized 401错误
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

// Forbidden 403错误
func Forbidden(c *gin.Context, message string) {
	Error(c, 403, message)
}

// NotFound 404错误
func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// InternalError 500错误
func InternalError(c *gin.Context, message string) {
	Error(c, 500, message)
}
