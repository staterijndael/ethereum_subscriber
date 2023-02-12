package greedy_redis_repository

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const redisHostBlockRepository = "0.0.0.0"
const redisPortBlockRepository = "6379"
const redisPasswordBlockRepository = ""

func TestBlockRepository_SetMaxCurrentBlock(t *testing.T) {
	ctx := context.TODO()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHostBlockRepository + ":" + redisPortBlockRepository,
		Password: redisPasswordBlockRepository,
		DB:       0,
	})

	// first testcase
	blockRepository := NewBlockRepository(redisClient, 10*time.Second)
	err := blockRepository.SetMaxCurrentBlock(ctx, 1)
	assert.NoError(t, err)

	err = blockRepository.SetMaxCurrentBlock(ctx, 5)

	currentBlock, err := blockRepository.GetCurrentBlock(ctx)
	assert.NoError(t, err)

	assert.Equal(t, currentBlock, uint64(5))

	// second testcase

	blockRepository = NewBlockRepository(redisClient, 10*time.Second)
	err = blockRepository.SetMaxCurrentBlock(ctx, 5)
	assert.NoError(t, err)

	currentBlock, err = blockRepository.GetCurrentBlock(ctx)
	assert.NoError(t, err)

	assert.Equal(t, currentBlock, uint64(5))

	err = blockRepository.SetMaxCurrentBlock(ctx, 2)
	assert.NoError(t, err)

	currentBlock, err = blockRepository.GetCurrentBlock(ctx)
	assert.NoError(t, err)

	assert.Equal(t, currentBlock, uint64(5))
}

func TestBlockRepository_GetCurrentBlock(t *testing.T) {
	ctx := context.TODO()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHostBlockRepository + ":" + redisPortBlockRepository,
		Password: redisPasswordBlockRepository,
		DB:       0,
	})

	// first testcase
	blockRepository := NewBlockRepository(redisClient, 10*time.Second)
	err := blockRepository.SetMaxCurrentBlock(ctx, 7)
	assert.NoError(t, err)

	currentBlock, err := blockRepository.GetCurrentBlock(ctx)
	assert.NoError(t, err)

	assert.Equal(t, currentBlock, uint64(7))
}
