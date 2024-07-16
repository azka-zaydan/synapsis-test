package infras

import (
	"context"
	"fmt"

	"github.com/azka-zaydan/synapsis-test/configs"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type Redis struct {
	Client *redis.Client
}

// RedisNewClient create new instance of redis
func RedisNewClient(config *configs.Config) *Redis {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Cache.Redis.Primary.Host, config.Cache.Redis.Primary.Port),
		Password: config.DB.MySQL.Read.Password,
		DB:       config.Cache.Redis.Primary.DB,
	})

	_, err := client.Ping(context.TODO()).Result()
	if err != nil {
		panic(err)
	}
	log.Info().Msg("Redis Connected!")

	return &Redis{
		Client: client,
	}
}
