package get_transactions

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/utils"
)

// Current scenario name that attached to this scenario and using for spotting method based on scenario name
const scenarioName = "GetTransactions"

// Current scenario number that attached to this scenario and using for spotting method based on scenario number
const scenarioNumber = 3

// IParserService interface of representation of parser that will be using for handling user request
type IParserService interface {
	GetTransactions(ctx context.Context, address string) ([]*models.Transaction, error)
}

// GetTransactionsScenario represents scenario object for further handling
type GetTransactionsScenario struct {
	parserService IParserService
}

// NewGetTransactionsScenario just returns pointer to GetTransactionsScenario object with filled service field
func NewGetTransactionsScenario(parserService IParserService) *GetTransactionsScenario {
	return &GetTransactionsScenario{
		parserService: parserService,
	}
}

// GetScenarioName returns scenario method name that was attached for current scenario for further handling
// generally using for routing scenarios based on user input when user type method name represented as a string
func (s *GetTransactionsScenario) GetScenarioName() string {
	return scenarioName
}

// GetScenarioNumber returns scenario method number that was attached for current scenario for further handling
// generally using for routing scenarios based on user input when user type method name represented as a number
func (s *GetTransactionsScenario) GetScenarioNumber() int {
	return scenarioNumber
}

// Present represents user scenario for GetTransactionScenario method and trying to handle it based on user input
func (s *GetTransactionsScenario) Present(ctx context.Context, reader *bufio.Reader) error {
	fmt.Println("Enter subscriber address: ")
	subscriberAddress, _ := reader.ReadString('\n')
	subscriberAddress = utils.ClearString(subscriberAddress)

	transactions, err := s.parserService.GetTransactions(ctx, subscriberAddress)
	if err != nil {
		return err
	}

	rawTransactions, err := json.MarshalIndent(transactions, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(rawTransactions))

	return nil
}
