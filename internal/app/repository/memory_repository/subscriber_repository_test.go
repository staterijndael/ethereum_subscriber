package memory_repository

import (
	"context"
	"github.com/bluntenpassant/ethereum_subscriber/internal/app/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubscriberRepository_AddNewSubscriber(t *testing.T) {
	ctx := context.TODO()

	subscriberRepository := NewSubscriberRepository()
	err := subscriberRepository.AddNewSubscriber(ctx, models.Subscriber{
		Address:              "0x690B9A9E9aa1C9dB991C7721a92d351Db4FaC990",
		SubscribeBlockNumber: 5,
		SubscribeTxCount:     30,
	})
	assert.NoError(t, err)
}

func TestSubscriberRepository_GetSubscriberByAddress(t *testing.T) {
	ctx := context.TODO()

	subscriberRepository := NewSubscriberRepository()

	subscriber := models.Subscriber{
		Address:              "0x690B9A9E9aa1C9dB991C7721a92d351Db4FaC990",
		SubscribeBlockNumber: 5,
		SubscribeTxCount:     30,
	}

	err := subscriberRepository.AddNewSubscriber(ctx, subscriber)
	assert.NoError(t, err)

	gotSubscriber, err := subscriberRepository.GetSubscriberByAddress(ctx, subscriber.Address)
	assert.NoError(t, err)

	assert.EqualValues(t, gotSubscriber, subscriber)
}
