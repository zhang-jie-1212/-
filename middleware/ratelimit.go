package middleware

//令桶限流策略
import (
	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"net/http"
	"time"
)

func RateLimitMiddleWare(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	//1.创建一个令牌桶,2秒钟填充一个，总容量为1
	bucket := ratelimit.NewBucket(fillInterval, cap)
	//2.从桶中取令牌
	return func(c *gin.Context) {
		//娶不到直接返回响应并不继续向下执行,取的个数为1，返回的是取到的令牌数目（或者从桶中删去的令牌数目）
		if bucket.TakeAvailable(1) == 0 {
			c.String(http.StatusOK, "rate limit....")
			c.Abort()
		}
		//否则向下执行
		c.Next()
	}
}
