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

	redis, err := redis_driver.NewRedisClient(ctx, &redis2.Options{
		Addr:     internalConfig.Storage.Redis.Host,
		Password: internalConfig.Storage.Redis.Password,
		DB:       internalConfig.Storage.Redis.DB,
	})
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	reader := bufio.NewReader(os.Stdin)

	container := cmd.NewContainer()
	container.Init(redis, internalConfig)

	presentScenariosByParams := container.GetPresentScenario(reader)

	presentScenario, ok := presentScenariosByParams[cmd.PresentScenarioByParams{
		Approach:   internalConfig.General.Approach,
		Processing: internalConfig.General.Processing,
		Storage:    internalConfig.General.Storage,
	}]
	if !ok {
		fmt.Println("Present scenario not found for given approach, processing and storage")
		return
	}

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

		if num, err := strconv.Atoi(methodNumberOrName); err == nil {
			err = presentScenario.PresentScenarioByNum(ctx, num)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				fmt.Println()
			}

			fmt.Println()
			continue
		}

		err = presentScenario.PresentScenarioByName(ctx, methodNumberOrName)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			fmt.Println()
			continue
		}

		fmt.Println()
	}
}
