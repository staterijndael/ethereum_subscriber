package scenarios

import (
	"bufio"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/cmd/scenarios/get_current_block"
	"github.com/bluntenpassant/ethereum_subscriber/cmd/scenarios/get_transactions"
	"github.com/bluntenpassant/ethereum_subscriber/cmd/scenarios/subscribe"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
)

type IParserService interface {
	GetCurrentBlock() uint64
	GetTransactions(address string) ([]*models.Transaction, error)
	Subscribe(address string) error
}

type Scenario interface {
	Present(reader *bufio.Reader) error
	GetScenarioName() string
	GetScenarioNumber() int
}

type Scenarios struct {
	scenariosByName   map[string]Scenario
	scenariosByNumber map[int]Scenario

	parserService IParserService
	reader        *bufio.Reader
}

func NewScenarios(reader *bufio.Reader, parserService IParserService) *Scenarios {
	return &Scenarios{
		parserService: parserService,
		reader:        reader,
	}
}

func (s *Scenarios) PresentScenarioByNum(num int) error {
	scenario, ok := s.scenariosByNumber[num]
	if !ok {
		return errors.New("scenario does not exist")
	}

	err := scenario.Present(s.reader)

	return err
}

func (s *Scenarios) PresentScenarioByName(name string) error {
	scenario, ok := s.scenariosByName[name]
	if !ok {
		return errors.New("scenario does not exist")
	}

	err := scenario.Present(s.reader)

	return err
}

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
