package async_parser

import (
	"errors"
	ethereum_jsonrpc "github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc"
	ethereum_jsonrpc_models "github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc/models"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"math/big"
	"sync"
	"sync/atomic"
)

type EthereumJsonRPCClient interface {
	GetTxCount(req *ethereum_jsonrpc.GetTxCountReq) (*ethereum_jsonrpc.GetTxCountResp, error)
	GetBlockByNumber(req *ethereum_jsonrpc.GetBlockByNumberReq) (*ethereum_jsonrpc.GetBlockByNumberResp, error)
	GetBlockNumber() (*ethereum_jsonrpc.GetBlockNumberResp, error)
}

type SubscriberRepository interface {
	AddNewSubscriber(subscriber models.Subscriber) error
	GetSubscriberByAddress(address string) (models.Subscriber, error)
}

type BlockRepository interface {
	SetMaxCurrentBlock(newCurrentBlock uint64)
	GetCurrentBlock() uint64
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

func (p *Parser) GetCurrentBlock() uint64 {
	currentBlock := p.blockRepository.GetCurrentBlock()

	return currentBlock
}

func (p *Parser) Subscribe(address string) error {
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

	err = p.subscriberRepository.AddNewSubscriber(models.Subscriber{
		Address:              address,
		SubscribeBlockNumber: uint64(blockNumberResp.BlockNumber),
		SubscribeTxCount:     uint64(txCountResp.Nonce),
	})

	return err
}

func (p *Parser) GetTransactions(address string) ([]*models.Transaction, error) {
	subscriber, err := p.subscriberRepository.GetSubscriberByAddress(address)
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

			bigBlockNumber := new(big.Int).SetUint64(blockNumber)
			blockResp, err := p.ethereumJsonRPCClient.GetBlockByNumber(&ethereum_jsonrpc.GetBlockByNumberReq{
				BlockNumber: ethereum_jsonrpc_models.HexBigInt(*bigBlockNumber),
				IsGetFullTx: true,
			})
			if err != nil {
				errChan <- err
				return
			}

			if blockNumber == uint64(currentBlockNumberResp.BlockNumber) {
				p.blockRepository.SetMaxCurrentBlock(blockNumber)
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
