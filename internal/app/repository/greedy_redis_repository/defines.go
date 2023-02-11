package greedy_redis_repository

import (
	"encoding/json"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
)

const redisNilErrMsg = "redis: nil"

const currentBlockKey = "sync_greedy_key_currentBlock"
const subscribersKey = "sync_greedy_key_Subscriber-"
const subscribersTxsKey = "sync_greedy_key_SubscriberTxs-"

func getCurrentBlockKey() string {
	return currentBlockKey
}

func getSubscribersTxsKey(address string) string {
	return subscribersTxsKey + address
}

func serializeSubscribersTxsValue(txs []*models.Transaction) ([]byte, error) {
	return json.Marshal(txs)
}

func deserializeSubscribersTxsValue(rawTxs []byte) ([]*models.Transaction, error) {
	var txs []*models.Transaction
	err := json.Unmarshal(rawTxs, &txs)
	if err != nil {
		return nil, err
	}

	return txs, nil
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
