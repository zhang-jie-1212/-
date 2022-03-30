package main

import (
	"bluebell/controller"
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/logger"
	"bluebell/pkg/snowflake"
	"bluebell/router"
	"bluebell/settings"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//1.加载配置文件到viper中 settings文件夹下的settings.go实现
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed,err:%v\n", err)
		return
	}
	//2.初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		fmt.Printf("init  logger failed,err:%v\n", err)
		return
	}
	//最后将缓冲区的日志追加到logger中
	defer zap.L().Sync()
	//3.初始化mysql链接
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		fmt.Printf("init  mysql failed,err:%v\n", err)
		return
	}
	//最后关闭db连接
	defer mysql.Close()
	//4.初始化redis链接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init  redis failed,err:%v\n", err)
		return
	}
	defer redis.Close()
	//5.初始化ID生成器
	if err := snowflake.Init(settings.Conf.AppConfig.StartTime, settings.Conf.AppConfig.MachineId); err != nil {
		fmt.Printf("init  sonyflake failed,err:%v\n", err)
		return
	}
	//6.初始化gin框架内置的校验器使用的翻译器,在shouldBindJSON校验错误时翻译错误
	if err := controller.InitTrans("zh"); err != nil {
		fmt.Printf("Init validator trans failed,err:%v\n", err)
		return
	}
	//7.创建路由引擎
	r := router.Setup(settings.Conf.Mode)
	//8.连接坚监听（优雅关机）
	//监听链接
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			zap.L().Error("listen failed", zap.Error(err))
		}
	}()
	//注册管道接收退出信号
	quit := make(chan os.Signal, 1)
	//监听设置的信号并转发给管道syscall.SIGTERM:kill默认发送的信号；syscall.SIGINT：kill -2发送（ctrl+c);syscall.SIGKILL:kill-9
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	//从管道读取信号，没有信号阻塞
	<-quit
	//有信号执行下面的代码：优雅关机
	zap.L().Info("Shutdown Server...")
	//创建一个超时5秒的context(监听port下的请求）,返回一个封装的context和一个cancel方法，此方法检测超时就关闭连接
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel() //5秒之后还没处理完请求就退出。多个goroutinr可以同时调用concel(),第一次调用之后看，对cancel之后调用不起作用（一个goroutine知调用一次即可）
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server shutdown:", zap.Error(err))
	}
	zap.L().Info("server exiting")
}
