package cache

import (
	"context"
	"log"

	"github.com/jj-style/chain-react/src/config"
	"github.com/redis/go-redis/v9"
)

func NewRedis(rconf *config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     rconf.Address,
		Password: rconf.Password,
		DB:       rconf.Db,
	})
	if err := client.Ping(context.TODO()).Err(); err != nil {
		log.Fatal(err)
	}
	return client
}
