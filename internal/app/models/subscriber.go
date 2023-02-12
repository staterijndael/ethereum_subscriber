package models

// Subscriber represent an address that subscribed for a listening and keeping data required by service
// for getting all address transaction
type Subscriber struct {
	// Subscriber address represented as a hexadecimal number in a string
	Address string
	// Block number that had subscriber in a moment of subscription or last parsed block number (depending on mode)
	SubscribeBlockNumber uint64
	// Number of transactions that user had in a moment of subscription or last parsed block (depending on mode)
	SubscribeTxCount uint64
}
