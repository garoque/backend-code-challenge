package entity

import (
	"time"

	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/garoque/backend-code-challenge-snapfi/pkg/uuid"
)

type StatesTransaction int

const (
	OPEN StatesTransaction = iota
	BOOKED
	FAILED
)

var StatesTransactionString = []string{
	"OPEN", "BOOKED", "FAILED",
}

func (st StatesTransaction) String() string {
	return StatesTransactionString[st]
}

type Transaction struct {
	ID            string            `json:"id"`
	SourceId      string            `json:"senderId" db:"id_source"`
	DestinationId string            `json:"receiverId" db:"id_destination"`
	Amount        float64           `json:"amount"`
	State         StatesTransaction `json:"-" db:"state"`
	StateString   string            `json:"state,omitempty"`
	CreatedAt     *time.Time        `json:"createdAt" db:"created_at"`
}

func NewTransaction(tr dto.CreateTransaction) *Transaction {
	return &Transaction{
		ID:            uuid.NewId(),
		SourceId:      tr.SourceUserId,
		DestinationId: tr.DestinationUserId,
		Amount:        tr.Amount,
	}
}

type TransactionIncreaseBalanceUser struct {
	ID     string  `json:"-"`
	UserId string  `json:"userId"`
	Value  float64 `json:"value"`
}

func NewIncreaseBalanceUser(tr dto.IncreaseBalanceUser) *TransactionIncreaseBalanceUser {
	return &TransactionIncreaseBalanceUser{
		ID:     uuid.NewId(),
		UserId: tr.UserId,
		Value:  tr.Value,
	}
}
