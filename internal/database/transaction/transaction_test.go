package transaction

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	"github.com/garoque/backend-code-challenge-snapfi/internal/test"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo/v4"
)

func TestCreate(t *testing.T) {
	query := "INSERT INTO transactions (id, id_source, id_destination, amount, state) VALUES (?, ?, ?, ?, ?)"

	transaction := entity.NewTransaction(dto.CreateTransaction{
		SourceUserId:      "source-user-id",
		DestinationUserId: "destination-user-id",
		Amount:            100,
	})

	cases := map[string]struct {
		InputTransaction *entity.Transaction
		ExpectedErr      error
		PrepareMock      func(mock sqlmock.Sqlmock)
	}{
		"deve retornar sucesso": {
			InputTransaction: transaction,
			ExpectedErr:      nil,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(query).
					WithArgs(transaction.ID, transaction.SourceId, transaction.DestinationId, transaction.Amount, transaction.State.String()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		"deve retornar erro: ao criar transaction": {
			InputTransaction: transaction,
			ExpectedErr:      echo.ErrInternalServerError,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(query).
					WithArgs(transaction.ID, transaction.SourceId, transaction.DestinationId, transaction.Amount, transaction.State.String()).
					WillReturnError(echo.ErrInternalServerError)
				mock.ExpectRollback()
			},
		},
		"deve retornar erro: ao comitar a transaction": {
			InputTransaction: transaction,
			ExpectedErr:      echo.ErrInternalServerError,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(query).
					WithArgs(transaction.ID, transaction.SourceId, transaction.DestinationId, transaction.Amount, transaction.State.String()).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().
					WillReturnError(echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			dbConn, mock := test.GetDB()
			cs.PrepareMock(mock)

			db := NewDatabaseTransaction(dbConn)
			ctx := context.Background()

			err := db.Create(ctx, cs.InputTransaction)
			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestUpdateState(t *testing.T) {
	query := "UPDATE transactions SET state = ? WHERE id = ?"

	cases := map[string]struct {
		InputState  string
		InputId     string
		ExpectedErr error
		PrepareMock func(mock sqlmock.Sqlmock)
	}{
		"deve retornar sucesso": {
			InputState:  entity.BOOKED.String(),
			InputId:     "transaction-id",
			ExpectedErr: nil,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(query).
					WithArgs(entity.BOOKED.String(), "transaction-id").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		"deve retornar erro: ao criar transaction": {
			InputState:  entity.BOOKED.String(),
			InputId:     "transaction-id",
			ExpectedErr: echo.ErrInternalServerError,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(query).
					WithArgs(entity.BOOKED.String(), "transaction-id").
					WillReturnError(echo.ErrInternalServerError)
				mock.ExpectRollback()
			},
		},
		"deve retornar erro: ao comitar a transaction": {
			InputState:  entity.BOOKED.String(),
			InputId:     "transaction-id",
			ExpectedErr: echo.ErrInternalServerError,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(query).
					WithArgs(entity.BOOKED.String(), "transaction-id").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().
					WillReturnError(echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			dbConn, mock := test.GetDB()
			cs.PrepareMock(mock)

			db := NewDatabaseTransaction(dbConn)
			ctx := context.Background()

			err := db.UpdateState(ctx, cs.InputState, cs.InputId)
			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestReadBalance(t *testing.T) {
	query := "SELECT balance FROM users WHERE id = ?"

	var balance float64
	userId := "user-id"

	cases := map[string]struct {
		InputUserId    string
		ExpectedResult float64
		ExpectedErr    error
		PrepareMock    func(mock sqlmock.Sqlmock)
	}{
		"deve retornar sucesso": {
			InputUserId:    userId,
			ExpectedResult: balance,
			ExpectedErr:    nil,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(query).
					WithArgs(userId).
					WillReturnRows(test.NewRows("balance").AddRow(balance))
				mock.ExpectCommit()
			},
		},
		"deve retornar erro: ao criar transaction": {
			InputUserId:    userId,
			ExpectedResult: balance,
			ExpectedErr:    echo.ErrNotFound,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(query).
					WithArgs(userId).
					WillReturnError(echo.ErrNotFound)
				mock.ExpectRollback()
			},
		},
		"deve retornar erro: ao comitar a transaction": {
			InputUserId:    userId,
			ExpectedResult: balance,
			ExpectedErr:    echo.ErrInternalServerError,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(query).
					WithArgs(userId).
					WillReturnRows(test.NewRows("balance").AddRow(balance))
				mock.ExpectCommit().
					WillReturnError(echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			dbConn, mock := test.GetDB()
			cs.PrepareMock(mock)

			db := NewDatabaseTransaction(dbConn)
			ctx := context.Background()

			balance, err := db.ReadBalance(ctx, cs.InputUserId)
			if diff := cmp.Diff(balance, cs.ExpectedResult); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestUpdateBalanceUser(t *testing.T) {
	query := "UPDATE users SET balance = ? WHERE id = ?"

	value := 100.10
	userId := "user-id"

	cases := map[string]struct {
		InputValue  float64
		InputUserId string
		ExpectedErr error
		PrepareMock func(mock sqlmock.Sqlmock)
	}{
		"deve retornar sucesso": {
			InputValue:  value,
			InputUserId: userId,
			ExpectedErr: nil,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(query).
					WithArgs(value, userId).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		"deve retornar erro: ao criar transaction": {
			InputValue:  value,
			InputUserId: userId,
			ExpectedErr: echo.ErrNotFound,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(query).
					WithArgs(value, userId).
					WillReturnError(echo.ErrNotFound)
				mock.ExpectRollback()
			},
		},
		"deve retornar erro: ao comitar a transaction": {
			InputValue:  value,
			InputUserId: userId,
			ExpectedErr: echo.ErrInternalServerError,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(query).
					WithArgs(value, userId).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().
					WillReturnError(echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			dbConn, mock := test.GetDB()
			cs.PrepareMock(mock)

			db := NewDatabaseTransaction(dbConn)
			ctx := context.Background()

			err := db.UpdateBalanceUser(ctx, cs.InputUserId, cs.InputValue)
			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}
