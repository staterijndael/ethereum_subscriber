package sync_greedy_redis_repository

import (
	"context"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	redis_driver "github.com/redis/go-redis/v9"
	"time"
)

type SubscriberRepository struct {
	redis          *redis_driver.Client
	expirationTime time.Duration
}

func NewSubscriberRepository(redis *redis_driver.Client, expirationTime time.Duration) *SubscriberRepository {
	return &SubscriberRepository{
		redis:          redis,
		expirationTime: expirationTime,
	}
}

func (r *SubscriberRepository) GetTransactionsReversed(ctx context.Context, address string) ([]*models.Transaction, error) {
	err := r.redis.Get(ctx, getSubscribersKey(address)).Err()
	if err != nil {
		if err.Error() == redisNilErrMsg {
			return nil, errors.New("address is not registered")
		}

		return nil, err
	}

	rawTxs, err := r.redis.Get(ctx, getSubscribersTxsKey(address)).Bytes()
	if err != nil {
		return nil, err
	}

	txs, err := deserializeSubscribersTxsValue(rawTxs)
	if err != nil {
		return nil, err
	}

	reversedSharedTransactions := models.ReverseTransactionsCopy(txs)

	return reversedSharedTransactions, nil
}

func (r *SubscriberRepository) GetLastTransaction(ctx context.Context, address string) (*models.Transaction, error) {
	err := r.redis.Get(ctx, getSubscribersKey(address)).Err()
	if err != nil {
		if err.Error() == redisNilErrMsg {
			return nil, errors.New("address is not registered")
		}

		return nil, err
	}

	rawTxs, err := r.redis.Get(ctx, getSubscribersTxsKey(address)).Bytes()
	if err != nil {
		if err.Error() == redisNilErrMsg {
			return nil, nil
		}

		return nil, err
	}

	txs, err := deserializeSubscribersTxsValue(rawTxs)
	if err != nil {
		return nil, err
	}
	if len(txs) == 0 {
		return nil, nil
	}

	return txs[len(txs)-1], nil
}

func (r *SubscriberRepository) AddTransactions(ctx context.Context, address string, txs []*models.Transaction) error {
	rawSubscriber, err := r.redis.Get(ctx, getSubscribersKey(address)).Bytes()
	if err != nil {
		if err.Error() == redisNilErrMsg {
			return errors.New("address is not registered")
		}

		return err
	}

	subscriber, err := deserealizeSubscribersValue(rawSubscriber)
	if err != nil {
		return err
	}

	rawTxs, err := r.redis.Get(ctx, getSubscribersTxsKey(address)).Bytes()
	if err != nil {
		if err.Error() == redisNilErrMsg {
			currentRawTxs, err2 := serializeSubscribersTxsValue(txs)
			if err2 != nil {
				return err2
			}
			err3 := r.redis.Set(ctx, getSubscribersTxsKey(address), currentRawTxs, r.expirationTime).Err()
			if err3 != nil {
				return err3
			}

			return nil
		}
		return err
	}

	storedTxs, err := deserializeSubscribersTxsValue(rawTxs)
	if err != nil {
		return err
	}

	storedTxs = append(storedTxs, txs...)

	var subscribeBlockNumber uint64
	if len(storedTxs) == 0 {
		subscribeBlockNumber = subscriber.SubscribeBlockNumber
	} else {
		subscribeBlockNumber = storedTxs[len(storedTxs)-1].BlockNumber
	}

	storedTxsRaw, err := serializeSubscribersTxsValue(storedTxs)
	if err != nil {
		return err
	}

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

func (r *SubscriberRepository) AddNewSubscriber(ctx context.Context, subscriber models.Subscriber) error {
	_, err := r.redis.Get(ctx, getSubscribersKey(subscriber.Address)).Result()
	if err == nil {
		return errors.New("subscriber already registered")
	}

	if err.Error() != redisNilErrMsg {
		return err
	}

	serializedSubscriber, err := serializeSubscribersValue(subscriber)
	if err != nil {
		return err
	}

	setSubCmd := r.redis.Set(ctx, getSubscribersKey(subscriber.Address), serializedSubscriber, r.expirationTime)
	if setSubCmd.Err() != nil {
		return setSubCmd.Err()
	}

	return nil
}

func (r *SubscriberRepository) GetSubscriberByAddress(ctx context.Context, address string) (models.Subscriber, error) {
	subscriberRawData, err := r.redis.Get(ctx, getSubscribersKey(address)).Bytes()
	if err != nil {
		if err.Error() == redisNilErrMsg {
			return models.Subscriber{}, errors.New("subscriber is not subscribed")
		}

		return models.Subscriber{}, err
	}

	subscriber, err := deserealizeSubscribersValue(subscriberRawData)
	if err != nil {
		return models.Subscriber{}, err
	}

	return subscriber, nil
}
