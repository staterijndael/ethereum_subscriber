package scenarios

import (
	"bufio"
	"context"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/cmd/scenarios/get_current_block"
	"github.com/bluntenpassant/ethereum_subscriber/cmd/scenarios/get_transactions"
	"github.com/bluntenpassant/ethereum_subscriber/cmd/scenarios/subscribe"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
)

// IParserService interface of representation of parser that will be using for handling user request
type IParserService interface {
	GetCurrentBlock(ctx context.Context) (uint64, error)
	GetTransactions(ctx context.Context, address string) ([]*models.Transaction, error)
	Subscribe(ctx context.Context, address string) error
}

// Scenario represents scenario interface for further handling
type Scenario interface {
	Present(ctx context.Context, reader *bufio.Reader) error
	GetScenarioName() string
	GetScenarioNumber() int
}

// Scenarios helps to handle different user scenarios based on user input
type Scenarios struct {
	scenariosByName   map[string]Scenario
	scenariosByNumber map[int]Scenario

	parserService IParserService
	reader        *bufio.Reader
}

// NewScenarios just returns pointer to Scenarios object with filled service field
func NewScenarios(reader *bufio.Reader, parserService IParserService) *Scenarios {
	return &Scenarios{
		parserService: parserService,
		reader:        reader,
	}
}

// PresentScenarioByNum handles scenario method that was attached for current scenario
// based on scenario name represented as a integer
func (s *Scenarios) PresentScenarioByNum(ctx context.Context, num int) error {
	scenario, ok := s.scenariosByNumber[num]
	if !ok {
		return errors.New("scenario does not exist")
	}

	err := scenario.Present(ctx, s.reader)

	return err
}

// PresentScenarioByName handles scenario method that was attached for current scenario
// based on scenario name represented as a string
func (s *Scenarios) PresentScenarioByName(ctx context.Context, name string) error {
	scenario, ok := s.scenariosByName[name]
	if !ok {
		return errors.New("scenario does not exist")
	}

	err := scenario.Present(ctx, s.reader)

	return err
}

// Init container with all helper services and repositories for further pass into usecase layer
func (s *Scenarios) Init() {
	getCurrentBlockScenario := get_current_block.NewGetCurrentBlockScenario(s.parserService)
	getTransactionsScenario := get_transactions.NewGetTransactionsScenario(s.parserService)
	subscribeScenario := subscribe.NewSubscribeScenario(s.parserService)

	s.scenariosByName = map[string]Scenario{
		getCurrentBlockScenario.GetScenarioName(): getCurrentBlockScenario,
		getTransactionsScenario.GetScenarioName(): getTransactionsScenario,
		subscribeScenario.GetScenarioName():       subscribeScenario,
	}

	s.scenariosByNumber = map[int]Scenario{
		getCurrentBlockScenario.GetScenarioNumber(): getCurrentBlockScenario,
		getTransactionsScenario.GetScenarioNumber(): getTransactionsScenario,
		subscribeScenario.GetScenarioNumber():       subscribeScenario,
	}
}
