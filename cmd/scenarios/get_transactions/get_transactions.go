package get_transactions

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"strings"
)

const scenarioName = "getTransactions"
const scenarioNumber = 3

type IParserService interface {
	GetTransactions(ctx context.Context, address string) ([]*models.Transaction, error)
}

type GetTransactionsScenario struct {
	parserService IParserService
}

func NewGetTransactionsScenario(parserService IParserService) *GetTransactionsScenario {
	return &GetTransactionsScenario{
		parserService: parserService,
	}
}

func (s *GetTransactionsScenario) GetScenarioName() string {
	return scenarioName
}

func (s *GetTransactionsScenario) GetScenarioNumber() int {
	return scenarioNumber
}

func (s *GetTransactionsScenario) Present(ctx context.Context, reader *bufio.Reader) error {
	fmt.Println("Enter subscriber address: ")
	subscriberAddress, _ := reader.ReadString('\n')
	subscriberAddress = strings.TrimSpace(subscriberAddress)
	subscriberAddress = strings.TrimRight(subscriberAddress, "\n")
	subscriberAddress = strings.ToLower(subscriberAddress)

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
