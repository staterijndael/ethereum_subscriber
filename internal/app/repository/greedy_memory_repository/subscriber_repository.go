package greedy_memory_repository

import (
	"context"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"sync"
)

type SubscriberRepository struct {
	subscribers   map[string]models.Subscriber
	subscriberTxs map[string][]*models.Transaction

	subscribersMx    sync.RWMutex
	subscribersTxsMx sync.RWMutex
}

func NewSubscriberRepository() *SubscriberRepository {
	return &SubscriberRepository{
		subscribers:      make(map[string]models.Subscriber),
		subscriberTxs:    make(map[string][]*models.Transaction),
		subscribersMx:    sync.RWMutex{},
		subscribersTxsMx: sync.RWMutex{},
	}
}

func (r *SubscriberRepository) GetTransactionsReversed(ctx context.Context, address string) ([]*models.Transaction, error) {
	r.subscribersMx.RLock()
	if _, ok := r.subscribers[address]; !ok {
		r.subscribersMx.RUnlock()
		return nil, errors.New("address is not registered")
	}
	r.subscribersMx.RUnlock()

	r.subscribersTxsMx.RLock()
	reversedSharedTransactions := models.ReverseTransactionsCopy(r.subscriberTxs[address])
	r.subscribersTxsMx.RUnlock()

	return reversedSharedTransactions, nil
}

func (r *SubscriberRepository) GetLastTransaction(ctx context.Context, address string) (*models.Transaction, error) {
	r.subscribersMx.RLock()
	if _, ok := r.subscribers[address]; !ok {
		r.subscribersMx.RUnlock()
		return nil, errors.New("address is not registered")
	}
	r.subscribersMx.RUnlock()

	r.subscribersTxsMx.RLock()
	txs := r.subscriberTxs[address]
	if len(txs) == 0 {
		r.subscribersTxsMx.RUnlock()
		return nil, nil
	}

	lastTx := txs[len(txs)-1]
	r.subscribersTxsMx.RUnlock()

	return lastTx, nil
}

func (r *SubscriberRepository) AddTransactions(ctx context.Context, address string, txs []*models.Transaction) error {
	r.subscribersMx.RLock()
	subscriber, ok := r.subscribers[address]
	if !ok {
		r.subscribersMx.RUnlock()
		return errors.New("address is not registered")
	}
	r.subscribersMx.RUnlock()

	r.subscribersTxsMx.Lock()
	r.subscriberTxs[address] = append(r.subscriberTxs[address], txs...)

	internalTxs := r.subscriberTxs[address]
	r.subscribersTxsMx.Unlock()

	var subscribeBlockNumber uint64
	if len(internalTxs) == 0 {
		subscribeBlockNumber = subscriber.SubscribeBlockNumber
	} else {
		subscribeBlockNumber = internalTxs[len(r.subscriberTxs)-1].BlockNumber
	}

	r.subscribersMx.Lock()
	r.subscribers[address] = models.Subscriber{
		Address:              address,
		SubscribeBlockNumber: subscribeBlockNumber,
		SubscribeTxCount:     subscriber.SubscribeTxCount + uint64(len(txs)),
	}
	r.subscribersMx.Unlock()

	return nil
}

func (r *SubscriberRepository) AddNewSubscriber(ctx context.Context, subscriber models.Subscriber) error {
	r.subscribersMx.Lock()
	if _, ok := r.subscribers[subscriber.Address]; ok {
		r.subscribersMx.Unlock()
		return errors.New("subscriber already registered")
	}

	r.subscribers[subscriber.Address] = subscriber
	r.subscriberTxs[subscriber.Address] = []*models.Transaction{}

	r.subscribersMx.Unlock()

	return nil
}

func (r *SubscriberRepository) GetSubscriberByAddress(ctx context.Context, address string) (models.Subscriber, error) {
	r.subscribersMx.RLock()
	subscriber, ok := r.subscribers[address]
	if !ok {
		r.subscribersMx.RUnlock()
		return models.Subscriber{}, errors.New("address is not subscribed")
	}
	r.subscribersMx.RUnlock()

	return subscriber, nil
}
