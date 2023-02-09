package mini_blog

import (
	"github.com/gin-gonic/gin"
	"github.com/liaomars/mini-blog/internal/mini-blog/controller/v1/user"
	"github.com/liaomars/mini-blog/internal/mini-blog/store"
	"github.com/liaomars/mini-blog/internal/pkg/core"
	"github.com/liaomars/mini-blog/internal/pkg/errno"
	"github.com/liaomars/mini-blog/internal/pkg/log"
	mw "github.com/liaomars/mini-blog/internal/pkg/middleware"
	"github.com/liaomars/mini-blog/pkg/auth"
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

	authz, err := auth.NewAuthz(store.S.DB())
	if err != nil {
		return nil
	}

	uc := user.New(store.S, authz)
	// 登录
	g.POST("/login", uc.Login)
	v1 := g.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)
			userv1.PUT(":name/change-password", uc.ChangePassword)
			userv1.Use(mw.Authn(), mw.Authz(authz))
			userv1.GET(":name", uc.GetUserInfo)
		}
	}
	return nil
}
