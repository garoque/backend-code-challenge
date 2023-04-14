package transaction

import (
	"context"
	"database/sql"
	"log"

	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type DabataseTransactionInterface interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	UpdateState(ctx context.Context, state entity.StatesTransaction, id string) error
	ReadBalance(ctx context.Context, userId string) (float64, error)
	UpdateBalanceUser(ctx context.Context, userId string, value float64) error
	ReadAll(ctx context.Context) ([]entity.Transaction, error)
}

type dbImpl struct {
	dbConn *sqlx.DB
}

func NewDatabaseTransaction(dbConn *sqlx.DB) DabataseTransactionInterface {
	return &dbImpl{dbConn}
}

func (tr *dbImpl) Create(ctx context.Context, transaction *entity.Transaction) error {
	tx, _ := tr.dbConn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	query := "INSERT INTO transactions (id, id_source, id_destination, amount, state) VALUES (?, ?, ?, ?, ?)"

	_, err := tx.ExecContext(ctx, query,
		transaction.ID,
		transaction.SourceId,
		transaction.DestinationId,
		transaction.Amount,
		transaction.State,
	)
	if err != nil {
		tx.Rollback()
		log.Println("Error create transaction: ", err.Error())
		return echo.ErrInternalServerError
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error create transaction tx.Commit: ", err.Error())
		return echo.ErrInternalServerError
	}

	return nil
}

func (tr *dbImpl) UpdateState(ctx context.Context, state entity.StatesTransaction, id string) error {
	tx, _ := tr.dbConn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	query := "UPDATE transactions SET state = ? WHERE id = ?"

	_, err := tx.ExecContext(ctx, query, state, id)
	if err != nil {
		tx.Rollback()
		log.Println("Error create transaction: ", err.Error())
		return echo.ErrInternalServerError
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error create transaction tx.Commit: ", err.Error())
		return echo.ErrInternalServerError
	}

	return nil
}

func (tr *dbImpl) ReadBalance(ctx context.Context, userId string) (float64, error) {
	tx, _ := tr.dbConn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	query := "SELECT balance FROM users WHERE id = ?"

	var balance float64

	err := tx.QueryRowContext(ctx, query, userId).Scan(&balance)
	if err != nil {
		tx.Rollback()
		log.Println("Error read balance: ", err.Error())
		return balance, echo.ErrNotFound
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error read balance tx.Commit: ", err.Error())
		return balance, echo.ErrInternalServerError
	}

	return balance, nil
}

func (tr *dbImpl) UpdateBalanceUser(ctx context.Context, userId string, value float64) error {
	tx, _ := tr.dbConn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	query := "UPDATE users SET balance = ? WHERE id = ?"

	_, err := tx.ExecContext(ctx, query, value, userId)
	if err != nil {
		tx.Rollback()
		log.Println("Error increase balance: ", err.Error())
		return echo.ErrNotFound
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error increase balance tx.Commit: ", err.Error())
		return echo.ErrInternalServerError
	}

	return nil
}

func (tr *dbImpl) ReadAll(ctx context.Context) ([]entity.Transaction, error) {
	transactions := make([]entity.Transaction, 0)
	query := "SELECT id, id_source, id_destination, amount, state, created_at FROM transactions ORDER BY created_at DESC"

	err := tr.dbConn.SelectContext(ctx, &transactions, query)
	if err != nil {
		log.Println("Error ReadAll transactions: ", err.Error())
		return nil, echo.ErrInternalServerError
	}

	return transactions, nil
}
