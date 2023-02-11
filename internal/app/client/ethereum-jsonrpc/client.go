package ethereum_jsonrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync/atomic"
)

type jsonRPCReq struct {
	JsonRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type jsonRPCResp struct {
	JsonRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Error   RpcError        `json:"error"`
	Result  json.RawMessage `json:"result"`
}

type RpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Client struct {
	Host string

	JsonRPC      string
	CurrentReqID atomic.Uint32
}

func NewClient(Host string, JsonRPC string) *Client {
	return &Client{
		Host:    Host,
		JsonRPC: JsonRPC,
	}
}

func (c *Client) sendJSONRPCRequest(method string, params []interface{}) (json.RawMessage, error) {
	request := jsonRPCReq{
		JsonRPC: c.JsonRPC,
		Method:  method,
		Params:  params,
		ID:      int(c.CurrentReqID.Add(1)),
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(c.Host, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var rpcResp jsonRPCResp
	err = json.Unmarshal(body, &rpcResp)
	if err != nil {
		return nil, err
	}

	if rpcResp.Error.Message != "" {
		return nil, errors.New(rpcResp.Error.Message)
	}

	return rpcResp.Result, nil
}
