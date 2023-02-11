package greedy_memory_repository

import (
	"context"
	"sync"
)

type BlockRepository struct {
	currentBlock   uint64
	currentBlockMx sync.RWMutex
}

func NewBlockRepository() *BlockRepository {
	return &BlockRepository{
		currentBlock:   0,
		currentBlockMx: sync.RWMutex{},
	}
}

func (r *BlockRepository) SetMaxCurrentBlock(ctx context.Context, newCurrentBlock uint64) error {
	r.currentBlockMx.RLock()

	if newCurrentBlock > r.currentBlock {
		r.currentBlockMx.RUnlock()

		r.currentBlockMx.Lock()
		r.currentBlock = newCurrentBlock
		r.currentBlockMx.Unlock()

		return nil
	}

	r.currentBlockMx.RUnlock()

	return nil
}

func (r *BlockRepository) GetCurrentBlock(ctx context.Context) (uint64, error) {
	r.currentBlockMx.RLock()
	currentBlock := r.currentBlock
	r.currentBlockMx.RUnlock()

	return currentBlock, nil
}
