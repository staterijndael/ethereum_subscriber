package main

import (
	"bufio"
	"context"
	"fmt"
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
	redis_driver "github.com/bluntenpassant/ethereum_subscriber/internal/drivers/redis"
	redis2 "github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	ctx := context.Background()

	internalConfig := config.Config{}

	// Read the YAML configuration file into a byte slice
	configData, err := ioutil.ReadFile("./config/config.local.yaml")
	if err != nil {
		panic("Error reading config.yaml: " + err.Error())
	}

	// Unmarshal the byte slice into a Config struct
	err = yaml.Unmarshal(configData, &internalConfig)
	if err != nil {
		panic("Error unmarshaling config.yaml: " + err.Error())
	}

	asyncSubscriberRepository := async_memory_repository.NewSubscriberRepository()
	asyncBlockRepository := async_memory_repository.NewBlockRepository()

	syncSubscriberRepository := sync_memory_repository.NewSubscriberRepository()
	syncBlockRepository := sync_memory_repository.NewBlockRepository()

	syncGreedySubscriberRepository := sync_greedy_memory_repository.NewSubscriberRepository()
	syncGreedyBlockRepository := sync_greedy_memory_repository.NewBlockRepository()

	redis, err := redis_driver.NewRedisClient(ctx, &redis2.Options{
		Addr:     internalConfig.Storage.Redis.Host,
		Password: internalConfig.Storage.Redis.Password,
		DB:       internalConfig.Storage.Redis.DB,
	})
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	syncRedisSubscriberRepository := sync_redis_repository.NewSubscriberRepository(redis, internalConfig.Storage.Redis.DataKeepAliveDuration)
	syncRedisBlockRepository := sync_redis_repository.NewBlockRepository(redis, internalConfig.Storage.Redis.DataKeepAliveDuration)

	syncGreedyRedisSubscriberRepository := sync_greedy_redis_repository.NewSubscriberRepository(redis, internalConfig.Storage.Redis.DataKeepAliveDuration)
	syncGreedyRedisBlockRepository := sync_greedy_redis_repository.NewBlockRepository(redis, internalConfig.Storage.Redis.DataKeepAliveDuration)

	ethereumJsonRPCClient := ethereum_jsonrpc.NewClient(internalConfig.EthereumJsonRPC.Host, internalConfig.EthereumJsonRPC.Version)

	asyncParserService := async_parser.NewParser(ethereumJsonRPCClient, asyncSubscriberRepository, asyncBlockRepository)
	syncParserService := sync_parser.NewParser(ethereumJsonRPCClient, syncSubscriberRepository, syncBlockRepository)
	syncGreedyParserService := sync_greedy_parser.NewParser(ethereumJsonRPCClient, syncGreedySubscriberRepository, syncGreedyBlockRepository)
	syncRedisParserService := sync_parser.NewParser(ethereumJsonRPCClient, syncRedisSubscriberRepository, syncRedisBlockRepository)
	syncGreedyRedisParserService := sync_greedy_parser.NewParser(ethereumJsonRPCClient, syncGreedyRedisSubscriberRepository, syncGreedyRedisBlockRepository)

	content, err := os.ReadFile("./cmd/hello_text")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n\n\n", content)

	reader := bufio.NewReader(os.Stdin)

	var presentedScenarios *scenarios.Scenarios
	if internalConfig.General.Processing == config.SyncProcessing {
		if internalConfig.General.Approach == config.GreedyApproach {
			if internalConfig.General.Storage == config.RedisStorage {
				presentedScenarios = scenarios.NewScenarios(reader, syncGreedyRedisParserService)
			} else if internalConfig.General.Storage == config.MemoryStorage {
				presentedScenarios = scenarios.NewScenarios(reader, syncGreedyParserService)
			} else {
				fmt.Println("Error: unknown storage (General -> Storage)")
				return
			}
		} else if internalConfig.General.Approach == config.ReleasingApproach {
			if internalConfig.General.Storage == config.RedisStorage {
				presentedScenarios = scenarios.NewScenarios(reader, syncRedisParserService)
			} else if internalConfig.General.Storage == config.MemoryStorage {
				presentedScenarios = scenarios.NewScenarios(reader, syncParserService)
			} else {
				fmt.Println("Error: unknown storage (General -> Storage)")
				return
			}
		} else {
			fmt.Println("Error: unknown approach type in config (General -> Approach)")
			return
		}
	} else if internalConfig.General.Processing == config.AsyncProcessing {
		if internalConfig.General.Approach == config.GreedyApproach {
			fmt.Println("Error: Greedy approach cannot be async")
			return
		}
		presentedScenarios = scenarios.NewScenarios(reader, asyncParserService)
	} else {
		fmt.Println("Error: unknown processing type in config (General -> Processing)")
		return
	}

	presentedScenarios.Init()

	for {
		// Read user input from the terminal
		fmt.Print("Enter a name or number of method: ")
		methodNumberOrName, _ := reader.ReadString('\n')
		methodNumberOrName = strings.TrimRight(methodNumberOrName, "\n")
		fmt.Println()

		if num, err := strconv.Atoi(methodNumberOrName); err == nil {
			err = presentedScenarios.PresentScenarioByNum(ctx, num)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				fmt.Println()
			}

			fmt.Println()
			continue
		}

		err = presentedScenarios.PresentScenarioByName(ctx, methodNumberOrName)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			fmt.Println()
			continue
		}

		fmt.Println()
	}
}
