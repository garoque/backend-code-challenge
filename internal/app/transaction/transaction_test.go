package transaction

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/garoque/backend-code-challenge-snapfi/internal/database"
	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	"github.com/garoque/backend-code-challenge-snapfi/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo/v4"
)

func TestCreate(t *testing.T) {
	sourceUserId := "source-user-id"
	sourceUserId2 := "source-user-id-2"
	destinationUserId := "destination-user-id"
	transaction := entity.NewTransaction(dto.CreateTransaction{
		SourceUserId:      sourceUserId,
		DestinationUserId: destinationUserId,
		Amount:            100.10,
	})

	bookedTransaction := entity.Transaction{
		ID:            transaction.ID,
		SourceId:      transaction.SourceId,
		DestinationId: transaction.DestinationId,
		Amount:        transaction.Amount,
		State:         entity.BOOKED,
		StateString:   entity.BOOKED.String(),
	}

	failedTransaction := entity.NewTransaction(dto.CreateTransaction{
		SourceUserId:      sourceUserId,
		DestinationUserId: destinationUserId,
		Amount:            250.10,
	})

	sourceUser := entity.User{
		ID:        sourceUserId,
		Name:      "Gabriel",
		Balance:   200.0,
		Mutex:     &sync.Mutex{},
		CreatedAt: time.Now(),
	}

	sourceUser2 := entity.User{
		ID:        sourceUserId2,
		Name:      "Gabriel",
		Balance:   2000.0,
		Mutex:     &sync.Mutex{},
		CreatedAt: time.Now(),
	}

	destinationUser := entity.User{
		ID:        destinationUserId,
		Name:      "João",
		Balance:   0,
		Mutex:     &sync.Mutex{},
		CreatedAt: time.Now(),
	}

	destinationUser2 := entity.User{
		ID:        destinationUserId,
		Name:      "João",
		Balance:   0,
		Mutex:     &sync.Mutex{},
		CreatedAt: time.Now(),
	}

	sourceUserBalanceUpdated := sourceUser.Balance - transaction.Amount
	destinationUserBalanceUpdated := destinationUser.Balance + transaction.Amount

	sourceUserFailedTrBalanceUpdated := sourceUser2.Balance - failedTransaction.Amount
	destinationUserFailedTrBalanceUpdated := destinationUser2.Balance + failedTransaction.Amount

	cases := map[string]struct {
		InputTransaction *entity.Transaction
		ExpectedResult   *entity.Transaction
		ExpectedErr      error
		PrepareMock      func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface)
	}{
		"deve retornar sucesso": {
			InputTransaction: transaction,
			ExpectedResult:   &bookedTransaction,
			ExpectedErr:      nil,
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().Create(gomock.Any(), transaction).Times(1).Return(nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), transaction.SourceId).Times(1).Return(&sourceUser, nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), transaction.DestinationId).Times(1).Return(&destinationUser, nil)

				mockTransactionDb.EXPECT().UpdateBalanceUser(gomock.Any(), sourceUser.ID, sourceUserBalanceUpdated).
					Times(1).Return(nil)

				mockTransactionDb.EXPECT().UpdateBalanceUser(gomock.Any(), destinationUser.ID, destinationUserBalanceUpdated).
					Times(1).Return(nil)

				mockTransactionDb.EXPECT().UpdateState(gomock.Any(), bookedTransaction.State, transaction.ID).Times(1).Return(nil)
			},
		},
		"deve retornar erro: ao registrar transaction": {
			InputTransaction: transaction,
			ExpectedResult:   nil,
			ExpectedErr:      echo.ErrInternalServerError,
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().Create(gomock.Any(), transaction).Times(1).Return(echo.ErrInternalServerError)
			},
		},
		"deve retornar erro: ao ler source user": {
			InputTransaction: transaction,
			ExpectedResult:   transaction,
			ExpectedErr:      echo.ErrInternalServerError,
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().Create(gomock.Any(), transaction).Times(1).Return(nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), transaction.SourceId).Times(1).Return(nil, echo.ErrInternalServerError)

				mockTransactionDb.EXPECT().UpdateState(gomock.Any(), entity.FAILED, transaction.ID).Times(1).Return(nil)
			},
		},
		"deve retornar erro: ao ler destination user": {
			InputTransaction: transaction,
			ExpectedResult:   transaction,
			ExpectedErr:      echo.ErrInternalServerError,
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().Create(gomock.Any(), transaction).Times(1).Return(nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), transaction.SourceId).Times(1).Return(&sourceUser, nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), transaction.DestinationId).Times(1).Return(nil, echo.ErrInternalServerError)

				mockTransactionDb.EXPECT().UpdateState(gomock.Any(), entity.FAILED, transaction.ID).Times(1).Return(nil)
			},
		},
		"deve retornar erro: 'Insufficient balance'": {
			InputTransaction: failedTransaction,
			ExpectedResult:   failedTransaction,
			ExpectedErr:      echo.NewHTTPError(echo.ErrBadRequest.Code, "Insufficient balance"),
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().Create(gomock.Any(), failedTransaction).Times(1).Return(nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), failedTransaction.SourceId).Times(1).Return(&sourceUser, nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), failedTransaction.DestinationId).Times(1).Return(&destinationUser, nil)

				mockTransactionDb.EXPECT().UpdateState(gomock.Any(), entity.FAILED, failedTransaction.ID).Times(1).Return(nil)
			},
		},
		"deve retornar erro: ao atualizar saldo source user": {
			InputTransaction: failedTransaction,
			ExpectedResult:   failedTransaction,
			ExpectedErr:      echo.ErrInternalServerError,
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().Create(gomock.Any(), failedTransaction).Times(1).Return(nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), failedTransaction.SourceId).Times(1).Return(&sourceUser2, nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), failedTransaction.DestinationId).Times(1).Return(&destinationUser, nil)

				mockTransactionDb.EXPECT().UpdateBalanceUser(gomock.Any(), sourceUser2.ID, sourceUserFailedTrBalanceUpdated).
					Times(1).Return(echo.ErrInternalServerError)

				mockTransactionDb.EXPECT().UpdateBalanceUser(gomock.Any(), sourceUser2.ID, sourceUser2.Balance).
					Times(1).Return(nil)

				mockTransactionDb.EXPECT().UpdateState(gomock.Any(), entity.FAILED, failedTransaction.ID).Times(1).Return(nil)
			},
		},
		"deve retornar erro: ao atualizar saldo destination user": {
			InputTransaction: failedTransaction,
			ExpectedResult:   failedTransaction,
			ExpectedErr:      echo.ErrInternalServerError,
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().Create(gomock.Any(), failedTransaction).Times(1).Return(nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), failedTransaction.SourceId).Times(1).Return(&sourceUser2, nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), failedTransaction.DestinationId).Times(1).Return(&destinationUser2, nil)

				mockTransactionDb.EXPECT().UpdateBalanceUser(gomock.Any(), sourceUser2.ID, sourceUserFailedTrBalanceUpdated).
					Times(1).Return(nil)

				mockTransactionDb.EXPECT().UpdateBalanceUser(gomock.Any(), destinationUser2.ID, destinationUserFailedTrBalanceUpdated).
					Times(1).Return(echo.ErrInternalServerError)

				mockTransactionDb.EXPECT().UpdateBalanceUser(gomock.Any(), destinationUser2.ID, destinationUser2.Balance).
					Times(1).Return(nil)

				mockTransactionDb.EXPECT().UpdateBalanceUser(gomock.Any(), sourceUser2.ID, sourceUser2.Balance).
					Times(1).Return(nil)

				mockTransactionDb.EXPECT().UpdateState(gomock.Any(), entity.FAILED, failedTransaction.ID).Times(1).Return(nil)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockTransactionDb := mocks.NewMockDabataseTransactionInterface(ctrl)
			mockUserDb := mocks.NewMockDabataseUserInterface(ctrl)
			cs.PrepareMock(mockTransactionDb, mockUserDb)

			app := NewAppTransaction(&database.Container{Transaction: mockTransactionDb, User: mockUserDb})

			transaction, err := app.Create(ctx, cs.InputTransaction)
			if diff := cmp.Diff(transaction, cs.ExpectedResult); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestIncreaseBalanceUser(t *testing.T) {
	destinationUserId := "destination-user-id"
	balance := entity.NewIncreaseBalanceUser(dto.IncreaseBalanceUser{
		UserId: destinationUserId,
		Value:  110.0,
	})
	transaction := &entity.Transaction{
		ID:            balance.ID,
		DestinationId: balance.UserId,
		Amount:        balance.Value,
	}

	bookedTransaction := entity.Transaction{
		ID:            transaction.ID,
		SourceId:      transaction.SourceId,
		DestinationId: transaction.DestinationId,
		Amount:        transaction.Amount,
		State:         entity.BOOKED,
		StateString:   entity.BOOKED.String(),
	}

	destinationUser := &entity.User{
		ID:        destinationUserId,
		Name:      "Gabriel",
		Balance:   200.0,
		Mutex:     &sync.Mutex{},
		CreatedAt: time.Now(),
	}

	destinationUser2 := &entity.User{
		ID:        destinationUserId,
		Name:      "Gabriel",
		Balance:   200.0,
		Mutex:     &sync.Mutex{},
		CreatedAt: time.Now(),
	}

	destinationUser3 := &entity.User{
		ID:        destinationUserId,
		Name:      "Gabriel",
		Balance:   200.0,
		Mutex:     &sync.Mutex{},
		CreatedAt: time.Now(),
	}

	balanceUserUpdated := destinationUser.Balance + transaction.Amount
	balanceUser := destinationUser2.Balance + transaction.Amount

	cases := map[string]struct {
		InputBalance   *entity.TransactionIncreaseBalanceUser
		ExpectedResult float64
		ExpectedErr    error
		PrepareMock    func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface)
	}{
		"deve retornar sucesso": {
			InputBalance:   balance,
			ExpectedResult: balanceUserUpdated,
			ExpectedErr:    nil,
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().Create(gomock.Any(), transaction).Times(1).Return(nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), balance.UserId).Times(1).Return(destinationUser, nil)
				mockTransactionDb.EXPECT().UpdateBalanceUser(gomock.Any(), destinationUser.ID, balanceUserUpdated).
					Times(1).Return(nil)

				mockTransactionDb.EXPECT().UpdateState(gomock.Any(), bookedTransaction.State, transaction.ID).Times(1).Return(nil)
				mockTransactionDb.EXPECT().ReadBalance(gomock.Any(), destinationUser.ID).Times(1).Return(balanceUserUpdated, nil)
			},
		},
		"deve retornar erro: ao registrar transaction": {
			InputBalance:   balance,
			ExpectedResult: 0,
			ExpectedErr:    echo.ErrInternalServerError,
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().Create(gomock.Any(), transaction).Times(1).Return(echo.ErrInternalServerError)
			},
		},
		"deve retornar erro: ao ler destination user": {
			InputBalance:   balance,
			ExpectedResult: 0,
			ExpectedErr:    echo.ErrInternalServerError,
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().Create(gomock.Any(), transaction).Times(1).Return(nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), balance.UserId).Times(1).Return(nil, echo.ErrInternalServerError)

				mockTransactionDb.EXPECT().UpdateState(gomock.Any(), entity.FAILED, transaction.ID).Times(1).Return(nil)
			},
		},
		"deve retornar erro: ao atualizar saldo destination user": {
			InputBalance:   balance,
			ExpectedResult: 0,
			ExpectedErr:    echo.ErrInternalServerError,
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().Create(gomock.Any(), transaction).Times(1).Return(nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), balance.UserId).Times(1).Return(destinationUser2, nil)

				mockTransactionDb.EXPECT().UpdateBalanceUser(gomock.Any(), destinationUser2.ID, balanceUser).
					Times(1).Return(echo.ErrInternalServerError)

				mockTransactionDb.EXPECT().UpdateBalanceUser(gomock.Any(), destinationUser2.ID, destinationUser2.Balance).
					Times(1).Return(nil)

				mockTransactionDb.EXPECT().UpdateState(gomock.Any(), entity.FAILED, transaction.ID).Times(1).Return(nil)
			},
		},
		"deve retornar erro: ao ler o saldo": {
			InputBalance:   balance,
			ExpectedResult: 0,
			ExpectedErr:    echo.ErrInternalServerError,
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().Create(gomock.Any(), transaction).Times(1).Return(nil)
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), balance.UserId).Times(1).Return(destinationUser3, nil)
				mockTransactionDb.EXPECT().UpdateBalanceUser(gomock.Any(), destinationUser3.ID, balanceUserUpdated).
					Times(1).Return(nil)

				mockTransactionDb.EXPECT().UpdateState(gomock.Any(), bookedTransaction.State, transaction.ID).Times(1).Return(nil)
				mockTransactionDb.EXPECT().ReadBalance(gomock.Any(), destinationUser3.ID).Times(1).Return(0.0, echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockTransactionDb := mocks.NewMockDabataseTransactionInterface(ctrl)
			mockUserDb := mocks.NewMockDabataseUserInterface(ctrl)
			cs.PrepareMock(mockTransactionDb, mockUserDb)

			app := NewAppTransaction(&database.Container{Transaction: mockTransactionDb, User: mockUserDb})

			balance, err := app.IncreaseBalanceUser(ctx, cs.InputBalance)
			if diff := cmp.Diff(balance, cs.ExpectedResult); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestReadAll(t *testing.T) {
	transaction := entity.NewTransaction(dto.CreateTransaction{
		SourceUserId:      "source-user-id",
		DestinationUserId: "destination-user-id",
		Amount:            100.10,
	})
	transactions := []entity.Transaction{{
		ID:            transaction.ID,
		SourceId:      transaction.SourceId,
		DestinationId: transaction.DestinationId,
		Amount:        transaction.Amount,
	}}

	cases := map[string]struct {
		ExpectedResult []entity.Transaction
		ExpectedErr    error
		PrepareMock    func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface)
	}{
		"deve retornar sucesso": {
			ExpectedResult: transactions,
			ExpectedErr:    nil,
			PrepareMock: func(mockTransactionDb *mocks.MockDabataseTransactionInterface, mockUserDb *mocks.MockDabataseUserInterface) {
				mockTransactionDb.EXPECT().ReadAll(gomock.Any()).Times(1).Return(transactions, nil)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockTransactionDb := mocks.NewMockDabataseTransactionInterface(ctrl)
			mockUserDb := mocks.NewMockDabataseUserInterface(ctrl)
			cs.PrepareMock(mockTransactionDb, mockUserDb)

			app := NewAppTransaction(&database.Container{Transaction: mockTransactionDb, User: mockUserDb})

			transactions, err := app.ReadAll(ctx)
			if diff := cmp.Diff(transactions, cs.ExpectedResult); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}
