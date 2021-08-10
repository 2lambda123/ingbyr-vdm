/*
 @Author: ingbyr
 @Description:
*/

package r

import (
	"github.com/gin-gonic/gin"
	"github.com/ingbyr/vdm/pkg/e"
	"net/http"
)

type Response struct {
	Code uint        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func OK(c *gin.Context, data interface{}) {
	R(c, http.StatusOK, e.Ok, data)
}

func Failed(c *gin.Context, errorCode uint) {
	R(c, http.StatusOK, errorCode, "")
}

func R(c *gin.Context, httpCode int, code uint, data interface{}) {
	c.JSON(httpCode, Response{
		Code: code,
		Msg:  e.ToMsg(code),
		Data: data,
	})
}
