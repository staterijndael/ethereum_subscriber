package subscribe

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/utils"
	"strconv"
	"unicode/utf8"
)

// Current scenario name that attached to this scenario and using for spotting method based on scenario name
const scenarioName = "Subscribe"

// Current scenario name that attached to this scenario and using for spotting method based on scenario name
const scenarioNumber = 2

// Length of admissible Ethereum address length
const addressLength = 42

// IParserService interface of representation of parser that will be using for handling user request
type IParserService interface {
	Subscribe(ctx context.Context, address string) error
}

// SubscribeScenario represents scenario object for further handling
type SubscribeScenario struct {
	parserService IParserService
}

// NewSubscribeScenario just returns pointer to GetTransactionsScenario object with filled service field
func NewSubscribeScenario(parserService IParserService) *SubscribeScenario {
	return &SubscribeScenario{
		parserService: parserService,
	}
}

// GetScenarioName returns scenario method name that was attached for current scenario for further handling
// generally using for routing scenarios based on user input when user type method name represented as a string
func (s *SubscribeScenario) GetScenarioName() string {
	return scenarioName
}

// GetScenarioNumber returns scenario method number that was attached for current scenario for further handling
// generally using for routing scenarios based on user input when user type method name represented as a number
func (s *SubscribeScenario) GetScenarioNumber() int {
	return scenarioNumber
}

// Present represents user scenario for GetTransactionScenario method and trying to handle it based on user input
func (s *SubscribeScenario) Present(ctx context.Context, reader *bufio.Reader) error {
	fmt.Println("Enter subscribe address: ")
	subscribeAddress, _ := reader.ReadString('\n')
	subscribeAddress = utils.ClearString(subscribeAddress)

	if utf8.RuneCountInString(subscribeAddress) != 42 {
		return errors.New("address length should be " + strconv.Itoa(addressLength))
	}

	err := s.parserService.Subscribe(ctx, subscribeAddress)
	if err != nil {
		return err
	}

	fmt.Println("true")

	return nil
}
