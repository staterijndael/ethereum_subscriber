package ethereum_jsonrpc

import (
	"encoding/json"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc/models"
)

// Defining the constant getBlockByNumberRPCName which is the name of the JSON-RPC method to be called.
const getBlockByNumberRPCName = "eth_getBlockByNumber"

// GetBlockByNumberReq represents the request payload for the "eth_getBlockByNumber" JSON-RPC method
type GetBlockByNumberReq struct {
	// BlockNumber is the block number or block tag that specifies the block to be returned
	BlockNumber models.HexUint64

	// IsGetFullTx is a boolean value that specifies whether to return the full transaction objects or just the hashes
	IsGetFullTx bool
}

// GetBlockByNumberResp represents the response payload for the "eth_getBlockByNumber" JSON-RPC method
type GetBlockByNumberResp struct {
	// Block is a Block struct that contains information about the block
	Block Block
}

// Block represents the information about a block in the Ethereum blockchain
type Block struct {
	// Number is the block number
	Number models.HexBigInt `json:"number"`

	// Hash is the block hash
	Hash string `json:"hash"`

	// Transactions is an array of transaction objects that belong to the block
	Transactions []*Transaction `json:"transactions"`
}

// Transaction represents a transaction in the Ethereum blockchain
type Transaction struct {
	// BlockHash is the hash of the block that this transaction belongs to
	BlockHash string `json:"blockHash"`

	// BlockNumber is the number of the block that this transaction belongs to
	BlockNumber models.HexUint64 `json:"blockNumber"`

	// From is the address of the account that initiated the transaction
	From string `json:"from"`

	// Gas is the amount of gas used by the transaction
	Gas models.HexBigInt `json:"gas"`

	// GasPrice is the price of gas used by the transaction
	GasPrice models.HexBigInt `json:"gasPrice"`

	// Hash is the transaction hash
	Hash string `json:"hash"`

	// Input is the input data for the transaction
	Input string `json:"input"`

	// Nonce is the nonce of the account that initiated the transaction
	Nonce models.HexUint64 `json:"nonce"`

	// To is the address of the account that the transaction was sent to
	To string `json:"to"`

	// TransactionIndex is the index of the transaction in the block
	TransactionIndex models.HexUint64 `json:"transactionIndex"`

	// Value is the amount of Ether transferred in the transaction
	Value models.HexBigInt `json:"value"`

	// V is the Ethereum network protocol version
	V models.HexBigInt `json:"v"`

	// R is a component of the signature of the transaction
	R models.HexBigInt `json:"r"`

	// S is a component of the signature of the transaction
	S models.HexBigInt `json:"s"`
}

// GetBlockByNumber method retrieves a block from the Ethereum blockchain, using its block number as the identifier.
// The method takes a GetBlockByNumberReq struct as an input argument, which contains the block number and a boolean flag indicating whether or not
// to retrieve the full transaction details for each transaction within the block.
func (c *Client) GetBlockByNumber(req *GetBlockByNumberReq) (*GetBlockByNumberResp, error) {
	// send the JSON-RPC request to the Ethereum client
	rawReqResp, err := c.sendJSONRPCRequest(getBlockByNumberRPCName, []interface{}{req.BlockNumber, req.IsGetFullTx})
	if err != nil {
		// if there was an error, return it
		return nil, err
	}

	// parse the raw JSON-RPC response into a GetBlockByNumberResp struct
	var getBlockByNumberResp GetBlockByNumberResp
	err = json.Unmarshal(rawReqResp, &getBlockByNumberResp.Block)
	if err != nil {
		// if there was an error during parsing, return it
		return nil, err
	}

	// return the response
	return &getBlockByNumberResp, nil
}
