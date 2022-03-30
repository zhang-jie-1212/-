package middleware

import (
	"bluebell/controller"
	"bluebell/pkg/jwt"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strings"
)

func JWTTokenMiddelWare(c *gin.Context) {
	var err error
	//1.从请求投中取出Authorization字段，对其进行分离
	authHeader := c.Request.Header.Get("Authorization")
	//2.分离头部出token,authHeader是string类型，Bearer和Token按空格分割
	if authHeader == "" {
		//没有此字段，参数有误，直接返回响应
		err = errors.New("参数为空")
		zap.L().Error("invalid auto", zap.Error(err))
		controller.ResponseError(c, controller.CodeNeedLogin)
		c.Abort()
		return
	}
	//按空格分割得到token
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		controller.ResponseError(c, controller.CodeInvaildToken)
		c.Abort()
		return
	}
	//对得到的token进行解析
	mc, err := jwt.ParseTocken(parts[1])
	if err != nil {
		zap.L().Error("jwt parse token failed", zap.Error(err))
		controller.ResponseError(c, controller.CodeInvaildToken)
		c.Abort()
		return
	}
	//将userid传给后续请求，userIDs是用户传递过来的，下一步要和数据库中做对比
	//controller.ContextUserID这个全局常量为什么不直接放在moddleware包：上面引用了controller包中的响应码
	//后面controller包要取出这个用户传过来的USERID进行操作，若放在middleware包，则相当于循环引用，会报错
	c.Set(controller.ContextUserID, mc.UserID)
	c.Next()
}
