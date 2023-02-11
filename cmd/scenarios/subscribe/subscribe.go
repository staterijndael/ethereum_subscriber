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

const scenarioName = "Subscribe"
const scenarioNumber = 2

const addressLength = 42

type IParserService interface {
	Subscribe(ctx context.Context, address string) error
}

type SubscribeScenario struct {
	parserService IParserService
}

func NewSubscribeScenario(parserService IParserService) *SubscribeScenario {
	return &SubscribeScenario{
		parserService: parserService,
	}
}

func (s *SubscribeScenario) GetScenarioName() string {
	return scenarioName
}

func (s *SubscribeScenario) GetScenarioNumber() int {
	return scenarioNumber
}

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
