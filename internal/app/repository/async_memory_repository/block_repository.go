package async_memory_repository

import (
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

func (r *BlockRepository) SetMaxCurrentBlock(newCurrentBlock uint64) {
	r.currentBlockMx.RLock()

	if newCurrentBlock > r.currentBlock {
		r.currentBlockMx.RUnlock()

		r.currentBlockMx.Lock()
		r.currentBlock = newCurrentBlock
		r.currentBlockMx.Unlock()

		return
	}

	r.currentBlockMx.RUnlock()
}

func (r *BlockRepository) GetCurrentBlock() uint64 {
	r.currentBlockMx.RLock()
	currentBlock := r.currentBlock
	r.currentBlockMx.RUnlock()

	return currentBlock
}
