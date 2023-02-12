package get_current_block

import (
	"bufio"
	"context"
	"errors"
	"fmt"
)

// Current scenario name that attached to this scenario and using for spotting method based on scenario name
const scenarioName = "GetCurrentBlock"

// Current scenario number that attached to this scenario and using for spotting method based on scenario number
const scenarioNumber = 1

// IParserService interface of representation of parser that will be using for handling user request
type IParserService interface {
	GetCurrentBlock(ctx context.Context) (uint64, error)
}

// GetCurrentBlockScenario represents scenario object for further handling
type GetCurrentBlockScenario struct {
	parserService IParserService
}

// NewGetCurrentBlockScenario just returns pointer to GetCurrentBlockScenario object with filled service field
func NewGetCurrentBlockScenario(parserService IParserService) *GetCurrentBlockScenario {
	return &GetCurrentBlockScenario{
		parserService: parserService,
	}
}

// GetScenarioName returns scenario method name that was attached for current scenario for further handling
// generally using for routing scenarios based on user input when user type method name represented as a string
func (s *GetCurrentBlockScenario) GetScenarioName() string {
	return scenarioName
}

// GetScenarioNumber returns scenario method number that was attached for current scenario for further handling
// generally using for routing scenarios based on user input when user type method name represented as a number
func (s *GetCurrentBlockScenario) GetScenarioNumber() int {
	return scenarioNumber
}

// Present represents user scenario for GetCurrentBlock method and trying to handle it based on user input
func (s *GetCurrentBlockScenario) Present(ctx context.Context, reader *bufio.Reader) error {
	currentBlock, err := s.parserService.GetCurrentBlock(ctx)
	if err != nil {
		return err
	}
	if currentBlock == 0 {
		return errors.New("current block is not parsed yet")
	}

	fmt.Println(currentBlock)

	return nil
}
