package user

import (
	"context"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	"github.com/garoque/backend-code-challenge-snapfi/internal/test"
	"github.com/garoque/backend-code-challenge-snapfi/pkg/custom_err"
	"github.com/garoque/backend-code-challenge-snapfi/pkg/uuid"
	"github.com/google/go-cmp/cmp"
)

func TestCreate(t *testing.T) {
	query := "INSERT INTO users (id, name, balance) VALUES (?, ?, ?)"

	user := entity.NewUser(dto.CreateUser{Name: "Gabriel"})

	cases := map[string]struct {
		InputUser   entity.User
		ExpectedErr error
		PrepareMock func(mock sqlmock.Sqlmock)
	}{
		"deve retornar sucesso": {
			InputUser:   *user,
			ExpectedErr: nil,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(query).
					WithArgs(user.ID, user.Name, user.Balance).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
		},
		"deve retornar erro: ao criar user": {
			InputUser:   *user,
			ExpectedErr: custom_err.New(http.StatusInternalServerError, custom_err.INTERNAL_ERROR),
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(query).
					WithArgs(user.ID, user.Name, user.Balance).
					WillReturnError(custom_err.New(http.StatusInternalServerError, custom_err.INTERNAL_ERROR))
				mock.ExpectRollback()
			},
		},
		"deve retornar erro: ao comitar a transaction": {
			InputUser:   *user,
			ExpectedErr: custom_err.New(http.StatusInternalServerError, custom_err.INTERNAL_ERROR),
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(query).
					WithArgs(user.ID, user.Name, user.Balance).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit().
					WillReturnError(custom_err.New(http.StatusInternalServerError, custom_err.INTERNAL_ERROR))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			dbConn, mock := test.GetDB()
			cs.PrepareMock(mock)

			db := NewDatabaseUser(dbConn)
			ctx := context.Background()

			err := db.Create(ctx, cs.InputUser)
			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestReadAll(t *testing.T) {
	query := "SELECT id, name, balance, created_at, updated_at FROM users"

	user := entity.NewUser(dto.CreateUser{Name: "Gabriel"})
	users := []entity.User{{
		ID:        user.ID,
		Name:      user.Name,
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt,
		UpdatedAt: nil,
	}}

	cases := map[string]struct {
		ExpectedResult []entity.User
		ExpectedErr    error
		PrepareMock    func(mock sqlmock.Sqlmock)
	}{
		"deve retornar sucesso": {
			ExpectedResult: users,
			ExpectedErr:    nil,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WillReturnRows(
						test.NewRows("id", "name", "balance", "created_at", "updated_at").
							AddRow(user.ID, user.Name, user.Balance, user.CreatedAt, nil),
					)
			},
		},
		"deve retornar erro": {
			ExpectedResult: nil,
			ExpectedErr:    custom_err.New(http.StatusInternalServerError, custom_err.INTERNAL_ERROR),
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WillReturnError(custom_err.New(http.StatusInternalServerError, custom_err.INTERNAL_ERROR))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			dbConn, mock := test.GetDB()
			cs.PrepareMock(mock)

			db := NewDatabaseUser(dbConn)
			ctx := context.Background()

			users, err := db.ReadAll(ctx)
			if diff := cmp.Diff(users, cs.ExpectedResult); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestReadOneById(t *testing.T) {
	query := "SELECT id, name, balance, created_at, updated_at FROM users WHERE id = ?"

	user := &entity.User{
		ID:   uuid.NewId(),
		Name: "Gabriel",
	}

	cases := map[string]struct {
		InputUserId    string
		ExpectedResult *entity.User
		ExpectedErr    error
		PrepareMock    func(mock sqlmock.Sqlmock)
	}{
		"deve retornar sucesso": {
			InputUserId:    user.ID,
			ExpectedResult: user,
			ExpectedErr:    nil,
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs(user.ID).
					WillReturnRows(
						test.NewRows("id", "name", "balance", "created_at", "updated_at").
							AddRow(user.ID, user.Name, user.Balance, user.CreatedAt, nil),
					)
			},
		},
		"deve retornar erro": {
			InputUserId:    user.ID,
			ExpectedResult: nil,
			ExpectedErr:    custom_err.New(http.StatusInternalServerError, custom_err.INTERNAL_ERROR),
			PrepareMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).
					WithArgs(user.ID).
					WillReturnError(
						custom_err.New(http.StatusInternalServerError, custom_err.INTERNAL_ERROR),
					)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			dbConn, mock := test.GetDB()
			cs.PrepareMock(mock)

			db := NewDatabaseUser(dbConn)
			ctx := context.Background()

			users, err := db.ReadOneById(ctx, cs.InputUserId)
			if diff := cmp.Diff(users, cs.ExpectedResult); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}
