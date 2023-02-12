package greedy_redis_repository

import (
	"encoding/json"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
)

// redisNilErrMsg is a constant that holds the value "redis: nil" and is used to check if the error returned by Redis client is "nil".
const redisNilErrMsg = "redis: nil"

// currentBlockKey is a constant that holds the key for the current block in the Redis database.
const currentBlockKey = "sync_greedy_key_currentBlock"

// subscribersKey is a constant that holds the prefix for the subscribers' keys in the Redis database.
const subscribersKey = "sync_greedy_key_Subscriber-"

// subscribersTxsKey is a constant that holds the prefix for the subscribers' transactions keys in the Redis database.
const subscribersTxsKey = "sync_greedy_key_SubscriberTxs-"

// getCurrentBlockKey returns the key for the current block in the Redis database.
func getCurrentBlockKey() string {
	return currentBlockKey
}

// getSubscribersTxsKey returns the key for the transactions of a subscriber in the Redis database.
func getSubscribersTxsKey(address string) string {
	return subscribersTxsKey + address
}

// serializeSubscribersTxsValue serializes an array of transactions into a JSON byte array.
func serializeSubscribersTxsValue(txs []*models.Transaction) ([]byte, error) {
	return json.Marshal(txs)
}

// deserializeSubscribersTxsValue deserializes a JSON byte array into an array of transactions.
func deserializeSubscribersTxsValue(rawTxs []byte) ([]*models.Transaction, error) {
	var txs []*models.Transaction
	err := json.Unmarshal(rawTxs, &txs)
	if err != nil {
		return nil, err
	}

	return txs, nil
}

// getSubscribersKey returns the key for a subscriber in the Redis database.
func getSubscribersKey(address string) string {
	return subscribersKey + address
}

// serializeCurrentBlockValue serializes a uint64 into a uint64 value.
func serializeCurrentBlockValue(currentBlock uint64) uint64 {
	return currentBlock
}

// serializeSubscribersValue serializes a Subscriber struct into a JSON byte array.
func serializeSubscribersValue(subscriber models.Subscriber) ([]byte, error) {
	rawData, err := json.Marshal(subscriber)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

// deserealizeSubscribersValue deserializes a JSON byte array into a Subscriber struct.
func deserealizeSubscribersValue(rawSubscriber []byte) (models.Subscriber, error) {
	var subscriber models.Subscriber
	err := json.Unmarshal(rawSubscriber, &subscriber)
	if err != nil {
		return models.Subscriber{}, err
	}

	return subscriber, nil
}
