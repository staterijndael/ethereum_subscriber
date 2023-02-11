package models

import (
	ethereum_jsonrpc "github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc"
	"math/big"
)

// swagger:model Transaction
type Transaction struct {
	BlockHash        string  `json:"blockHash"`
	BlockNumber      uint64  `json:"blockNumber"`
	From             string  `json:"from"`
	Gas              big.Int `json:"gas"`
	GasPrice         big.Int `json:"gasPrice"`
	Hash             string  `json:"hash"`
	Input            string  `json:"input"`
	Nonce            uint64  `json:"nonce"`
	To               string  `json:"to"`
	TransactionIndex uint64  `json:"transactionIndex"`
	Value            big.Int `json:"value"`
	V                big.Int `json:"v"`
	R                big.Int `json:"r"`
	S                big.Int `json:"s"`
}

func ConvertJsonRPCTxToInternal(tx *ethereum_jsonrpc.Transaction) *Transaction {
	if tx == nil {
		return nil
	}

	return &Transaction{
		BlockHash:        tx.BlockHash,
		BlockNumber:      uint64(tx.BlockNumber),
		From:             tx.From,
		Gas:              big.Int(tx.Gas),
		GasPrice:         big.Int(tx.GasPrice),
		Hash:             tx.Hash,
		Input:            tx.Input,
		Nonce:            uint64(tx.Nonce),
		To:               tx.To,
		TransactionIndex: uint64(tx.TransactionIndex),
		Value:            big.Int(tx.Value),
		V:                big.Int(tx.V),
		R:                big.Int(tx.R),
		S:                big.Int(tx.S),
	}
}

func ReverseTransactionsByLink(txs []*Transaction) {
	left := 0
	right := len(txs) - 1

	for left < right {
		tmp := txs[left]
		txs[left] = txs[right]
		txs[right] = tmp

		left++
		right--
	}
}

func ReverseTransactionsCopy(txs []*Transaction) []*Transaction {
	if len(txs) == 0 {
		return []*Transaction{}
	}
	newTxs := make([]*Transaction, len(txs))

	left := 0
	right := len(txs) - 1

	for left < right {
		copyLeft := *txs[left]
		copyRight := *txs[right]

		newTxs[left] = &copyRight
		newTxs[right] = &copyLeft

		left++
		right--
	}

	if len(txs)%2 != 0 {
		tmp := *txs[len(txs)/2]
		newTxs[len(txs)/2] = &tmp
	}

	return newTxs
}
