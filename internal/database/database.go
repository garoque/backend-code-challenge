package database

import (
	"github.com/garoque/backend-code-challenge-snapfi/internal/database/transaction"
	"github.com/garoque/backend-code-challenge-snapfi/internal/database/user"
	"github.com/jmoiron/sqlx"
)

type Container struct {
	User        user.DabataseUserInterface
	Transaction transaction.DabataseTransactionInterface
}

func New(dbConn *sqlx.DB) *Container {
	return &Container{
		User:        user.NewDatabaseUser(dbConn),
		Transaction: transaction.NewDatabaseTransaction(dbConn),
	}
}
