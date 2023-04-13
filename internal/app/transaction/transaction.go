package transaction

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/garoque/backend-code-challenge-snapfi/internal/database"
	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	"github.com/garoque/backend-code-challenge-snapfi/pkg/custom_err"
)

type AppTransactionInterface interface {
	Create(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error)
	IncreaseBalanceUser(ctx context.Context, transaction *entity.TransactionIncreaseBalanceUser) (float64, error)
}

type appTransactionImpl struct {
	db *database.Container
}

func NewAppTransaction(db *database.Container) AppTransactionInterface {
	return &appTransactionImpl{db}
}

func (tr *appTransactionImpl) Create(ctx context.Context, transaction *entity.Transaction) (*entity.Transaction, error) {
	err := tr.db.Transaction.Create(ctx, transaction)
	if err != nil {
		log.Println("Error app.Transaction.Create.db.Create: ", err.Error())
		return nil, err
	}

	sourceUser, err := tr.db.User.ReadOneById(ctx, transaction.SourceId)
	if err != nil {
		tr.updateStatusTransaction(ctx, transaction, 2)

		log.Println("Error app.Transaction.Create.db.ReadOneById.sourceId: ", err.Error())
		return nil, err
	}

	destinationUser, err := tr.db.User.ReadOneById(ctx, transaction.DestinationId)
	if err != nil {
		tr.updateStatusTransaction(ctx, transaction, 2)

		log.Println("Error app.Transaction.Create.db.ReadOneById.destinationId: ", err.Error())
		return nil, err
	}

	sourceUser.Mutex = &sync.Mutex{}
	destinationUser.Mutex = &sync.Mutex{}

	sourceUser.Mutex.Lock()
	destinationUser.Mutex.Lock()
	defer sourceUser.Mutex.Unlock()
	defer destinationUser.Mutex.Unlock()

	if sourceUser.Balance < transaction.Amount {
		tr.updateStatusTransaction(ctx, transaction, 2)

		log.Println("Error app.Transaction.Create sourceUser.Balance < transaction.Amount Insufficient balance")
		return nil, custom_err.New(http.StatusInternalServerError, "Insufficient balance")
	}

	sourceUser.Balance -= transaction.Amount
	err = tr.db.Transaction.UpdateBalanceUser(ctx, sourceUser.ID, sourceUser.Balance)
	if err != nil {
		tr.revertSourceBalanceTransaction(ctx, sourceUser, transaction.Amount)
		tr.updateStatusTransaction(ctx, transaction, 2)

		log.Println("Error app.Transaction.Create.db.UpdateBalanceUser.sourceUser: ", err.Error())
		return nil, err
	}

	destinationUser.Balance += transaction.Amount
	err = tr.db.Transaction.UpdateBalanceUser(ctx, destinationUser.ID, destinationUser.Balance)
	if err != nil {
		tr.revertDestinationBalanceTransaction(ctx, destinationUser, transaction.Amount)
		tr.updateStatusTransaction(ctx, transaction, 2)

		log.Println("Error app.Transaction.Create.db.UpdateBalanceUser.destinationUser: ", err.Error())
		return transaction, err
	}

	tr.updateStatusTransaction(ctx, transaction, 1)

	return transaction, nil
}

func (tr *appTransactionImpl) updateStatusTransaction(ctx context.Context, transaction *entity.Transaction, state entity.StatesTransaction) {
	transaction.State = state
	transaction.StateString = transaction.State.String()
	tr.db.Transaction.UpdateState(ctx, transaction.StateString, transaction.ID)
}

func (tr *appTransactionImpl) revertSourceBalanceTransaction(ctx context.Context, user *entity.User, amount float64) {
	user.Balance += amount
	if err := tr.db.Transaction.UpdateBalanceUser(ctx, user.ID, user.Balance); err != nil {
		log.Println("Error app.Transaction.Create.db.UpdateBalanceUser.revertSourceBalanceTransaction: ", err.Error())
	}
}

func (tr *appTransactionImpl) revertDestinationBalanceTransaction(ctx context.Context, user *entity.User, amount float64) {
	user.Balance -= amount
	if err := tr.db.Transaction.UpdateBalanceUser(ctx, user.ID, user.Balance); err != nil {
		log.Println("Error app.Transaction.Create.db.UpdateBalanceUser.revertDestinationBalanceTransaction: ", err.Error())
	}
}

func (tr *appTransactionImpl) IncreaseBalanceUser(ctx context.Context, balance *entity.TransactionIncreaseBalanceUser) (float64, error) {
	transaction := &entity.Transaction{
		ID:            balance.ID,
		DestinationId: balance.UserId,
		Amount:        balance.Value,
	}

	err := tr.db.Transaction.Create(ctx, transaction)
	if err != nil {
		log.Println("Error app.Transaction.IncreaseBalanceUser.db.Create: ", err.Error())
		return 0, err
	}

	user, err := tr.db.User.ReadOneById(ctx, balance.UserId)
	if err != nil {
		tr.updateStatusTransaction(ctx, transaction, 2)
		log.Println("Error app.Transaction.IncreaseBalanceUser.db.ReadOneById: ", err.Error())
		return 0, err
	}

	user.Mutex = &sync.Mutex{}
	user.Mutex.Lock()
	defer user.Mutex.Unlock()

	user.Balance += balance.Value

	err = tr.db.Transaction.UpdateBalanceUser(ctx, user.ID, user.Balance)
	if err != nil {
		tr.updateStatusTransaction(ctx, transaction, 2)
		tr.revertDestinationBalanceTransaction(ctx, user, balance.Value)
		log.Println("Error app.Transaction.UpdateBalanceUser.db.UpdateBalanceUser: ", err.Error())
		return 0, err
	}

	newBalance, err := tr.db.Transaction.ReadBalance(ctx, user.ID)
	if err != nil {
		log.Println("Error app.Transaction.IncreaseBalanceUser.db.ReadBalance: ", err.Error())
		return 0, err
	}

	tr.updateStatusTransaction(ctx, transaction, 1)

	return newBalance, nil
}
