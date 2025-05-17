package infras

import (
	"context"
	"fmt"

	"github.com/farganamar/evv-service/configs"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Prefix      string
	RedisClient *redis.Client
}

// ProvideRedis is the provider for Redis.
func ProvideRedis(config *configs.Config) *Redis {
	redisClient := RedisNewClient(*config)
	return &Redis{
		Prefix:      config.Cache.Redis.Primary.Prefix,
		RedisClient: redisClient,
	}
}

// RedisNewClient create new instance of redis
func RedisNewClient(config configs.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", config.Cache.Redis.Primary.Host, config.Cache.Redis.Primary.Port),
		Password:   config.Cache.Redis.Primary.Password,
		MaxRetries: config.Cache.Redis.Primary.RetryCount,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return client
}
