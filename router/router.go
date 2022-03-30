package router

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"time"
)

func Setup(mode string) *gin.Engine {
	//gin框架的调试信息的输出位置，和logger一样，也可以自己设置，dev开发模式输出到终端，release发布模式不输出任何信息
	//若为发布模式，不输出gin的调试信息
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	//准册路由引擎
	r := gin.New()
	//引入request和recover中间件
	r.Use(logger.SetGinRequestLogger(), logger.SetGinRecoveryLogger(true), middleware.RateLimitMiddleWare(2*time.Second, 1))
	//引入限流中间件
	//r.Use(logger.SetGi  nRequestLogger(), logger.SetGinRecoveryLogger(true), middleware.RateLimitMiddleWare(2*time.Second, 1))
	v := r.Group("/air/v1")
	//用户点击注册
	v.POST("/signup", controller.SignUpHandler)
	v.GET("/login", controller.LoginHandler)
	v.GET("/post/:id", controller.PostDetailHandler)
	v.GET("/posts", controller.PostListHandler)
	v.GET("/postsort", controller.GetSortPostList)
	//每次对请求进行token验证
	v.Use(middleware.JWTTokenMiddelWare)
	{
		v.GET("/community", controller.CommunityHandler)
		v.GET("/community/:id", controller.CommunityDetailHandler)
		v.POST("/post", controller.CreatePostHandler)
		v.POST("/vote", controller.PostVoteHandler)
	}
	pprof.Register(r)
	return r
}
