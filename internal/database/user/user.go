package user

import (
	"context"
	"database/sql"
	"log"

	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type DabataseUserInterface interface {
	Create(ctx context.Context, user entity.User) error
	ReadAll(ctx context.Context) ([]entity.User, error)
	ReadOneById(ctx context.Context, userId string) (*entity.User, error)
}

type dbImpl struct {
	dbConn *sqlx.DB
}

func NewDatabaseUser(dbConn *sqlx.DB) DabataseUserInterface {
	return &dbImpl{dbConn}
}

func (u *dbImpl) Create(ctx context.Context, user entity.User) error {
	tx, _ := u.dbConn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelDefault})
	query := "INSERT INTO users (id, name, balance) VALUES (?, ?, ?)"

	_, err := tx.ExecContext(ctx, query, user.ID, user.Name, user.Balance)
	if err != nil {
		tx.Rollback()
		log.Println("Error create user: ", err.Error())
		return echo.ErrInternalServerError
	}

	err = tx.Commit()
	if err != nil {
		log.Println("Error create user tx.Commit: ", err.Error())
		return echo.ErrInternalServerError
	}

	return nil
}

func (u *dbImpl) ReadAll(ctx context.Context) ([]entity.User, error) {
	users := make([]entity.User, 0)
	query := "SELECT id, name, balance, created_at, updated_at FROM users"

	err := u.dbConn.SelectContext(ctx, &users, query)
	if err != nil {
		log.Println("Error ReadAll user: ", err.Error())
		return nil, echo.ErrInternalServerError
	}

	return users, nil
}

func (u *dbImpl) ReadOneById(ctx context.Context, userId string) (*entity.User, error) {
	user := new(entity.User)
	query := "SELECT id, name, balance, created_at, updated_at FROM users WHERE id = ?"

	err := u.dbConn.GetContext(ctx, user, query, userId)
	if err != nil {
		log.Println("Error ReadOneById user: ", err.Error())
		return nil, echo.ErrNotFound
	}

	return user, nil
}
