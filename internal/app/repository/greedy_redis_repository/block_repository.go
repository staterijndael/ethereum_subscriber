package greedy_redis_repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

// BlockRepository is a structure that holds information about the current block.
// It uses a Redis client to store the current block and has an expiration time for the stored value.
type BlockRepository struct {
	redis          *redis.Client
	expirationTime time.Duration
}

// NewBlockRepository returns a new instance of the BlockRepository with the specified Redis client and expiration time.
func NewBlockRepository(redis *redis.Client, expirationTime time.Duration) *BlockRepository {
	return &BlockRepository{
		redis:          redis,
		expirationTime: expirationTime,
	}
}

// SetMaxCurrentBlock updates the current block in the repository.
// If the new current block is higher than the existing current block, it sets the new value and returns nil.
// If the value does not exist, it sets the new value and returns nil.
// If there is an error, it returns the error.
func (r *BlockRepository) SetMaxCurrentBlock(ctx context.Context, newCurrentBlock uint64) error {
	// Get the current block from the repository
	currentBlock, err := r.redis.Get(ctx, getCurrentBlockKey()).Uint64()
	if err != nil {
		// If the current block does not exist, set the new current block
		if err.Error() == redisNilErrMsg {
			err = r.redis.Set(ctx, getCurrentBlockKey(), serializeCurrentBlockValue(newCurrentBlock), r.expirationTime).Err()
			if err != nil {
				return err
			}
			return nil
		}

		// Return the error if there was a problem retrieving the current block
		return err
	}

	// If the new current block is higher than the existing current block, set the new value
	if newCurrentBlock > currentBlock {
		err := r.redis.Set(ctx, getCurrentBlockKey(), serializeCurrentBlockValue(newCurrentBlock), r.expirationTime).Err()
		if err != nil {
			return err
		}
	}

	// Return nil to indicate success
	return nil
}

// GetCurrentBlock returns the current block from the repository.
// If there is an error, it returns 0 and the error.
func (r *BlockRepository) GetCurrentBlock(ctx context.Context) (uint64, error) {
	// Get the current block from the repository
	currentBlock, err := r.redis.Get(ctx, getCurrentBlockKey()).Uint64()
	if err != nil {
		// Return 0 and the error if there was a problem retrieving the current block
		return 0, err
	}
	// Return the current block and nil to indicate success
	return currentBlock, nil
}
