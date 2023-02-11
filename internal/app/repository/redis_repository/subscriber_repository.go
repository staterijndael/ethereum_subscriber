package redis_repository

import (
	"context"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"github.com/redis/go-redis/v9"
	"time"
)

type SubscriberRepository struct {
	redis          *redis.Client
	expirationTime time.Duration
}

func NewSubscriberRepository(redis *redis.Client, expirationTime time.Duration) *SubscriberRepository {
	return &SubscriberRepository{
		redis:          redis,
		expirationTime: expirationTime,
	}
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
