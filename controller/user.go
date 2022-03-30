package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"strings"
)

//定义函数，去掉validator翻译成中文后的错误信息中的：后端定义的结构体字段名
func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	//遍历传过来的Map 错误，err代表value res的键为：找到key（field)中“."的索引，从它后一个开始进行切片）
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

func SignUpHandler(c *gin.Context) {
	//1.获取请求参数以及请求参数校验
	//前后端分离项目都使用JSON数据格式传递，所以用shouldJSONBind()函数来进行解析请求参数
	q := new(models.ParamSignUp)
	//参数出错，获取失败
	if err := c.ShouldBindJSON(q); err != nil {
		zap.L().Error("signup with invalid param", zap.Error(err))
		//首先将err转为validator.ValidationErrors类型，如果OK，则是validator校验产生的错误，进行翻译；如果不是，不能翻译
		errs, ok := err.(validator.ValidationErrors)
		//不是翻译错误，则是参数校验失败，无效参数
		if !ok {
			ResponseError(c, CodeInvalidParam)
			return
		}
		//如果是，对validator.ValidationErrors进行翻译后返回,调用errs.Translate，传入trans翻译器进行翻译
		//校验错误,返回的msg是校验错误信息
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务逻辑处理
	if err := logic.SignUp(q); err != nil {
		zap.L().Error("logic.SignUp failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeExistUser)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	//3.返回响应
	ResponseSuccess(c, "注册成功！")
}
func LoginHandler(c *gin.Context) {
	//1.获取请求参数和参数校验（校验密码不为空）
	q := new(models.LoginParams)
	if err := c.ShouldBindJSON(q); err != nil {
		//请求参数有误，写入error
		zap.L().Error("login with invalid param", zap.Error(err))
		//转为validator.ValidationErrors类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			//参数出错
			ResponseError(c, CodeInvalidParam)
			return
		}
		//校验出错，翻译出错信息
		ResponseErrorWithMsg(c, CodeInvalidParam, removeTopStruct(errs.Translate(trans)))
		return
	}
	//2.业务逻辑:判断数据库中是否有这个用户，有的话，判断密码是否正确
	token, err := logic.Login(q)
	if err != nil {
		//记录错误
		zap.L().Error("logic.login failed", zap.String("username", q.Username), zap.Error(err))
		//参数校验错误
		if token == "" {
			if errors.Is(err, mysql.ErrorPassword) {
				ResponseError(c, CodeInvalidPassword)
				return
			}
			ResponseError(c, CodeNotExistUser)
			return
		}
		//token生成错误
		ResponseError(c, CodeServerBusy)
	}
	//3.返回success响应和token
	ResponseSuccess(c, token)
}

func PongHandler(c *gin.Context) {

}
