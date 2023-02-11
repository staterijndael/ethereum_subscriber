package redis_repository

import (
	"encoding/json"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
)

const redisNilErrMsg = "redis: nil"

const currentBlockKey = "sync_key_currentBlock"
const subscribersKey = "sync_key_Subscriber-"

func getCurrentBlockKey() string {
	return currentBlockKey
}

func getSubscribersKey(address string) string {
	return subscribersKey + address
}

func serializeCurrentBlockValue(currentBlock uint64) uint64 {
	return currentBlock
}

func serializeSubscribersValue(subscriber models.Subscriber) ([]byte, error) {
	rawData, err := json.Marshal(subscriber)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func deserealizeSubscribersValue(rawSubscriber []byte) (models.Subscriber, error) {
	var subscriber models.Subscriber
	err := json.Unmarshal(rawSubscriber, &subscriber)
	if err != nil {
		return models.Subscriber{}, err
	}

	return subscriber, nil
}
