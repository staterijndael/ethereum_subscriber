package greedy_memory_repository

import (
	"context"
	"errors"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"sync"
)

// SubscriberRepository is a struct that holds the map of subscribers and their transactions
type SubscriberRepository struct {
	subscribers   map[string]models.Subscriber
	subscriberTxs map[string][]*models.Transaction

	// Mutexes to ensure concurrency safety while accessing subscribers and subscriber transactions maps
	subscribersMx    sync.RWMutex
	subscribersTxsMx sync.RWMutex
}

// NewSubscriberRepository returns a new instance of the SubscriberRepository struct
func NewSubscriberRepository() *SubscriberRepository {
	return &SubscriberRepository{
		subscribers:      make(map[string]models.Subscriber),
		subscriberTxs:    make(map[string][]*models.Transaction),
		subscribersMx:    sync.RWMutex{},
		subscribersTxsMx: sync.RWMutex{},
	}
}

// GetTransactionsReversed returns the reversed transactions of a subscriber by address
func (r *SubscriberRepository) GetTransactionsReversed(ctx context.Context, address string) ([]*models.Transaction, error) {
	// Lock the subscribers map for reading to ensure concurrency safety
	r.subscribersMx.RLock()
	if _, ok := r.subscribers[address]; !ok {
		r.subscribersMx.RUnlock()
		return nil, errors.New("address is not registered")
	}
	r.subscribersMx.RUnlock()

	// Lock the subscribers transactions map for reading to ensure concurrency safety
	r.subscribersTxsMx.RLock()
	reversedSharedTransactions := models.ReverseTransactionsCopy(r.subscriberTxs[address])
	r.subscribersTxsMx.RUnlock()

	return reversedSharedTransactions, nil
}

// GetLastTransaction returns the last transaction of a subscriber by address
func (r *SubscriberRepository) GetLastTransaction(ctx context.Context, address string) (*models.Transaction, error) {
	// Lock the subscribers map for reading to ensure concurrency safety
	r.subscribersMx.RLock()
	if _, ok := r.subscribers[address]; !ok {
		r.subscribersMx.RUnlock()
		return nil, errors.New("address is not registered")
	}
	r.subscribersMx.RUnlock()

	// Lock the subscribers transactions map for reading to ensure concurrency safety
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

// AddTransactions adds transactions to the subscriber with the given address.
// If the address is not registered, it returns an error.
func (r *SubscriberRepository) AddTransactions(ctx context.Context, address string, txs []*models.Transaction) error {
	// Read lock on subscribers to safely access the map
	r.subscribersMx.RLock()
	subscriber, ok := r.subscribers[address]
	if !ok {
		// Unlock the map before returning
		r.subscribersMx.RUnlock()
		return errors.New("address is not registered")
	}
	r.subscribersMx.RUnlock()

	// Write lock on subscriberTxs to safely access and modify the map
	r.subscribersTxsMx.Lock()
	r.subscriberTxs[address] = append(r.subscriberTxs[address], txs...)

	internalTxs := r.subscriberTxs[address]
	r.subscribersTxsMx.Unlock()

	var subscribeBlockNumber uint64
	if len(internalTxs) == 0 {
		subscribeBlockNumber = subscriber.SubscribeBlockNumber
	} else {
		subscribeBlockNumber = internalTxs[len(internalTxs)-1].BlockNumber
	}

	// Write lock on subscribers to safely access and modify the map
	r.subscribersMx.Lock()
	r.subscribers[address] = models.Subscriber{
		Address:              address,
		SubscribeBlockNumber: subscribeBlockNumber,
		SubscribeTxCount:     subscriber.SubscribeTxCount + uint64(len(txs)),
	}
	r.subscribersMx.Unlock()

	return nil
}

// AddNewSubscriber adds a new subscriber to the repository.
// If the subscriber already exists, it returns an error.
// The function uses a lock to ensure thread-safety while accessing the `subscribers` map.
func (r *SubscriberRepository) AddNewSubscriber(ctx context.Context, subscriber models.Subscriber) error {
	// Acquire the lock for the subscribers map
	r.subscribersMx.Lock()

	// Check if the subscriber already exists
	if _, ok := r.subscribers[subscriber.Address]; ok {
		// Release the lock
		r.subscribersMx.Unlock()
		// Return an error if the subscriber already exists
		return errors.New("subscriber already registered")
	}

	// Add the new subscriber to the subscribers map
	r.subscribers[subscriber.Address] = subscriber
	// Initialize the transactions slice for the subscriber
	r.subscriberTxs[subscriber.Address] = []*models.Transaction{}

	// Release the lock
	r.subscribersMx.Unlock()

	// Return nil to indicate success
	return nil
}

// GetSubscriberByAddress returns the subscriber with the specified address.
// If the address is not subscribed, it returns an error.
// The function uses a read lock to ensure thread-safety while accessing the `subscribers` map.
func (r *SubscriberRepository) GetSubscriberByAddress(ctx context.Context, address string) (models.Subscriber, error) {
	// Acquire the read lock for the subscribers map
	r.subscribersMx.RLock()

	// Get the subscriber with the specified address
	subscriber, ok := r.subscribers[address]
	if !ok {
		// Release the read lock
		r.subscribersMx.RUnlock()
		// Return an error if the address is not subscribed
		return models.Subscriber{}, errors.New("address is not subscribed")
	}

	// Release the read lock
	r.subscribersMx.RUnlock()

	// Return the subscriber and nil to indicate success
	return subscriber, nil
}
