package cmd

import (
	"bufio"
	"context"
	"github.com/bluntenpassant/ethereum_subscriber/cmd/scenarios"
	"github.com/bluntenpassant/ethereum_subscriber/config"
	ethereum_jsonrpc "github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/repository/greedy_memory_repository"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/repository/greedy_redis_repository"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/repository/memory_repository"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/repository/redis_repository"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/service/async_parser"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/service/sync_greedy_parser"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/service/sync_parser"
	redis2 "github.com/redis/go-redis/v9"
)

type IParserService interface {
	GetCurrentBlock(ctx context.Context) (uint64, error)
	GetTransactions(ctx context.Context, address string) ([]*models.Transaction, error)
	Subscribe(ctx context.Context, address string) error
}

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

type ModeParams struct {
	Approach   config.ApproachParam
	Processing config.ProcessingParam
	Storage    config.StorageParam
}

func (c *Container) GetServiceByParams() map[ModeParams]IParserService {
	presentScenarioByParams := map[ModeParams]IParserService{
		ModeParams{
			Processing: config.SyncProcessing,
			Approach:   config.GreedyApproach,
			Storage:    config.RedisStorage,
		}: c.syncGreedyRedisParserService,

		ModeParams{
			Processing: config.SyncProcessing,
			Approach:   config.GreedyApproach,
			Storage:    config.MemoryStorage,
		}: c.syncGreedyParserService,

		ModeParams{
			Processing: config.SyncProcessing,
			Approach:   config.ReleasingApproach,
			Storage:    config.RedisStorage,
		}: c.syncRedisParserService,

		ModeParams{
			Processing: config.SyncProcessing,
			Approach:   config.ReleasingApproach,
			Storage:    config.MemoryStorage,
		}: c.syncParserService,

		ModeParams{
			Processing: config.AsyncProcessing,
			Approach:   config.ReleasingApproach,
			Storage:    config.MemoryStorage,
		}: c.asyncParserService,
	}

	return presentScenarioByParams
}

func (c *Container) GetPresentScenarioByParams(reader *bufio.Reader) map[ModeParams]*scenarios.Scenarios {
	presentScenarioByParams := map[ModeParams]*scenarios.Scenarios{
		ModeParams{
			Processing: config.SyncProcessing,
			Approach:   config.GreedyApproach,
			Storage:    config.RedisStorage,
		}: scenarios.NewScenarios(reader, c.syncGreedyRedisParserService),

		ModeParams{
			Processing: config.SyncProcessing,
			Approach:   config.GreedyApproach,
			Storage:    config.MemoryStorage,
		}: scenarios.NewScenarios(reader, c.syncGreedyParserService),

		ModeParams{
			Processing: config.SyncProcessing,
			Approach:   config.ReleasingApproach,
			Storage:    config.RedisStorage,
		}: scenarios.NewScenarios(reader, c.syncRedisParserService),

		ModeParams{
			Processing: config.SyncProcessing,
			Approach:   config.ReleasingApproach,
			Storage:    config.MemoryStorage,
		}: scenarios.NewScenarios(reader, c.syncParserService),

		ModeParams{
			Processing: config.AsyncProcessing,
			Approach:   config.ReleasingApproach,
			Storage:    config.MemoryStorage,
		}: scenarios.NewScenarios(reader, c.asyncParserService),
	}

	return presentScenarioByParams
}

func (c *Container) Init(redis *redis2.Client, config config.Config) {
	asyncSubscriberRepository := memory_repository.NewSubscriberRepository()
	asyncBlockRepository := memory_repository.NewBlockRepository()

	syncGreedySubscriberRepository := greedy_memory_repository.NewSubscriberRepository()
	syncGreedyBlockRepository := greedy_memory_repository.NewBlockRepository()

	syncRedisSubscriberRepository := redis_repository.NewSubscriberRepository(redis, config.Storage.Redis.DataKeepAliveDuration)
	syncRedisBlockRepository := redis_repository.NewBlockRepository(redis, config.Storage.Redis.DataKeepAliveDuration)

	syncGreedyRedisSubscriberRepository := greedy_redis_repository.NewSubscriberRepository(redis, config.Storage.Redis.DataKeepAliveDuration)
	syncGreedyRedisBlockRepository := greedy_redis_repository.NewBlockRepository(redis, config.Storage.Redis.DataKeepAliveDuration)

	ethereumJsonRPCClient := ethereum_jsonrpc.NewClient(config.EthereumJsonRPC.Host, config.EthereumJsonRPC.Version)

	c.asyncParserService = async_parser.NewParser(ethereumJsonRPCClient, asyncSubscriberRepository, asyncBlockRepository)
	c.syncParserService = sync_parser.NewParser(ethereumJsonRPCClient, asyncSubscriberRepository, asyncBlockRepository)
	c.syncGreedyParserService = sync_greedy_parser.NewParser(ethereumJsonRPCClient, syncGreedySubscriberRepository, syncGreedyBlockRepository)
	c.syncRedisParserService = sync_parser.NewParser(ethereumJsonRPCClient, syncRedisSubscriberRepository, syncRedisBlockRepository)
	c.syncGreedyRedisParserService = sync_greedy_parser.NewParser(ethereumJsonRPCClient, syncGreedyRedisSubscriberRepository, syncGreedyRedisBlockRepository)
}
