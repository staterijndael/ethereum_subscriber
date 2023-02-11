package models

import (
	ethereum_jsonrpc "github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc"
	"math/big"
)

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
