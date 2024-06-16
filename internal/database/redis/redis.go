package redis

import (
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client      *redis.Client
	environment string
}

func NewClient(redisURL string, environment string) (*RedisClient, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	return &RedisClient{
		client:      redis.NewClient(opts),
		environment: environment,
	}, nil
}
