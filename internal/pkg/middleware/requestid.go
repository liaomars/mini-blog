package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/liaomars/mini-blog/internal/pkg/know"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头中是否有 `X-Request-ID`，如果有则复用，没有则新建
		requestId := c.Request.Header.Get(know.XRequestIDKey)

		if requestId == "" {
			requestId = uuid.New().String()
		}

		// 将 RequestID 保存在 gin.Context 中，方便后边程序使用
		c.Set(know.XRequestIDKey, requestId)

		// 将 RequestID 保存在 HTTP 返回头中，Header 的键为 `X-Request-ID`
		c.Writer.Header().Set(know.XRequestIDKey, requestId)
		c.Next()
	}
}
