package greedy_redis_repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type BlockRepository struct {
	redis          *redis.Client
	expirationTime time.Duration
}

func NewBlockRepository(redis *redis.Client, expirationTime time.Duration) *BlockRepository {
	return &BlockRepository{
		redis:          redis,
		expirationTime: expirationTime,
	}
}

func (r *BlockRepository) SetMaxCurrentBlock(ctx context.Context, newCurrentBlock uint64) error {
	currentBlock, err := r.redis.Get(ctx, getCurrentBlockKey()).Uint64()
	if err != nil {
		if err.Error() == redisNilErrMsg {
			err = r.redis.Set(ctx, getCurrentBlockKey(), serializeCurrentBlockValue(newCurrentBlock), r.expirationTime).Err()
			if err != nil {
				return err
			}

			return nil
		}

		return err
	}

	if newCurrentBlock > currentBlock {
		err := r.redis.Set(ctx, getCurrentBlockKey(), serializeCurrentBlockValue(newCurrentBlock), r.expirationTime).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *BlockRepository) GetCurrentBlock(ctx context.Context) (uint64, error) {
	currentBlock, err := r.redis.Get(ctx, getCurrentBlockKey()).Uint64()
	if err != nil {
		return 0, err
	}

	return currentBlock, nil
}
