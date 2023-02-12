package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/bluntenpassant/ethereum_subscriber/cmd"
	"github.com/bluntenpassant/ethereum_subscriber/config"
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

	var redis *redis2.Client

	// Check if our current storage is redis then we will create redis client, otherwise
	// here could be error cause unable to connect to redis, because we are using redis.Ping() inside
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

	// Init reader for further user input reading in different interfaces. We are trying to read input in main entrance
	// and pass it after in scenario interface to further handling depending on user input
	reader := bufio.NewReader(os.Stdin)

	// Init container with all helper services and repositories for further pass into usecase layer
	container := cmd.NewContainer()
	container.Init(redis, internalConfig)

	// Map that contains all application scenarios with services by 3 main parameters that can mutate current scenario choice.
	// Depends on this 3 parameters we choose convenient scneario for current usecase layer.
	presentScenariosByParams := container.GetPresentScenarioByParams(reader)

	// Pass 3 main parameters for user cli scenario choice
	presentScenario, ok := presentScenariosByParams[cmd.ModeParams{
		Approach:   internalConfig.General.Approach,
		Processing: internalConfig.General.Processing,
		Storage:    internalConfig.General.Storage,
	}]
	if !ok {
		fmt.Println("Present scenario not found for given approach, processing and storage")
		return
	}

	// Init user cli scenaior for further showing
	presentScenario.Init()

	content, err := os.ReadFile("./cmd/hello_text")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n\n\n", content)

	for {
		// Read user input from the terminal
		fmt.Print("Enter a name or number of method: ")
		methodNumberOrName, _ := reader.ReadString('\n')
		methodNumberOrName = strings.TrimRight(methodNumberOrName, "\n")
		fmt.Println()

		// If user input represented as a method number, then handle it as a number
		if num, err := strconv.Atoi(methodNumberOrName); err == nil {
			err = presentScenario.PresentScenarioByNum(ctx, num)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				fmt.Println()
			}

			fmt.Println()
			continue
		}

		// Otherwise, if user input is a method name, handle it as a string
		err = presentScenario.PresentScenarioByName(ctx, methodNumberOrName)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			fmt.Println()
			continue
		}

		fmt.Println()
	}
}
