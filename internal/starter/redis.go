package starter

import (
	"github.com/cuihairu/salon/internal/config"
	"github.com/gin-contrib/sessions/redis"
	"runtime"
)

func NewRedis(redisConfig *config.RedisConfig) (redis.Store, error) {
	if redisConfig == nil {
		redisConfig = &config.RedisConfig{
			NumConn:  runtime.NumCPU(),
			Address:  "127.0.0.1:6379",
			Db:       0,
			Secret:   "somerandomsecret",
			Password: "",
		}
	}
	return redis.NewStore(redisConfig.NumConn, "tcp", redisConfig.Address, redisConfig.Password, []byte(redisConfig.Secret))
}
