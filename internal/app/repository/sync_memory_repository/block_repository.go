package sync_memory_repository

import "context"

type BlockRepository struct {
	currentBlock uint64
}

func NewBlockRepository() *BlockRepository {
	return &BlockRepository{
		currentBlock: 0,
	}
}

func (r *BlockRepository) SetMaxCurrentBlock(ctx context.Context, newCurrentBlock uint64) error {
	if newCurrentBlock > r.currentBlock {
		r.currentBlock = newCurrentBlock
	}

	return nil
}

func (r *BlockRepository) GetCurrentBlock(ctx context.Context) (uint64, error) {
	currentBlock := r.currentBlock

	return currentBlock, nil
}
