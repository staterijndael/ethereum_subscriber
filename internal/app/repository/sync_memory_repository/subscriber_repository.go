package sync_memory_repository

import (
	"context"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
)

type SubscriberRepository struct {
	subscribers map[string]models.Subscriber
}

func NewSubscriberRepository() *SubscriberRepository {
	return &SubscriberRepository{
		subscribers: make(map[string]models.Subscriber),
	}
}

func (r *SubscriberRepository) AddNewSubscriber(ctx context.Context, subscriber models.Subscriber) error {
	if _, ok := r.subscribers[subscriber.Address]; ok {
		return errors.New("subscriber already registered")
	}

	r.subscribers[subscriber.Address] = subscriber

	return nil
}

func (r *SubscriberRepository) GetSubscriberByAddress(ctx context.Context, address string) (models.Subscriber, error) {
	subscriber, ok := r.subscribers[address]
	if !ok {
		return models.Subscriber{}, errors.New("address is not subscribed")
	}

	return subscriber, nil
}
