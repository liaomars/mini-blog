package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/liaomars/mini-blog/internal/pkg/core"
	"github.com/liaomars/mini-blog/internal/pkg/errno"
	"github.com/liaomars/mini-blog/internal/pkg/know"
	"github.com/liaomars/mini-blog/internal/pkg/log"
)

type Auther interface {
	Authorize(sub, obj, act string) (bool, error)
}

// Authz 在gin中间件，对请求进行授权检查
func Authz(a Auther) gin.HandlerFunc {
	return func(c *gin.Context) {
		sub := c.GetString(know.XUsernameKey)
		obj := c.Request.URL.Path
		act := c.Request.Method

		log.Debugw("Build authorize context", "sub", sub, "obj", obj, "act", act)
		if allowed, _ := a.Authorize(sub, obj, act); !allowed {
			core.WriteResponse(c, errno.ErrUnauthorized, nil)
			c.Abort()
			return
		}
	}
}
