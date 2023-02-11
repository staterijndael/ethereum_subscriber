package sync_parser

import (
	"context"
	"errors"
	ethereum_jsonrpc "github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc"
	ethereum_jsonrpc_models "github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc/models"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"sync"
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

	transactions := make([]*models.Transaction, 0, txCount)

	var txPool = sync.Pool{
		New: func() interface{} {
			return &models.Transaction{}
		},
	}

	for i := uint64(currentBlockNumberResp.BlockNumber); i >= subscriber.SubscribeBlockNumber && (txCount > 0); i-- {
		blockNumber := i

		blockResp, err := p.ethereumJsonRPCClient.GetBlockByNumber(&ethereum_jsonrpc.GetBlockByNumberReq{
			BlockNumber: ethereum_jsonrpc_models.HexUint64(blockNumber),
			IsGetFullTx: true,
		})
		if err != nil {
			return nil, err
		}

		if blockNumber == uint64(currentBlockNumberResp.BlockNumber) {
			err = p.blockRepository.SetMaxCurrentBlock(ctx, blockNumber)
			if err != nil {
				return nil, err
			}
		}

		for j := len(blockResp.Block.Transactions) - 1; j >= 0 && txCount > 0; j-- {
			tx := blockResp.Block.Transactions[j]
			if tx.From == address || tx.To == address {
				transaction := txPool.Get().(*models.Transaction)
				*transaction = *models.ConvertJsonRPCTxToInternal(tx)
				transactions = append(transactions, transaction)

				txCount--
			}
		}
	}

	return transactions, nil
}
