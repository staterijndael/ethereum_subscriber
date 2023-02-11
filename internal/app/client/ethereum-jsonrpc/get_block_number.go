package ethereum_jsonrpc

import (
	"encoding/json"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc/models"
)

const getBlockNumberRPCName = "eth_blockNumber"

type GetBlockNumberResp struct {
	BlockNumber models.HexUint64
}

func (c *Client) GetBlockNumber() (*GetBlockNumberResp, error) {
	rawReqResp, err := c.sendJSONRPCRequest(getBlockNumberRPCName, []interface{}{})
	if err != nil {
		return nil, err
	}

	var getTxClientResp GetBlockNumberResp
	err = json.Unmarshal(rawReqResp, &getTxClientResp.BlockNumber)
	if err != nil {
		return nil, err
	}

	return &getTxClientResp, nil
}
