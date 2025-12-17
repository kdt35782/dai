package middleware

import (
	"github.com/gin-gonic/gin"
	"sm-medical/pkg/utils"
	"strings"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// 检查格式: Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Unauthorized(c, "请求头格式错误")
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Unauthorized(c, "Token无效或已过期")
			c.Abort()
			return
		}

		// 将用户信息存入context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

// CORS 跨域中间件
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Authorization, Accept, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
