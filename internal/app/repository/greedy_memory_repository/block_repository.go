package greedy_memory_repository

import (
	"context"
	"sync"
)

// BlockRepository is a struct that contains information about the current block in a blockchain.
type BlockRepository struct {
	// currentBlock is an uint64 representing the current block in the blockchain.
	currentBlock uint64
	// currentBlockMx is a sync.RWMutex used to protect the currentBlock from concurrent access.
	currentBlockMx sync.RWMutex
}

// NewBlockRepository is a constructor function for BlockRepository that returns a new instance of BlockRepository.
func NewBlockRepository() *BlockRepository {
	return &BlockRepository{
		currentBlock:   0,
		currentBlockMx: sync.RWMutex{},
	}
}

// SetMaxCurrentBlock updates the currentBlock in the BlockRepository if the newCurrentBlock is greater than the currentBlock.
func (r *BlockRepository) SetMaxCurrentBlock(ctx context.Context, newCurrentBlock uint64) error {
	// Acquire a read lock to prevent concurrent writes to the currentBlock.
	r.currentBlockMx.RLock()

	// Check if the newCurrentBlock is greater than the currentBlock.
	if newCurrentBlock > r.currentBlock {
		// If so, release the read lock and acquire a write lock.
		r.currentBlockMx.RUnlock()
		r.currentBlockMx.Lock()

		// Update the currentBlock with the new value.
		r.currentBlock = newCurrentBlock

		// Release the write lock.
		r.currentBlockMx.Unlock()

		// Return nil to indicate success.
		return nil
	}

	// If the newCurrentBlock is not greater than the currentBlock, release the read lock.
	r.currentBlockMx.RUnlock()

	// Return nil to indicate success.
	return nil
}

// GetCurrentBlock returns the currentBlock in the BlockRepository.
func (r *BlockRepository) GetCurrentBlock(ctx context.Context) (uint64, error) {
	// Acquire a read lock to prevent concurrent writes to the currentBlock.
	r.currentBlockMx.RLock()

	// Get the currentBlock value.
	currentBlock := r.currentBlock

	// Release the read lock.
	r.currentBlockMx.RUnlock()

	// Return the currentBlock and nil to indicate success.
	return currentBlock, nil
}
