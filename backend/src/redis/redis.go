package cache

import (
	"github.com/jj-style/chain-react/src/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(rconf *config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     rconf.Address,
		Password: rconf.Password,
		DB:       rconf.Db,
	})
}
