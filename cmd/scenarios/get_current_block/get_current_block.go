package get_current_block

import (
	"bufio"
	"errors"
	"fmt"
)

const scenarioName = "GetCurrentBlock"
const scenarioNumber = 1

type IParserService interface {
	GetCurrentBlock() uint64
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

func (s *GetCurrentBlockScenario) Present(reader *bufio.Reader) error {
	currentBlock := s.parserService.GetCurrentBlock()
	if currentBlock == 0 {
		return errors.New("current block is not parsed yet")
	}

	fmt.Println(currentBlock)

	return nil
}
