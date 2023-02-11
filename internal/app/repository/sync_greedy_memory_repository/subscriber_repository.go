package sync_greedy_memory_repository

import (
	"context"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
)

type SubscriberRepository struct {
	subscribers   map[string]models.Subscriber
	subscriberTxs map[string][]*models.Transaction
}

func NewSubscriberRepository() *SubscriberRepository {
	return &SubscriberRepository{
		subscribers:   make(map[string]models.Subscriber),
		subscriberTxs: make(map[string][]*models.Transaction),
	}
}

func (r *SubscriberRepository) GetTransactionsReversed(ctx context.Context, address string) ([]*models.Transaction, error) {
	if _, ok := r.subscribers[address]; !ok {
		return nil, errors.New("address is not registered")
	}

	reversedSharedTransactions := models.ReverseTransactionsCopy(r.subscriberTxs[address])

	return reversedSharedTransactions, nil
}

func (r *SubscriberRepository) GetLastTransaction(ctx context.Context, address string) (*models.Transaction, error) {
	if _, ok := r.subscribers[address]; !ok {
		return nil, errors.New("address is not registered")
	}

	txs := r.subscriberTxs[address]
	if len(txs) == 0 {
		return nil, nil
	}

	return txs[len(txs)-1], nil
}

func (r *SubscriberRepository) AddTransactions(ctx context.Context, address string, txs []*models.Transaction) error {
	subscriber, ok := r.subscribers[address]
	if !ok {
		return errors.New("address is not registered")
	}

	r.subscriberTxs[address] = append(r.subscriberTxs[address], txs...)

	internalTxs := r.subscriberTxs[address]

	var subscribeBlockNumber uint64
	if len(internalTxs) == 0 {
		subscribeBlockNumber = subscriber.SubscribeBlockNumber
	} else {
		subscribeBlockNumber = internalTxs[len(r.subscriberTxs)-1].BlockNumber
	}

	r.subscribers[address] = models.Subscriber{
		Address:              address,
		SubscribeBlockNumber: subscribeBlockNumber,
		SubscribeTxCount:     subscriber.SubscribeTxCount + uint64(len(txs)),
	}

	return nil
}

func (r *SubscriberRepository) AddNewSubscriber(ctx context.Context, subscriber models.Subscriber) error {
	if _, ok := r.subscribers[subscriber.Address]; ok {
		return errors.New("subscriber already registered")
	}

	r.subscribers[subscriber.Address] = subscriber
	r.subscriberTxs[subscriber.Address] = []*models.Transaction{}

	return nil
}

func (r *SubscriberRepository) GetSubscriberByAddress(ctx context.Context, address string) (models.Subscriber, error) {
	subscriber, ok := r.subscribers[address]
	if !ok {
		return models.Subscriber{}, errors.New("address is not subscribed")
	}

	return subscriber, nil
}
