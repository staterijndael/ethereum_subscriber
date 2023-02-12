package greedy_redis_repository

import (
	"context"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

const redisHostSubscriberRepository = "0.0.0.0"
const redisPortSubscriberRepository = "6379"
const redisPasswordSubscriberRepository = ""

func TestSubscriberRepository_GetLastTransaction(t *testing.T) {
	ctx := context.TODO()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHostSubscriberRepository + ":" + redisPortSubscriberRepository,
		Password: redisPasswordSubscriberRepository,
		DB:       0,
	})
	subscriberRepository := NewSubscriberRepository(redisClient, 10*time.Second)

	type TestCase struct {
		Name         string
		Transactions []*models.Transaction
		Subscriber   models.Subscriber
		IsEqual      bool
	}

	testCases := []TestCase{
		{
			Name: "OK 1",
			Transactions: []*models.Transaction{{
				BlockHash:   "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50f",
				BlockNumber: 16614479,
				From:        "0x45849a974058661eb2128aceb60d2c6ed99e2a14",
				Gas:         *big.NewInt(16338636656),
				GasPrice:    *big.NewInt(16338636656),
				Hash:        "0x0f41f88706566a6e14bfc1f85d9761eb216a3b141c6eec13c820b6da569bf8a5",
				Input:       "Function: transfer(address recipient, uint256 amount)\n\nMethodID: 0xa9059cbb\n[0]:  0000000000000000000000006cc8dcbca746a6e4fdefb98e1d0df903b107fd21\n[1]:  00000000000000000000000000000000000000000000000000000081da747e22",
				Nonce:       0x0000000000000000,
				To:          "0x388c818ca8b9251b393131c08a736a67ccb19297",
			}},
			Subscriber: models.Subscriber{
				Address:              "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50f",
				SubscribeBlockNumber: 5,
				SubscribeTxCount:     10,
			},
			IsEqual: true,
		},
		{
			Name: "OK 2",
			Transactions: []*models.Transaction{{
				BlockHash:   "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50",
				BlockNumber: 16614479,
				From:        "0x45849a974058661eb2128aceb60d2c6ed99e2a14",
				Gas:         *big.NewInt(16338636656),
				GasPrice:    *big.NewInt(16338636656),
				Hash:        "0x0f41f88706566a6e14bfc1f85d9761eb216a3b141c6eec13c820b6da569bf8a5",
				Input:       "Function: transfer(address recipient, uint256 amount)\n\nMethodID: 0xa9059cbb\n[0]:  0000000000000000000000006cc8dcbca746a6e4fdefb98e1d0df903b107fd21\n[1]:  00000000000000000000000000000000000000000000000000000081da747e22",
				Nonce:       0x0000000000000000,
				To:          "0x388c818ca8b9251b393131c08a736a67ccb19297",
			}},
			Subscriber: models.Subscriber{
				Address:              "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50",
				SubscribeBlockNumber: 5,
				SubscribeTxCount:     10,
			},
			IsEqual: true,
		},
		{
			Name: "OK 2",
			Transactions: []*models.Transaction{{
				BlockHash:   "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50f",
				BlockNumber: 16614479,
				From:        "0x45849a974058661eb2128aceb60d2c6ed99e2a14",
				Gas:         *big.NewInt(16338636656),
				GasPrice:    *big.NewInt(16338636656),
				Hash:        "0x0f41f88706566a6e14bfc1f85d9761eb216a3b141c6eec13c820b6da569bf8a5",
				Input:       "Function: transfer(address recipient, uint256 amount)\n\nMethodID: 0xa9059cbb\n[0]:  0000000000000000000000006cc8dcbca746a6e4fdefb98e1d0df903b107fd21\n[1]:  00000000000000000000000000000000000000000000000000000081da747e22",
				Nonce:       0x0000000000000000,
				To:          "0x388c818ca8b9251b393131c08a736a67ccb19297",
			}},
			Subscriber: models.Subscriber{
				Address:              "",
				SubscribeBlockNumber: 5,
				SubscribeTxCount:     10,
			},
			IsEqual: false,
		},
	}

	for _, testCase := range testCases {
		// we do not use t.Run() due to requirement of synchronous execution for every testcase
		// because we must get error in a last testcase due to duplicating registered address
		err := subscriberRepository.AddNewSubscriber(ctx, testCase.Subscriber)
		assert.NoError(t, err)
		err = subscriberRepository.AddTransactions(ctx, testCase.Subscriber.Address, testCase.Transactions)
		assert.NoError(t, err)

		txs, err := subscriberRepository.GetTransactionsReversed(ctx, testCase.Subscriber.Address)
		assert.NoError(t, err)

		expectedTxs := models.ReverseTransactionsCopy(testCase.Transactions)

		if testCase.IsEqual {
			assert.EqualValues(t, expectedTxs[len(expectedTxs)-1], txs[len(txs)-1])
		}
	}
}

func TestSubscriberRepository_GetTransactionsReversed(t *testing.T) {
	ctx := context.TODO()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHostSubscriberRepository + ":" + redisPortSubscriberRepository,
		Password: redisPasswordSubscriberRepository,
		DB:       0,
	})
	subscriberRepository := NewSubscriberRepository(redisClient, 10*time.Second)

	type TestCase struct {
		Name         string
		Transactions []*models.Transaction
		Subscriber   models.Subscriber
		IsEqual      bool
	}

	testCases := []TestCase{
		{
			Name: "OK 1",
			Transactions: []*models.Transaction{{
				BlockHash:   "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50f",
				BlockNumber: 16614479,
				From:        "0x45849a974058661eb2128aceb60d2c6ed99e2a14",
				Gas:         *big.NewInt(16338636656),
				GasPrice:    *big.NewInt(16338636656),
				Hash:        "0x0f41f88706566a6e14bfc1f85d9761eb216a3b141c6eec13c820b6da569bf8a5",
				Input:       "Function: transfer(address recipient, uint256 amount)\n\nMethodID: 0xa9059cbb\n[0]:  0000000000000000000000006cc8dcbca746a6e4fdefb98e1d0df903b107fd21\n[1]:  00000000000000000000000000000000000000000000000000000081da747e22",
				Nonce:       0x0000000000000000,
				To:          "0x388c818ca8b9251b393131c08a736a67ccb19297",
			}},
			Subscriber: models.Subscriber{
				Address:              "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50f",
				SubscribeBlockNumber: 5,
				SubscribeTxCount:     10,
			},
			IsEqual: true,
		},
		{
			Name: "OK 2",
			Transactions: []*models.Transaction{{
				BlockHash:   "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50",
				BlockNumber: 16614479,
				From:        "0x45849a974058661eb2128aceb60d2c6ed99e2a14",
				Gas:         *big.NewInt(16338636656),
				GasPrice:    *big.NewInt(16338636656),
				Hash:        "0x0f41f88706566a6e14bfc1f85d9761eb216a3b141c6eec13c820b6da569bf8a5",
				Input:       "Function: transfer(address recipient, uint256 amount)\n\nMethodID: 0xa9059cbb\n[0]:  0000000000000000000000006cc8dcbca746a6e4fdefb98e1d0df903b107fd21\n[1]:  00000000000000000000000000000000000000000000000000000081da747e22",
				Nonce:       0x0000000000000000,
				To:          "0x388c818ca8b9251b393131c08a736a67ccb19297",
			}},
			Subscriber: models.Subscriber{
				Address:              "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50",
				SubscribeBlockNumber: 5,
				SubscribeTxCount:     10,
			},
			IsEqual: true,
		},
		{
			Name: "OK 2",
			Transactions: []*models.Transaction{{
				BlockHash:   "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50f",
				BlockNumber: 16614479,
				From:        "0x45849a974058661eb2128aceb60d2c6ed99e2a14",
				Gas:         *big.NewInt(16338636656),
				GasPrice:    *big.NewInt(16338636656),
				Hash:        "0x0f41f88706566a6e14bfc1f85d9761eb216a3b141c6eec13c820b6da569bf8a5",
				Input:       "Function: transfer(address recipient, uint256 amount)\n\nMethodID: 0xa9059cbb\n[0]:  0000000000000000000000006cc8dcbca746a6e4fdefb98e1d0df903b107fd21\n[1]:  00000000000000000000000000000000000000000000000000000081da747e22",
				Nonce:       0x0000000000000000,
				To:          "0x388c818ca8b9251b393131c08a736a67ccb19297",
			}},
			Subscriber: models.Subscriber{
				Address:              "",
				SubscribeBlockNumber: 5,
				SubscribeTxCount:     10,
			},
			IsEqual: false,
		},
	}

	for _, testCase := range testCases {
		// we do not use t.Run() due to requirement of synchronous execution for every testcase
		// because we must get error in a last testcase due to duplicating registered address
		err := subscriberRepository.AddNewSubscriber(ctx, testCase.Subscriber)
		assert.NoError(t, err)
		err = subscriberRepository.AddTransactions(ctx, testCase.Subscriber.Address, testCase.Transactions)
		assert.NoError(t, err)

		txs, err := subscriberRepository.GetTransactionsReversed(ctx, testCase.Subscriber.Address)
		assert.NoError(t, err)

		expectedTxs := models.ReverseTransactionsCopy(testCase.Transactions)

		if testCase.IsEqual {
			assert.EqualValues(t, expectedTxs, txs)
		}
	}
}

func TestSubscriberRepository_AddTransactions(t *testing.T) {
	ctx := context.TODO()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisHostSubscriberRepository + ":" + redisPortSubscriberRepository,
		Password: redisPasswordSubscriberRepository,
		DB:       0,
	})
	subscriberRepository := NewSubscriberRepository(redisClient, 10*time.Second)

	type TestCase struct {
		Name         string
		Transactions []*models.Transaction
		Subscriber   models.Subscriber
		IsErr        bool
	}

	testCases := []TestCase{
		{
			Name: "OK 1",
			Transactions: []*models.Transaction{{
				BlockHash:   "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50f",
				BlockNumber: 16614479,
				From:        "0x45849a974058661eb2128aceb60d2c6ed99e2a14",
				Gas:         *big.NewInt(16338636656),
				GasPrice:    *big.NewInt(16338636656),
				Hash:        "0x0f41f88706566a6e14bfc1f85d9761eb216a3b141c6eec13c820b6da569bf8a5",
				Input:       "Function: transfer(address recipient, uint256 amount)\n\nMethodID: 0xa9059cbb\n[0]:  0000000000000000000000006cc8dcbca746a6e4fdefb98e1d0df903b107fd21\n[1]:  00000000000000000000000000000000000000000000000000000081da747e22",
				Nonce:       0x0000000000000000,
				To:          "0x388c818ca8b9251b393131c08a736a67ccb19297",
			}},
			Subscriber: models.Subscriber{
				Address:              "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50f",
				SubscribeBlockNumber: 5,
				SubscribeTxCount:     10,
			},
			IsErr: false,
		},
		{
			Name: "OK 2",
			Transactions: []*models.Transaction{{
				BlockHash:   "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50",
				BlockNumber: 16614479,
				From:        "0x45849a974058661eb2128aceb60d2c6ed99e2a14",
				Gas:         *big.NewInt(16338636656),
				GasPrice:    *big.NewInt(16338636656),
				Hash:        "0x0f41f88706566a6e14bfc1f85d9761eb216a3b141c6eec13c820b6da569bf8a5",
				Input:       "Function: transfer(address recipient, uint256 amount)\n\nMethodID: 0xa9059cbb\n[0]:  0000000000000000000000006cc8dcbca746a6e4fdefb98e1d0df903b107fd21\n[1]:  00000000000000000000000000000000000000000000000000000081da747e22",
				Nonce:       0x0000000000000000,
				To:          "0x388c818ca8b9251b393131c08a736a67ccb19297",
			}},
			Subscriber: models.Subscriber{
				Address:              "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50",
				SubscribeBlockNumber: 5,
				SubscribeTxCount:     10,
			},
			IsErr: false,
		},
		{
			Name: "OK 2",
			Transactions: []*models.Transaction{{
				BlockHash:   "0xdfcbf46fb8ae7c1769e28a864b7e49d66e931214010a3d0c3c82248d2b0bd50f",
				BlockNumber: 16614479,
				From:        "0x45849a974058661eb2128aceb60d2c6ed99e2a14",
				Gas:         *big.NewInt(16338636656),
				GasPrice:    *big.NewInt(16338636656),
				Hash:        "0x0f41f88706566a6e14bfc1f85d9761eb216a3b141c6eec13c820b6da569bf8a5",
				Input:       "Function: transfer(address recipient, uint256 amount)\n\nMethodID: 0xa9059cbb\n[0]:  0000000000000000000000006cc8dcbca746a6e4fdefb98e1d0df903b107fd21\n[1]:  00000000000000000000000000000000000000000000000000000081da747e22",
				Nonce:       0x0000000000000000,
				To:          "0x388c818ca8b9251b393131c08a736a67ccb19297",
			}},
			Subscriber: models.Subscriber{
				Address:              "",
				SubscribeBlockNumber: 5,
				SubscribeTxCount:     10,
			},
			IsErr: true,
		},
	}

	for _, testCase := range testCases {
		// we do not use t.Run() due to requirement of synchronous execution for every testcase
		// because we must get error in a last testcase due to duplicating registered address
		err := subscriberRepository.AddNewSubscriber(ctx, testCase.Subscriber)
		assert.NoError(t, err)
		err = subscriberRepository.AddTransactions(ctx, testCase.Name, testCase.Transactions)
		if testCase.IsErr {
			assert.Error(t, err)
		}
	}
}

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

			err := subscriberRepository.AddNewSubscriber(ctx, testCase.SubscriberForAddition)
			assert.NoError(t, err)

			gotSubscriber, err := subscriberRepository.GetSubscriberByAddress(ctx, testCase.SubscriberForAddition.Address)
			assert.NoError(t, err)

			assert.EqualValues(t, gotSubscriber, testCase.SubscriberForAddition)
		})
	}
}
