package mini_blog

import (
	"github.com/gin-gonic/gin"
	"github.com/liaomars/mini-blog/internal/mini-blog/controller/v1/user"
	"github.com/liaomars/mini-blog/internal/mini-blog/store"
	"github.com/liaomars/mini-blog/internal/pkg/core"
	"github.com/liaomars/mini-blog/internal/pkg/errno"
	"github.com/liaomars/mini-blog/internal/pkg/log"
)

func installRouter(g *gin.Engine) error {
	// 注册404 处理handler
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	// 注册一个心跳检查访问路由
	g.GET("/healthz", func(c *gin.Context) {
		log.C(c).Infow("Healthz function called")
		core.WriteResponse(c, nil, map[string]string{"status": "ok"})
	})

	uc := user.New(store.S)
	v1 := g.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)
		}
	}
	return nil
}
