package infras

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/farganamar/evv-service/configs"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
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
	// Check if Redis is properly configured
	if config.Cache.Redis.Primary.Host == "" {
		log.Warn().Msg("Redis host not configured, skipping Redis initialization")
		return nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%s", config.Cache.Redis.Primary.Host, config.Cache.Redis.Primary.Port),
		Password:   config.Cache.Redis.Primary.Password,
		MaxRetries: config.Cache.Redis.Primary.RetryCount,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Info().Msg("Redis connection established")
			return nil
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Error().Err(err).Msg("Redis connection failed")
		// Consider returning nil here if Redis is optional
		panic(err)
	} else {
		log.Info().Msg("Redis connection successful")
	}

	return client
}
