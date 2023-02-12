package greedy_redis_repository

import (
	"context"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	redis_driver "github.com/redis/go-redis/v9"
	"time"
)

// SubscriberRepository is a struct that represents a repository for subscriber data. It holds a Redis client instance and an expiration time for data stored in the Redis cache.
type SubscriberRepository struct {
	redis          *redis_driver.Client
	expirationTime time.Duration
}

// NewSubscriberRepository is a constructor for SubscriberRepository that takes in a Redis client instance and an expiration time for data stored in the Redis cache, and returns a pointer to a SubscriberRepository instance.
func NewSubscriberRepository(redis *redis_driver.Client, expirationTime time.Duration) *SubscriberRepository {
	return &SubscriberRepository{
		redis:          redis,
		expirationTime: expirationTime,
	}
}

// GetTransactionsReversed retrieves the transactions associated with a subscriber, reversed, from the Redis cache. It takes in a context and the subscriber's address, and returns a slice of Transaction instances and an error if one occurred.
func (r *SubscriberRepository) GetTransactionsReversed(ctx context.Context, address string) ([]*models.Transaction, error) {
	// Check if the subscriber exists in the cache
	err := r.redis.Get(ctx, getSubscribersKey(address)).Err()
	if err != nil {
		// If the subscriber does not exist, return an error indicating that the address is not registered
		if err.Error() == redisNilErrMsg {
			return nil, errors.New("address is not registered")
		}
		// If there was a different error, return it
		return nil, err
	}

	// Get the raw byte representation of the subscriber's transactions from the cache
	rawTxs, err := r.redis.Get(ctx, getSubscribersTxsKey(address)).Bytes()
	if err != nil {
		// If there was an error retrieving the data, return it
		return nil, err
	}

	// Deserialize the raw transactions into a slice of Transaction instances
	txs, err := deserializeSubscribersTxsValue(rawTxs)
	if err != nil {
		// If there was an error deserializing the data, return it
		return nil, err
	}

	// Reverse the list of transactions and return it
	reversedSharedTransactions := models.ReverseTransactionsCopy(txs)
	return reversedSharedTransactions, nil
}

// GetLastTransaction returns the last transaction for a given address from the Redis cache.
// If the address is not registered or if there are no transactions associated with the address, it returns nil and an error.
func (r *SubscriberRepository) GetLastTransaction(ctx context.Context, address string) (*models.Transaction, error) {
	// Check if the address is registered in the Redis cache
	err := r.redis.Get(ctx, getSubscribersKey(address)).Err()
	if err != nil {
		// If the error message is "redis: nil", then the address is not registered
		if err.Error() == redisNilErrMsg {
			return nil, errors.New("address is not registered")
		}

		// Return the error if it is not "redis: nil"
		return nil, err
	}

	// Get the transactions associated with the address from the Redis cache
	rawTxs, err := r.redis.Get(ctx, getSubscribersTxsKey(address)).Bytes()
	if err != nil {
		// If the error message is "redis: nil", then there are no transactions associated with the address
		if err.Error() == redisNilErrMsg {
			return nil, nil
		}

		// Return the error if it is not "redis: nil"
		return nil, err
	}

	// Deserialize the transactions
	txs, err := deserializeSubscribersTxsValue(rawTxs)
	if err != nil {
		return nil, err
	}

	// Return nil if there are no transactions
	if len(txs) == 0 {
		return nil, nil
	}

	// Return the last transaction
	return txs[len(txs)-1], nil
}

// AddTransactions adds a list of transactions for a given address to the Redis cache.
// If the address is not registered, it returns an error.
func (r *SubscriberRepository) AddTransactions(ctx context.Context, address string, txs []*models.Transaction) error {
	// Get the raw byte data of the subscriber from Redis using the given address
	rawSubscriber, err := r.redis.Get(ctx, getSubscribersKey(address)).Bytes()
	if err != nil {
		// If the error is a "nil" error, then the address is not registered
		if err.Error() == redisNilErrMsg {
			return errors.New("address is not registered")
		}

		// Return any other error encountered
		return err
	}

	// Deserialize the raw byte data of the subscriber into a subscriber object
	subscriber, err := deserealizeSubscribersValue(rawSubscriber)
	if err != nil {
		return err
	}

	// Get the raw byte data of the transactions for the given address
	rawTxs, err := r.redis.Get(ctx, getSubscribersTxsKey(address)).Bytes()
	if err != nil {
		// If the error is a "nil" error, then there are no stored transactions for this address
		if err.Error() == redisNilErrMsg {
			// Serialize the current transactions
			currentRawTxs, err2 := serializeSubscribersTxsValue(txs)
			if err2 != nil {
				return err2
			}
			// Store the serialized transactions in Redis
			err3 := r.redis.Set(ctx, getSubscribersTxsKey(address), currentRawTxs, r.expirationTime).Err()
			if err3 != nil {
				return err3
			}

			return nil
		}
		return err
	}

	// Deserialize the raw byte data of the stored transactions into a transactions slice
	storedTxs, err := deserializeSubscribersTxsValue(rawTxs)
	if err != nil {
		return err
	}

	// Append the current transactions to the stored transactions
	storedTxs = append(storedTxs, txs...)

	// Determine the subscribe block number
	var subscribeBlockNumber uint64
	if len(storedTxs) == 0 {
		// If there are no stored transactions, use the subscriber's original subscribe block number
		subscribeBlockNumber = subscriber.SubscribeBlockNumber
	} else {
		// If there are stored transactions, use the block number of the last transaction
		subscribeBlockNumber = storedTxs[len(storedTxs)-1].BlockNumber
	}

	// Serialize the updated stored transactions
	storedTxsRaw, err := serializeSubscribersTxsValue(storedTxs)
	if err != nil {
		return err
	}

	// Store the serialized updated stored transactions in Redis
	err = r.redis.Set(ctx, getSubscribersTxsKey(address), storedTxsRaw, r.expirationTime).Err()
	if err != nil {
		return err
	}

	newSubscriber := models.Subscriber{
		Address:              address,
		SubscribeBlockNumber: subscribeBlockNumber,
		SubscribeTxCount:     subscriber.SubscribeTxCount + uint64(len(txs)),
	}

	rawNewSubscriber, err := serializeSubscribersValue(newSubscriber)
	if err != nil {
		return err
	}

	err = r.redis.Set(ctx, getSubscribersKey(address), rawNewSubscriber, r.expirationTime).Err()
	if err != nil {
		return err
	}

	return nil
}

// AddNewSubscriber adds a new subscriber to the repository.
// If the subscriber with the same address already exists in the repository, an error "subscriber already registered" will be returned.
func (r *SubscriberRepository) AddNewSubscriber(ctx context.Context, subscriber models.Subscriber) error {
	// Check if subscriber with the same address already exists
	_, err := r.redis.Get(ctx, getSubscribersKey(subscriber.Address)).Result()
	if err == nil {
		return errors.New("subscriber already registered")
	}

	// Return an error if it's not a nil error message
	if err.Error() != redisNilErrMsg {
		return err
	}

	// Serialize the subscriber data
	serializedSubscriber, err := serializeSubscribersValue(subscriber)
	if err != nil {
		return err
	}

	// Store the serialized subscriber data in the repository with the specified expiration time
	setSubCmd := r.redis.Set(ctx, getSubscribersKey(subscriber.Address), serializedSubscriber, r.expirationTime)
	if setSubCmd.Err() != nil {
		return setSubCmd.Err()
	}

	return nil
}

// GetSubscriberByAddress returns a subscriber with a given address from the repository.
// If the subscriber with the given address doesn't exist in the repository, an error "subscriber is not subscribed" will be returned.
func (r *SubscriberRepository) GetSubscriberByAddress(ctx context.Context, address string) (models.Subscriber, error) {
	// Get the raw subscriber data from the repository
	subscriberRawData, err := r.redis.Get(ctx, getSubscribersKey(address)).Bytes()
	if err != nil {
		// If the subscriber with the given address doesn't exist in the repository, return an error "subscriber is not subscribed"
		if err.Error() == redisNilErrMsg {
			return models.Subscriber{}, errors.New("subscriber is not subscribed")
		}

		return models.Subscriber{}, err
	}

	// Deserialize the subscriber data
	subscriber, err := deserealizeSubscribersValue(subscriberRawData)
	if err != nil {
		return models.Subscriber{}, err
	}

	return subscriber, nil
}
