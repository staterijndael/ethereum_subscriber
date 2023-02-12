package redis_repository

import (
	"context"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"github.com/redis/go-redis/v9"
	"time"
)

// SubscriberRepository is a structure to store subscribers in a Redis database
type SubscriberRepository struct {
	redis          *redis.Client
	expirationTime time.Duration
}

// NewSubscriberRepository creates a new instance of SubscriberRepository
// redis - an instance of redis.Client for communication with the Redis database
// expirationTime - the expiration time for subscribers stored in the Redis database
func NewSubscriberRepository(redis *redis.Client, expirationTime time.Duration) *SubscriberRepository {
	return &SubscriberRepository{
		redis:          redis,
		expirationTime: expirationTime,
	}
}

// AddNewSubscriber adds a new subscriber to the Redis database
// ctx - a context for the Redis request
// subscriber - the subscriber to be added to the Redis database
// returns an error if adding the subscriber to the Redis database failed
func (r *SubscriberRepository) AddNewSubscriber(ctx context.Context, subscriber models.Subscriber) error {
	// Check if the subscriber already exists in the Redis database
	_, err := r.redis.Get(ctx, getSubscribersKey(subscriber.Address)).Result()
	if err == nil {
		// If the subscriber already exists, return an error
		return errors.New("subscriber already registered")
	}

	// If the error returned from Redis is not nil, return the error
	if err.Error() != redisNilErrMsg {
		return err
	}

	// Serialize the subscriber to be stored in the Redis database
	serializedSubscriber, err := serializeSubscribersValue(subscriber)
	if err != nil {
		return err
	}

	// Set the subscriber in the Redis database
	setSubCmd := r.redis.Set(ctx, getSubscribersKey(subscriber.Address), serializedSubscriber, r.expirationTime)
	if setSubCmd.Err() != nil {
		return setSubCmd.Err()
	}

	return nil
}

// GetSubscriberByAddress retrieves a subscriber from the Redis database based on their address
// ctx - a context for the Redis request
// address - the address of the subscriber to retrieve
// returns the subscriber and an error if retrieving the subscriber from the Redis database failed
func (r *SubscriberRepository) GetSubscriberByAddress(ctx context.Context, address string) (models.Subscriber, error) {
	// Get the subscriber data from redis using the provided address and the result of the getSubscribersKey function.
	subscriberRawData, err := r.redis.Get(ctx, getSubscribersKey(address)).Bytes()
	// If there was an error, check if it was due to the subscriber not being subscribed.
	if err != nil {
		if err.Error() == redisNilErrMsg {
			return models.Subscriber{}, errors.New("subscriber is not subscribed")
		}

		return models.Subscriber{}, err
	}

	// Deserialize the raw data into a Subscriber struct.
	subscriber, err := deserealizeSubscribersValue(subscriberRawData)
	if err != nil {
		return models.Subscriber{}, err
	}

	// Return the subscriber and nil for the error.
	return subscriber, nil
}
