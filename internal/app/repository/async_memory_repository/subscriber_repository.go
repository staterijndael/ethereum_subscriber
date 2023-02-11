package async_memory_repository

import (
	"context"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"sync"
)

type SubscriberRepository struct {
	subscribers   map[string]models.Subscriber
	subscribersMx sync.RWMutex
}

func NewSubscriberRepository() *SubscriberRepository {
	return &SubscriberRepository{
		subscribers:   make(map[string]models.Subscriber),
		subscribersMx: sync.RWMutex{},
	}
}

func (r *SubscriberRepository) AddNewSubscriber(ctx context.Context, subscriber models.Subscriber) error {
	r.subscribersMx.Lock()
	defer r.subscribersMx.Unlock()

	if _, ok := r.subscribers[subscriber.Address]; ok {
		return errors.New("subscriber already registered")
	}

	r.subscribers[subscriber.Address] = subscriber

	return nil
}

func (r *SubscriberRepository) GetSubscriberByAddress(ctx context.Context, address string) (models.Subscriber, error) {
	r.subscribersMx.RLock()
	defer r.subscribersMx.RUnlock()

	subscriber, ok := r.subscribers[address]
	if !ok {
		return models.Subscriber{}, errors.New("address is not subscribed")
	}

	return subscriber, nil
}
