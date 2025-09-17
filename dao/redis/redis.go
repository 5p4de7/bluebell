package redis

import (
	"bluebell/settings"
	"fmt"
	"github.com/go-redis/redis"
)

var rdb *redis.Client

func Init(cfg *settings.RedisConf) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host, cfg.Port),
		DB: cfg.Db,
	})
	_, err = rdb.Ping().Result()
	return
}
func Close() {
	_ = rdb.Close()
}
