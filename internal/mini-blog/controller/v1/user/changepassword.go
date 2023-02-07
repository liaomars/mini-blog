package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/liaomars/mini-blog/internal/pkg/core"
	"github.com/liaomars/mini-blog/internal/pkg/errno"
	"github.com/liaomars/mini-blog/internal/pkg/log"
	v1 "github.com/liaomars/mini-blog/pkg/api/miniblog/v1"
)

func (ctrl *UserController) ChangePassword(c *gin.Context) {
	log.C(c).Infow("change password function called")

	var r v1.ChangePasswordRequest

	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}

	//if err := ctrl.b.Users().Login(c, &r); err != nil {
	//	core.WriteResponse(c, err, nil)
	//	return
	//}
	err := ctrl.b.Users().ChangePassword(c, c.Param("name"), &r)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}
