package app

import (
	"github.com/garoque/backend-code-challenge-snapfi/internal/app/transaction"
	"github.com/garoque/backend-code-challenge-snapfi/internal/app/user"
	"github.com/garoque/backend-code-challenge-snapfi/internal/database"
)

type Container struct {
	User        user.AppUserInterface
	Transaction transaction.AppTransactionInterface
}

func New(db *database.Container) *Container {
	return &Container{
		User:        user.NewAppUser(db),
		Transaction: transaction.NewAppTransaction(db),
	}
}
