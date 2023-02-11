package get_current_block

import (
	"bufio"
	"context"
	"errors"
	"fmt"
)

const scenarioName = "GetCurrentBlock"
const scenarioNumber = 1

type IParserService interface {
	GetCurrentBlock(ctx context.Context) (uint64, error)
}

type GetCurrentBlockScenario struct {
	parserService IParserService
}

func NewGetCurrentBlockScenario(parserService IParserService) *GetCurrentBlockScenario {
	return &GetCurrentBlockScenario{
		parserService: parserService,
	}
}

func (s *GetCurrentBlockScenario) GetScenarioName() string {
	return scenarioName
}

func (s *GetCurrentBlockScenario) GetScenarioNumber() int {
	return scenarioNumber
}

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
