package async_parser

import (
	"context"
	"errors"
	ethereum_jsonrpc "github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc"
	ethereum_jsonrpc_models "github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc/models"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"sync"
	"sync/atomic"
)

type EthereumJsonRPCClient interface {
	GetTxCount(req *ethereum_jsonrpc.GetTxCountReq) (*ethereum_jsonrpc.GetTxCountResp, error)
	GetBlockByNumber(req *ethereum_jsonrpc.GetBlockByNumberReq) (*ethereum_jsonrpc.GetBlockByNumberResp, error)
	GetBlockNumber() (*ethereum_jsonrpc.GetBlockNumberResp, error)
}

type SubscriberRepository interface {
	AddNewSubscriber(ctx context.Context, subscriber models.Subscriber) error
	GetSubscriberByAddress(ctx context.Context, address string) (models.Subscriber, error)
}

type BlockRepository interface {
	SetMaxCurrentBlock(ctx context.Context, newCurrentBlock uint64) error
	GetCurrentBlock(ctx context.Context) (uint64, error)
}

type Parser struct {
	ethereumJsonRPCClient EthereumJsonRPCClient
	subscriberRepository  SubscriberRepository
	blockRepository       BlockRepository
}

func NewParser(ethereumJsonRPCClient EthereumJsonRPCClient, subscriberRepository SubscriberRepository, blockRepository BlockRepository) *Parser {
	return &Parser{
		ethereumJsonRPCClient: ethereumJsonRPCClient,
		subscriberRepository:  subscriberRepository,
		blockRepository:       blockRepository,
	}
}

func (p *Parser) GetCurrentBlock(ctx context.Context) (uint64, error) {
	currentBlock, err := p.blockRepository.GetCurrentBlock(ctx)

	return currentBlock, err
}

func (p *Parser) Subscribe(ctx context.Context, address string) error {
	blockNumberResp, err := p.ethereumJsonRPCClient.GetBlockNumber()
	if err != nil {
		return errors.New("error getting current block number cause: " + err.Error())
	}

	txCountResp, err := p.ethereumJsonRPCClient.GetTxCount(&ethereum_jsonrpc.GetTxCountReq{
		Address:  address,
		EndBlock: blockNumberResp.BlockNumber,
	})
	if err != nil {
		return err
	}

	err = p.subscriberRepository.AddNewSubscriber(ctx, models.Subscriber{
		Address:              address,
		SubscribeBlockNumber: uint64(blockNumberResp.BlockNumber),
		SubscribeTxCount:     uint64(txCountResp.Nonce),
	})

	return err
}

// GetTransactions uses for a getting a full list of inbound or outbounds transactions by user since subscription.
// Due to simple transactions in Ethereum Network are not indexed, we cannot get it through simple JSONRPC methods.
// This method uses heuristic approach without having to process the entire chain.
// Algorithm:
// 1. Get current block number in ethereum network to set up point to which we must search and block number in a moment of last parsed block (currentBlockNumber, lastBlockNumber).
// 2. Get count of inbound or outbound transactions in current moment and in a moment of last parsed txCount related to user (currentTxCount, lastTxCount)
// 3. txCount (count of transactions that should be handled in a future starting from last parsed time) = currentTxCount - lastTxCount
// 4. Iterate in range currentBlockNumber..lastBlockNumber and get all transactions for every block where from==address or to==address until txCount > 0
// NOTE 1*: method starting loop from currentBlockNumber to lastBlockNumber, because we have a gap from currentBlockNumber,
// because we could start our subscription not from start transaction for block
// Example: transaction count in every block = 3 (blockSize), blockCount = 10, currentBlockNumber = 10,
// lastBlockNumber = 1 (1-indexed block numbers in this example), lastTxCount = 2, currentTxCount = 28
// txCount = 26, if we iterate from 1 block and 1 transaction then we will catch transaction 1..27 instead of 3..29,
// so we need to start from transaction 3 (subscriptionTxCount + 1), but if we had a lastTxCount > blockSize
// then we need to find our starting transaction first, it is a little annoying
// it is much easier to iterate from last block because we will already have reversed order, and we don`t need to count transactions
// NOTE 2*: as opposed to SyncParser (Releasing Approach) approach Greedy implies saving already parsed transaction,
// we can afford it due to Ethereum transaction immutability promise (https://ethereum.org/en/developers/docs/transactions/#:~:text=As%20time%20passes,billions%20of%20dollars.)
// this method might require distributed storage (like Redis or PostgreSQL) instead of default memory storage
// but it might significantly increase performance due to keeping already parsed transactions.
// This approach aimed at long-term program execution with long lifetime.
// NOTE 3* this approach by default does not guarantees order of transactions and it might require extra sorting for transactions
func (p *Parser) GetTransactions(ctx context.Context, address string) ([]*models.Transaction, error) {
	subscriber, err := p.subscriberRepository.GetSubscriberByAddress(ctx, address)
	if err != nil {
		return nil, err
	}

	currentBlockNumberResp, err := p.ethereumJsonRPCClient.GetBlockNumber()
	if err != nil {
		return nil, err
	}

	currentTxCountResp, err := p.ethereumJsonRPCClient.GetTxCount(&ethereum_jsonrpc.GetTxCountReq{
		Address:  address,
		EndBlock: currentBlockNumberResp.BlockNumber,
	})
	if err != nil {
		return nil, err
	}

	txCount := uint64(currentTxCountResp.Nonce) - subscriber.SubscribeTxCount

	var addressTxCountAtomic uint64

	atomic.StoreUint64(&addressTxCountAtomic, txCount)

	transactions := make([]*models.Transaction, 0, addressTxCountAtomic)
	transactionMx := sync.Mutex{}

	var txPool = sync.Pool{
		New: func() interface{} {
			return &models.Transaction{}
		},
	}

	wg := sync.WaitGroup{}

	errChan := make(chan error, txCount+1)

	for i := uint64(currentBlockNumberResp.BlockNumber); i >= subscriber.SubscribeBlockNumber; i-- {
		wg.Add(1)
		go func(blockNumber uint64) {
			defer wg.Done()

			blockResp, err := p.ethereumJsonRPCClient.GetBlockByNumber(&ethereum_jsonrpc.GetBlockByNumberReq{
				BlockNumber: ethereum_jsonrpc_models.HexUint64(blockNumber),
				IsGetFullTx: true,
			})
			if err != nil {
				errChan <- err
				return
			}

			if blockNumber == uint64(currentBlockNumberResp.BlockNumber) {
				err = p.blockRepository.SetMaxCurrentBlock(ctx, blockNumber)
				if err != nil {
					errChan <- err
					return
				}
			}

			for j := len(blockResp.Block.Transactions) - 1; j >= 0 && atomic.LoadUint64(&addressTxCountAtomic) > 0; j-- {
				tx := blockResp.Block.Transactions[j]
				if tx.From == address || tx.To == address {
					transactionMx.Lock()
					transaction := txPool.Get().(*models.Transaction)
					*transaction = *models.ConvertJsonRPCTxToInternal(tx)
					transactions = append(transactions, transaction)
					transactionMx.Unlock()

					atomic.AddUint64(&addressTxCountAtomic, ^uint64(0))
				}
			}
		}(i)
	}

	wg.Wait()

	var errorMsg string

	var done bool
	for !done {
		select {
		case err := <-errChan:
			if err != nil {
				errorMsg += "Error: " + err.Error() + "\n"
			}
		default:
			done = true
			break
		}
	}

	if errorMsg != "" {
		return nil, errors.New(errorMsg)
	}

	return transactions, nil
}
