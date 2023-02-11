package sync_memory_repository

type BlockRepository struct {
	currentBlock uint64
}

func NewBlockRepository() *BlockRepository {
	return &BlockRepository{
		currentBlock: 0,
	}
}

func (r *BlockRepository) SetMaxCurrentBlock(newCurrentBlock uint64) {
	if newCurrentBlock > r.currentBlock {
		r.currentBlock = newCurrentBlock
	}
}

func (r *BlockRepository) GetCurrentBlock() uint64 {
	currentBlock := r.currentBlock

	return currentBlock
}
