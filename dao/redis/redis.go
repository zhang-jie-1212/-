package redis

import (
	"bluebell/settings"
	"fmt"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConfig) error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: "", //取不到默认空
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		zap.L().Error("connect redis failed", zap.Error(err))
	}
	return err
}
func Close() {
	rdb.Close()
}
