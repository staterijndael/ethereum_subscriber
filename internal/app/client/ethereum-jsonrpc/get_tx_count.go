package ethereum_jsonrpc

import (
	"encoding/json"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc/models"
)

const getTxCountRPCName = "eth_getTransactionCount"

type GetTxCountReq struct {
	Address  string
	EndBlock models.HexUint64
}

func (r *GetTxCountReq) Validate() error {
	if r.Address == "" {
		return errors.New("address field is empty")
	}

	return nil
}

type GetTxCountResp struct {
	Nonce models.HexUint64
}

func (c *Client) GetTxCount(req *GetTxCountReq) (*GetTxCountResp, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	rawReqResp, err := c.sendJSONRPCRequest(getTxCountRPCName, []interface{}{req.Address, req.EndBlock})
	if err != nil {
		return nil, err
	}

	var getTxClientResp GetTxCountResp
	err = json.Unmarshal(rawReqResp, &getTxClientResp.Nonce)
	if err != nil {
		return nil, err
	}

	return &getTxClientResp, nil
}
