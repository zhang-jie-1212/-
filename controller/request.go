package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
)

const ContextUserID = "userID"

//对request进行操作
//eg:获取上一个中间件传过来的用户ID等
var ErrorUserNotLogin = errors.New("用户未登录")

func GetCurentUserID(c *gin.Context) (uid int64, err error) {
	//返回的是接口类型的userID
	userID, ok := c.Get(ContextUserID)
	if !ok {
		//用户未登录
		err = ErrorUserNotLogin
		return
	}
	uid, ok = userID.(int64)
	if !ok {
		err = ErrorUserNotLogin
		return
	}
	return
}

//获取参数中的帖子页面和数量
func GetPostLimit(c *gin.Context) (int64, int64) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")
	page, err := strconv.ParseInt(pageStr, 10, 64)
	if err != nil {
		page = 1
	}
	size, err := strconv.ParseInt(sizeStr, 10, 64)
	if err != nil {
		size = 10
	}
	return page, size
}
