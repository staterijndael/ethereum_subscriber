package config

import "time"

type ProcessingParam string
type ApproachParam string
type StorageParam string

var (
	SyncProcessing  ProcessingParam = "sync"
	AsyncProcessing ProcessingParam = "async"
)

var (
	GreedyApproach    ApproachParam = "greedy"
	ReleasingApproach ApproachParam = "releasing"
)

var (
	MemoryStorage StorageParam = "memory"
	RedisStorage  StorageParam = "redis"
)

type Config struct {
	EthereumJsonRPC EthereumJsonRPC `yaml:"ethereum_jsonrpc"`
	General         General         `yaml:"general"`
	Storage         Storage         `yaml:"storage"`
}

type EthereumJsonRPC struct {
	Host    string `yaml:"host"`
	Version string `yaml:"version"`
}

type General struct {
	Processing ProcessingParam `yaml:"processing"`
	Approach   ApproachParam   `yaml:"approach"`
	Storage    StorageParam    `yaml:"storage"`
}

type Storage struct {
	Redis Redis `yaml:"redis"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`

	DataKeepAliveDuration time.Duration `yaml:"data_keep_alive_duration"`
}
