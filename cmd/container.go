package cmd

import (
	"bufio"
	"github.com/bluntenpassant/ethereum_subscriber/cmd/scenarios"
	"github.com/bluntenpassant/ethereum_subscriber/config"
	ethereum_jsonrpc "github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/repository/async_memory_repository"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/repository/sync_greedy_memory_repository"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/repository/sync_greedy_redis_repository"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/repository/sync_memory_repository"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/repository/sync_redis_repository"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/service/async_parser"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/service/sync_greedy_parser"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/service/sync_parser"
	redis2 "github.com/redis/go-redis/v9"
)

type Container struct {
	asyncParserService           *async_parser.Parser
	syncParserService            *sync_parser.Parser
	syncGreedyParserService      *sync_greedy_parser.Parser
	syncRedisParserService       *sync_parser.Parser
	syncGreedyRedisParserService *sync_greedy_parser.Parser
}

func NewContainer() *Container {
	return &Container{}
}

type PresentScenarioByParams struct {
	Approach   config.ApproachParam
	Processing config.ProcessingParam
	Storage    config.StorageParam
}

func (c *Container) GetPresentScenario(reader *bufio.Reader) map[PresentScenarioByParams]*scenarios.Scenarios {
	presentScenarioByParams := map[PresentScenarioByParams]*scenarios.Scenarios{
		PresentScenarioByParams{
			Processing: config.SyncProcessing,
			Approach:   config.GreedyApproach,
			Storage:    config.RedisStorage,
		}: scenarios.NewScenarios(reader, c.syncGreedyRedisParserService),

		PresentScenarioByParams{
			Processing: config.SyncProcessing,
			Approach:   config.GreedyApproach,
			Storage:    config.MemoryStorage,
		}: scenarios.NewScenarios(reader, c.syncGreedyParserService),

		PresentScenarioByParams{
			Processing: config.SyncProcessing,
			Approach:   config.ReleasingApproach,
			Storage:    config.RedisStorage,
		}: scenarios.NewScenarios(reader, c.syncRedisParserService),

		PresentScenarioByParams{
			Processing: config.SyncProcessing,
			Approach:   config.ReleasingApproach,
			Storage:    config.MemoryStorage,
		}: scenarios.NewScenarios(reader, c.syncParserService),

		PresentScenarioByParams{
			Processing: config.AsyncProcessing,
			Approach:   config.ReleasingApproach,
			Storage:    config.MemoryStorage,
		}: scenarios.NewScenarios(reader, c.asyncParserService),
	}

	return presentScenarioByParams
}

func (c *Container) Init(redis *redis2.Client, config config.Config) {
	asyncSubscriberRepository := async_memory_repository.NewSubscriberRepository()
	asyncBlockRepository := async_memory_repository.NewBlockRepository()

	syncSubscriberRepository := sync_memory_repository.NewSubscriberRepository()
	syncBlockRepository := sync_memory_repository.NewBlockRepository()

	syncGreedySubscriberRepository := sync_greedy_memory_repository.NewSubscriberRepository()
	syncGreedyBlockRepository := sync_greedy_memory_repository.NewBlockRepository()

	syncRedisSubscriberRepository := sync_redis_repository.NewSubscriberRepository(redis, config.Storage.Redis.DataKeepAliveDuration)
	syncRedisBlockRepository := sync_redis_repository.NewBlockRepository(redis, config.Storage.Redis.DataKeepAliveDuration)

	syncGreedyRedisSubscriberRepository := sync_greedy_redis_repository.NewSubscriberRepository(redis, config.Storage.Redis.DataKeepAliveDuration)
	syncGreedyRedisBlockRepository := sync_greedy_redis_repository.NewBlockRepository(redis, config.Storage.Redis.DataKeepAliveDuration)

	ethereumJsonRPCClient := ethereum_jsonrpc.NewClient(config.EthereumJsonRPC.Host, config.EthereumJsonRPC.Version)

	c.asyncParserService = async_parser.NewParser(ethereumJsonRPCClient, asyncSubscriberRepository, asyncBlockRepository)
	c.syncParserService = sync_parser.NewParser(ethereumJsonRPCClient, syncSubscriberRepository, syncBlockRepository)
	c.syncGreedyParserService = sync_greedy_parser.NewParser(ethereumJsonRPCClient, syncGreedySubscriberRepository, syncGreedyBlockRepository)
	c.syncRedisParserService = sync_parser.NewParser(ethereumJsonRPCClient, syncRedisSubscriberRepository, syncRedisBlockRepository)
	c.syncGreedyRedisParserService = sync_greedy_parser.NewParser(ethereumJsonRPCClient, syncGreedyRedisSubscriberRepository, syncGreedyRedisBlockRepository)
}
