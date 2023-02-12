package ethereum_jsonrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync/atomic"
)

// jsonRPCReq represents a JSON-RPC request to be sent to Ethereum node.
type jsonRPCReq struct {
	// JsonRPC specifies the JSON-RPC version.
	JsonRPC string `json:"jsonrpc"`

	// Method is the name of the method to be called on the Ethereum node.
	Method string `json:"method"`

	// Params is an array of parameters to be passed to the method.
	Params []interface{} `json:"params"`

	// ID is a unique identifier for the request.
	ID int `json:"id"`
}

// jsonRPCResp represents a JSON-RPC response from Ethereum node.
type jsonRPCResp struct {
	// JsonRPC specifies the JSON-RPC version.
	JsonRPC string `json:"jsonrpc"`

	// ID is a unique identifier for the request, corresponding to the ID field in jsonRPCReq.
	ID int `json:"id"`

	// Error represents an error returned by the Ethereum node.
	Error RpcError `json:"error"`

	// Result is the result of the method call.
	Result json.RawMessage `json:"result"`
}

// RpcError represents an error returned by the Ethereum node.
type RpcError struct {
	// Code is the error code returned by the Ethereum node.
	Code int `json:"code"`

	// Message is the error message returned by the Ethereum node.
	Message string `json:"message"`
}

// Client represents a client that can send JSON RPC requests to an Ethereum node
type Client struct {
	// Host is the address of the Ethereum node to connect to
	Host string

	// JsonRPC is the version of the JSON-RPC protocol to use
	JsonRPC string
	// CurrentReqID is an atomic counter used to generate unique identifiers for requests
	CurrentReqID atomic.Uint32
}

// NewClient creates a new client that can send JSON RPC requests to an Ethereum node
func NewClient(Host string, JsonRPC string) *Client {
	return &Client{
		Host:    Host,
		JsonRPC: JsonRPC,
	}
}

// sendJSONRPCRequest sends a JSON-RPC request to the Ethereum node defined in the Client struct.
// The method takes a method name and an array of parameters and returns the raw result in JSON format or an error.
func (c *Client) sendJSONRPCRequest(method string, params []interface{}) (json.RawMessage, error) {
	// Create a JSON-RPC request struct with the provided method name, parameters and a unique ID.
	request := jsonRPCReq{
		JsonRPC: c.JsonRPC,                  // JSON-RPC version string
		Method:  method,                     // Name of the method to call
		Params:  params,                     // Method parameters
		ID:      int(c.CurrentReqID.Add(1)), // Unique request ID
	}

	// Marshal the JSON-RPC request struct into a JSON payload.
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// Send a POST request to the Ethereum node with the JSON payload.
	resp, err := http.Post(c.Host, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Check if the status code of the response is not OK (200).
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	// Unmarshal the response body into a JSON-RPC response struct.
	var rpcResp jsonRPCResp
	err = json.Unmarshal(body, &rpcResp)
	if err != nil {
		return nil, err
	}

	// Check if the response contains an error message.
	if rpcResp.Error.Message != "" {
		return nil, errors.New(rpcResp.Error.Message)
	}

	// Return the raw result from the JSON-RPC response.
	return rpcResp.Result, nil
}
