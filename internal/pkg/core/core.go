package core

import (
	"github.com/gin-gonic/gin"
	"github.com/liaomars/mini-blog/internal/pkg/errno"
	"net/http"
)

type ErrResponse struct {
	// Code 指定了业务错误码.
	Code string `json:"code"`

	// Message 指定了错误内容.
	Message string `json:"message"`
}

func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		hcode, code, message := errno.Decode(err)
		c.JSON(hcode, ErrResponse{
			Code:    code,
			Message: message,
		})
		return
	}
	c.JSON(http.StatusOK, data)
}
