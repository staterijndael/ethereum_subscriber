package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, options *redis.Options) (*redis.Client, error) {
	rdb := redis.NewClient(options)

	cmd := rdb.Ping(ctx)
	err := cmd.Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
