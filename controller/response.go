package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//将返回给前端的数据封装成一个struct,使前端能够看得懂
type ResponseData struct {
	Code TypeCode
	Msg  interface{}
	Data interface{}
}

//根据code返回错误的响应信息和数据
func ResponseError(c *gin.Context, code TypeCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.GetMsg(),
		Data: nil,
	})
}

// ResponseSuccess 返回自定义的code和msg
func ResponseErrorWithMsg(c *gin.Context, code TypeCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

//ResponseSuccess返回正确的信息
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.GetMsg(),
		Data: data,
	})
}
