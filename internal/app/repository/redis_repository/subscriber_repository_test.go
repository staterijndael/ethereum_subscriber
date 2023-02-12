package redis_repository

import (
	"context"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const redisHostSubscriberRepository = "0.0.0.0"
const redisPortSubscriberRepository = "6379"
const redisPasswordSubscriberRepository = ""

func TestSubscriberRepository_AddNewSubscriber(t *testing.T) {
	ctx := context.TODO()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHostSubscriberRepository + ":" + redisPortSubscriberRepository,
		Password: redisPasswordSubscriberRepository,
		DB:       0,
	})

	type TestCase struct {
		Name                  string
		SubscriberForAddition models.Subscriber
	}

	testCases := []TestCase{
		{
			Name: "first testcase",
			SubscriberForAddition: models.Subscriber{
				Address:              "0x7885d4adcbb79d7ae83bd60b4d990206b5a357c5aa24bb5098d83788d0f1e6d2",
				SubscribeTxCount:     14,
				SubscribeBlockNumber: 15,
			},
		},
		{
			Name: "second testcase",
			SubscriberForAddition: models.Subscriber{
				Address:              "0x4eeaf2090f49b91f02b90a5797e3d73f6e0532f4111172b9e012208f81f470a6",
				SubscribeTxCount:     7,
				SubscribeBlockNumber: 19,
			},
		},
		{
			Name: "third testcase",
			SubscriberForAddition: models.Subscriber{
				Address:              "0x0f41f88706566a6e14bfc1f85d9761eb216a3b141c6eec13c820b6da569bf8a5",
				SubscribeBlockNumber: 2,
				SubscribeTxCount:     22,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			subscriberRepository := NewSubscriberRepository(redisClient, 10*time.Second)
			err := subscriberRepository.AddNewSubscriber(ctx, models.Subscriber{
				Address:              testCase.SubscriberForAddition.Address,
				SubscribeBlockNumber: testCase.SubscriberForAddition.SubscribeBlockNumber,
				SubscribeTxCount:     testCase.SubscriberForAddition.SubscribeTxCount,
			})
			assert.NoError(t, err)
		})
	}
}

func TestSubscriberRepository_GetSubscriberByAddress(t *testing.T) {
	ctx := context.TODO()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHostSubscriberRepository + ":" + redisPortSubscriberRepository,
		Password: redisPasswordSubscriberRepository,
		DB:       0,
	})

	subscriberRepository := NewSubscriberRepository(redisClient, 10*time.Second)

	type TestCase struct {
		Name                  string
		SubscriberForAddition models.Subscriber
		IsErr                 bool
	}

	testCases := []TestCase{
		{
			Name: "OK 1",
			SubscriberForAddition: models.Subscriber{
				Address:              "0x7885d4adcbb79d7ae83bd60b4d990206b5a357c5aa24bb5098d83788d0f1e6d2",
				SubscribeTxCount:     14,
				SubscribeBlockNumber: 15,
			},
			IsErr: false,
		},
		{
			Name: "OK 2",
			SubscriberForAddition: models.Subscriber{
				Address:              "0x4eeaf2090f49b91f02b90a5797e3d73f6e0532f4111172b9e012208f81f470a6",
				SubscribeTxCount:     7,
				SubscribeBlockNumber: 19,
			},
			IsErr: false,
		},
		{
			Name: "OK 3",
			SubscriberForAddition: models.Subscriber{
				Address:              "0x0f41f88706566a6e14bfc1f85d9761eb216a3b141c6eec13c820b6da569bf8a5",
				SubscribeBlockNumber: 2,
				SubscribeTxCount:     22,
			},
			IsErr: false,
		},
		{
			Name: "OK 4",
			SubscriberForAddition: models.Subscriber{
				Address:              "0x0f41f88706566a6e14bfc1f85d9761eb216a3b141c6eec13c820b6da569bf8a5",
				SubscribeBlockNumber: 2,
				SubscribeTxCount:     22,
			},
			IsErr: true,
		},
	}

	for _, testCase := range testCases {
		// we do not use t.Run() due to requirement of synchronous execution for every testcase
		// because we must get error in a last testcase due to duplicating registered address
		err := subscriberRepository.AddNewSubscriber(ctx, testCase.SubscriberForAddition)
		if testCase.IsErr {
			assert.Error(t, err)
		}

		gotSubscriber, err := subscriberRepository.GetSubscriberByAddress(ctx, testCase.SubscriberForAddition.Address)
		assert.NoError(t, err)

		assert.EqualValues(t, gotSubscriber, testCase.SubscriberForAddition)
	}
}
