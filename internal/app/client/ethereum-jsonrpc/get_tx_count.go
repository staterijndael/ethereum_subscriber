package ethereum_jsonrpc

import (
	"encoding/json"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc/models"
)

// getTxCountRPCName is the name of the JSON-RPC method for getting the transaction count of a specific address.
const getTxCountRPCName = "eth_getTransactionCount"

// GetTxCountReq represents the request for the GetTxCount method.
// It contains the address of the account for which the transaction count should be retrieved and the end block number.
type GetTxCountReq struct {
	Address  string
	EndBlock models.HexUint64
}

// Validate checks if the address field in the request is empty.
// If the address field is empty, it returns an error.
func (r *GetTxCountReq) Validate() error {
	if r.Address == "" {
		return errors.New("address field is empty")
	}

	return nil
}

// GetTxCountResp represents the response of the GetTxCount method.
// It contains the nonce of the account specified in the request.
type GetTxCountResp struct {
	Nonce models.HexUint64
}

// GetTxCount is a method of the Client struct that sends a JSON-RPC request to retrieve the transaction count of a specific address.
// It takes a GetTxCountReq as input and returns a GetTxCountResp and error.
// If the address field in the request is empty, it returns an error.
// If the request is successful, it returns the nonce of the account specified in the request.
func (c *Client) GetTxCount(req *GetTxCountReq) (*GetTxCountResp, error) {
	// Validate the input request
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	// Send the JSON-RPC request
	rawReqResp, err := c.sendJSONRPCRequest(getTxCountRPCName, []interface{}{req.Address, req.EndBlock})
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON-RPC response into GetTxCountResp
	var getTxClientResp GetTxCountResp
	err = json.Unmarshal(rawReqResp, &getTxClientResp.Nonce)
	if err != nil {
		return nil, err
	}

	// Return the response
	return &getTxClientResp, nil
}
