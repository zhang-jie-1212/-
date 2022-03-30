package logger

import (
	"bluebell/settings"
	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var logger *zap.Logger

//zap日志库创建日志文件
//设置core的文件打开路径并设置切割
func setLogWriter(filename string, max_size, max_age, max_backup int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    max_size,
		MaxAge:     max_age,
		MaxBackups: max_backup,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

//设置编码格式（JSON）并设置日期展示方式
func setEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
func Init(cfg *settings.LogConfig, mode string) error {
	//1、设置core的文件打开路径并设置切割
	writeSyncer := setLogWriter(
		cfg.Filename,
		cfg.MaxSize,
		cfg.MaxAge,
		cfg.MaxBackups)
	//2、设置编码格式（JSON）并设置日期展示方式
	encoder := setEncoder()
	//3、设置日志级别
	//type zapcore.Level int8：new初始化一个zapcore.Level类型的指针
	var l = new(zapcore.Level)
	//反序列化：将字节序列转换为l对象的过程，转换后为DebugLevel
	err := l.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		return err
	}
	//4、初始化core：两种：一种将信息输入到日志文件，一种输入到终端；取决于config.ymal:APP:mode
	//dev:开发模式：将信息既输入到终端，又输入到日志库中
	var core zapcore.Core
	if mode == "dev" {
		//输入到终端的编码
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		//zapcore.NewTee创建多个输出位置的core 其中输出到终端店额输出位置表示为：zapcore.Lock(os.Stdout)
		core = zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), l),
			zapcore.NewCore(encoder, writeSyncer, l),
		)
	} else {
		core = zapcore.NewCore(encoder, writeSyncer, l)
	}
	//5、zap.new(core,option...)创建日志，设置显示调用信息
	lg := zap.New(core, zap.AddCaller())
	//6.替换zap全局库中的logger,此后从别的库中取logger就用：zap.L(),防止logger.logger调用麻烦
	zap.ReplaceGlobals(lg)
	return nil
}

//路由中间件：像日志中写入请求信息
func SetGinRequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)
		//将请求写入logger;以及响应状态（响应吗
		zap.L().Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			//错误信息ErrorTypePrivate ErrorType = 1 << 0；表示一个私有错误
			//Errors是附加到使用此上下文的所有处理程序/中间件的错误列表。
			//ByType返回对gin.ErrorTypePrivate字节进行筛选的只读副本。
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// SetGinRecoveryLogger:程序运行中出现任何panic(error)，将err信息写入日志
func SetGinRecoveryLogger(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
