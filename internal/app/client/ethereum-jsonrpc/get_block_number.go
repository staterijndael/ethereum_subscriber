package ethereum_jsonrpc

import (
	"encoding/json"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc/models"
)

// getBlockNumberRPCName is a constant that stores the name of the JSON-RPC method "eth_blockNumber".
const getBlockNumberRPCName = "eth_blockNumber"

// GetBlockNumberResp represents the response of the JSON-RPC method "eth_blockNumber".
type GetBlockNumberResp struct {
	// BlockNumber is a hexadecimal representation of the latest block number on the blockchain.
	BlockNumber models.HexUint64
}

// GetBlockNumber makes a JSON-RPC call to the Ethereum node to retrieve the latest block number on the blockchain.
// It returns the block number in hexadecimal representation and an error if the call fails.
func (c *Client) GetBlockNumber() (*GetBlockNumberResp, error) {
	// Send a JSON-RPC request to the Ethereum node with method name "eth_blockNumber" and no parameters.
	rawReqResp, err := c.sendJSONRPCRequest(getBlockNumberRPCName, []interface{}{})
	if err != nil {
		// If the JSON-RPC request fails, return the error.
		return nil, err
	}
	// Unmarshal the response from the Ethereum node into a GetBlockNumberResp struct.
	var getTxClientResp GetBlockNumberResp
	err = json.Unmarshal(rawReqResp, &getTxClientResp.BlockNumber)
	if err != nil {
		// If the response from the Ethereum node cannot be unmarshaled, return the error.
		return nil, err
	}

	// Return the GetBlockNumberResp struct and a nil error.
	return &getTxClientResp, nil
}
