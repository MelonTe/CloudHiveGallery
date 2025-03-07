package common

import (
	"chg/internal/ecode"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 统一响应
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func BaseResponse(g *gin.Context, data interface{}, msg string, code int) {
	g.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}
func Success(g *gin.Context, data interface{}) {
	BaseResponse(g, data, "", 0)
}

// 失败响应
func Error(g *gin.Context, code int) {
	BaseResponse(g, nil, ecode.GetErrMsg(code), code)
}
