package config

type Processing string

var (
	Sync  Processing = "sync"
	Async Processing = "async"
)

type Config struct {
	EthereumJsonRPC EthereumJsonRPC `yaml:"ethereum_jsonrpc"`
	General         General         `yaml:"general"`
}

type EthereumJsonRPC struct {
	Host    string `yaml:"host"`
	Version string `yaml:"version"`
}

type General struct {
	Processing Processing `yaml:"processing"`
}
