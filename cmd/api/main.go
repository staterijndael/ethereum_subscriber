package main

import (
	"context"
	"fmt"
	"github.com/bluntenpassant/ethereum_subscriber/cmd"
	"github.com/bluntenpassant/ethereum_subscriber/config"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/handlers"
	redis_driver "github.com/bluntenpassant/ethereum_subscriber/internal/drivers/redis"
	redis2 "github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"io/ioutil"
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

	var redis *redis2.Client

	if internalConfig.General.Storage == config.RedisStorage {
		redis, err = redis_driver.NewRedisClient(ctx, &redis2.Options{
			Addr:     internalConfig.Storage.Redis.Host,
			Password: internalConfig.Storage.Redis.Password,
			DB:       internalConfig.Storage.Redis.DB,
		})
		if err != nil {
			fmt.Println("Error: " + err.Error())
			return
		}
	}

	container := cmd.NewContainer()
	container.Init(redis, internalConfig)

	getServiceByParams := container.GetServiceByParams()

	parserService, ok := getServiceByParams[cmd.ModeParams{
		Approach:   internalConfig.General.Approach,
		Processing: internalConfig.General.Processing,
		Storage:    internalConfig.General.Storage,
	}]
	if !ok {
		fmt.Println("Present scenario not found for given approach, processing and storage")
		return
	}

	httpHandler := handlers.NewHandler(parserService)

	fmt.Println("HTTP Server started...")
	httpHandler.Start(internalConfig.Http)
}
