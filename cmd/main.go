package main

import (
	"bufio"
	"fmt"
	"github.com/bluntenpassant/ethereum_subscriber/cmd/scenarios"
	"github.com/bluntenpassant/ethereum_subscriber/config"
	ethereum_jsonrpc "github.com/bluntenpassant/ethereum_subscriber/internal/app/client/ethereum-jsonrpc"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/repository/async_memory_repository"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/service/async_parser"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/service/sync_parser"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
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

	syncSubscriberRepository := async_memory_repository.NewSubscriberRepository()
	syncBlockRepository := async_memory_repository.NewBlockRepository()

	ethereumJsonRPCClient := ethereum_jsonrpc.NewClient(internalConfig.EthereumJsonRPC.Host, internalConfig.EthereumJsonRPC.Version)

	asyncParserService := async_parser.NewParser(ethereumJsonRPCClient, asyncSubscriberRepository, asyncBlockRepository)
	syncParserService := sync_parser.NewParser(ethereumJsonRPCClient, syncSubscriberRepository, syncBlockRepository)

	content, err := os.ReadFile("./cmd/hello_text")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s\n\n\n", content)

	reader := bufio.NewReader(os.Stdin)

	var presentedScenarios *scenarios.Scenarios
	if internalConfig.General.Processing == config.Sync {
		presentedScenarios = scenarios.NewScenarios(reader, syncParserService)
	} else if internalConfig.General.Processing == config.Async {
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
			err = presentedScenarios.PresentScenarioByNum(num)
			if err != nil {
				fmt.Println("Error: " + err.Error())
				fmt.Println()
			}

			fmt.Println()
			continue
		}

		err = presentedScenarios.PresentScenarioByName(methodNumberOrName)
		if err != nil {
			fmt.Println("Error: " + err.Error())
			fmt.Println()
			continue
		}

		fmt.Println()
	}
}
