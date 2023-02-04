package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/liaomars/mini-blog/internal/pkg/core"
	"github.com/liaomars/mini-blog/internal/pkg/errno"
	"github.com/liaomars/mini-blog/internal/pkg/know"
	"github.com/liaomars/mini-blog/pkg/token"
)

func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {
		username, err := token.ParseRequest(c)

		if err != nil {
			core.WriteResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}

		// 将 username 保存在 gin.Context 中，方便后边程序使用
		c.Set(know.XUsernameKey, username)

		c.Next()
	}
}
