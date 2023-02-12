package redis_repository

import (
	"encoding/json"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
)

// redisNilErrMsg is a constant string representing the error message for a nil value in redis
const redisNilErrMsg = "redis: nil"

// currentBlockKey is a constant string representing the key for the current block in redis
const currentBlockKey = "sync_key_currentBlock"

// subscribersKey is a constant string representing the prefix for subscribers' keys in redis
const subscribersKey = "sync_key_Subscriber-"

// getCurrentBlockKey returns the key for the current block
func getCurrentBlockKey() string {
	return currentBlockKey
}

// getSubscribersKey returns the key for a subscriber with the given address
func getSubscribersKey(address string) string {
	return subscribersKey + address
}

// serializeCurrentBlockValue serializes the current block value as a uint64
func serializeCurrentBlockValue(currentBlock uint64) uint64 {
	return currentBlock
}

// serializeSubscribersValue serializes a subscriber as a byte slice
func serializeSubscribersValue(subscriber models.Subscriber) ([]byte, error) {
	rawData, err := json.Marshal(subscriber)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

// deserealizeSubscribersValue deserializes a subscriber from a byte slice
func deserealizeSubscribersValue(rawSubscriber []byte) (models.Subscriber, error) {
	var subscriber models.Subscriber
	err := json.Unmarshal(rawSubscriber, &subscriber)
	if err != nil {
		return models.Subscriber{}, err
	}

	return subscriber, nil
}
