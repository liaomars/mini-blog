package user

import (
	"github.com/gin-gonic/gin"
	"github.com/liaomars/mini-blog/internal/pkg/core"
	"github.com/liaomars/mini-blog/internal/pkg/log"
)

func (ctrl *UserController) GetUserInfo(c *gin.Context) {
	log.C(c).Infow("getuserinfo function called")

	resp, err := ctrl.b.Users().GetUserInfo(c, c.Param("name"))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, resp)
}
