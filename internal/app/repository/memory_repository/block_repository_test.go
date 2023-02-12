package memory_repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockRepository_SetMaxCurrentBlock(t *testing.T) {
	ctx := context.TODO()

	// first testcase
	blockRepository := NewBlockRepository()
	err := blockRepository.SetMaxCurrentBlock(ctx, 1)
	assert.NoError(t, err)

	err = blockRepository.SetMaxCurrentBlock(ctx, 5)

	currentBlock, err := blockRepository.GetCurrentBlock(ctx)
	assert.NoError(t, err)

	assert.Equal(t, currentBlock, uint64(5))

	// second testcase

	blockRepository = NewBlockRepository()
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

	blockRepository := NewBlockRepository()
	err := blockRepository.SetMaxCurrentBlock(ctx, 7)
	assert.NoError(t, err)

	currentBlock, err := blockRepository.GetCurrentBlock(ctx)
	assert.NoError(t, err)

	assert.Equal(t, currentBlock, uint64(7))
}
