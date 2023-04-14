package entity

import (
	"testing"

	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	transaction := NewTransaction(dto.CreateTransaction{
		SourceUserId:      "source-user-id",
		DestinationUserId: "destination-user-id",
		Amount:            100.10,
	})
	assert.NotNil(t, transaction)
	assert.NotEmpty(t, transaction.ID)
	assert.Equal(t, 100.10, transaction.Amount)
	assert.Equal(t, "source-user-id", transaction.SourceId)
	assert.Equal(t, "destination-user-id", transaction.DestinationId)
}

func TestNewIncreaseBalanceUser(t *testing.T) {
	transaction := NewIncreaseBalanceUser(dto.IncreaseBalanceUser{
		UserId: "user-id",
		Value:  100.10,
	})
	assert.NotNil(t, transaction)
	assert.NotEmpty(t, transaction.ID)
	assert.Equal(t, 100.10, transaction.Value)
	assert.Equal(t, "user-id", transaction.UserId)
}

func TestStatesTransactionString(t *testing.T) {
	stateOpened := OPEN.String()
	assert.Equal(t, "OPEN", stateOpened)

	stateBooked := BOOKED.String()
	assert.Equal(t, "BOOKED", stateBooked)

	stateFailed := FAILED.String()
	assert.Equal(t, "FAILED", stateFailed)
}
