package sync_memory_repository

import (
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

func (r *SubscriberRepository) AddNewSubscriber(subscriber models.Subscriber) error {
	if _, ok := r.subscribers[subscriber.Address]; ok {
		return errors.New("subscriber already registered")
	}

	r.subscribers[subscriber.Address] = subscriber

	return nil
}

func (r *SubscriberRepository) GetSubscriberByAddress(address string) (models.Subscriber, error) {
	subscriber, ok := r.subscribers[address]
	if !ok {
		return models.Subscriber{}, errors.New("address is not subscribed")
	}

	return subscriber, nil
}
