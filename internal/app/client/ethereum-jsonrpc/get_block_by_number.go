package ethereum_jsonrpc

import (
	"encoding/json"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc/models"
)

const getBlockByNumberRPCName = "eth_getBlockByNumber"

type GetBlockByNumberReq struct {
	BlockNumber models.HexBigInt
	IsGetFullTx bool
}

type GetBlockByNumberResp struct {
	Block Block
}

type Block struct {
	Number       models.HexBigInt `json:"number"`
	Hash         string           `json:"hash"`
	Transactions []*Transaction   `json:"transactions"`
}

type Transaction struct {
	BlockHash        string           `json:"blockHash"`
	BlockNumber      models.HexUint64 `json:"blockNumber"`
	From             string           `json:"from"`
	Gas              models.HexBigInt `json:"gas"`
	GasPrice         models.HexBigInt `json:"gasPrice"`
	Hash             string           `json:"hash"`
	Input            string           `json:"input"`
	Nonce            models.HexUint64 `json:"nonce"`
	To               string           `json:"to"`
	TransactionIndex models.HexUint64 `json:"transactionIndex"`
	Value            models.HexBigInt `json:"value"`
	V                models.HexBigInt `json:"v"`
	R                models.HexBigInt `json:"r"`
	S                models.HexBigInt `json:"s"`
}

func (c *Client) GetBlockByNumber(req *GetBlockByNumberReq) (*GetBlockByNumberResp, error) {
	rawReqResp, err := c.sendJSONRPCRequest(getBlockByNumberRPCName, []interface{}{req.BlockNumber, req.IsGetFullTx})
	if err != nil {
		return nil, err
	}

	var getBlockByNumberResp GetBlockByNumberResp
	err = json.Unmarshal(rawReqResp, &getBlockByNumberResp.Block)
	if err != nil {
		return nil, err
	}

	return &getBlockByNumberResp, nil
}
