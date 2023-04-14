package transaction

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/garoque/backend-code-challenge-snapfi/internal/app"
	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	"github.com/garoque/backend-code-challenge-snapfi/internal/mocks"
	"github.com/garoque/backend-code-challenge-snapfi/pkg/uuid"
	"github.com/garoque/backend-code-challenge-snapfi/pkg/validator"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	request := dto.CreateTransaction{
		SourceUserId:      "1234",
		DestinationUserId: "5678",
		Amount:            100.0,
	}

	transaction := &entity.Transaction{
		ID:            uuid.NewId(),
		SourceId:      request.SourceUserId,
		DestinationId: request.DestinationUserId,
		Amount:        request.Amount,
	}

	cases := map[string]struct {
		InputTransaction dto.CreateTransaction
		ExpectedResult   *entity.Transaction
		ExpectedErr      error
		PrepareMock      func(mockTransactionApp *mocks.MockAppTransactionInterface)
	}{
		"deve retornar sucesso": {
			InputTransaction: request,
			ExpectedResult:   transaction,
			ExpectedErr:      nil,
			PrepareMock: func(mockTransactionApp *mocks.MockAppTransactionInterface) {
				mockTransactionApp.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(transaction, nil)
			},
		},
		"deve retornar erro": {
			InputTransaction: request,
			ExpectedResult:   transaction,
			ExpectedErr:      echo.ErrInternalServerError,
			PrepareMock: func(mockTransactionApp *mocks.MockAppTransactionInterface) {
				mockTransactionApp.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(transaction, echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockTransactionApp := mocks.NewMockAppTransactionInterface(ctrl)
			cs.PrepareMock(mockTransactionApp)

			api := handler{
				app: &app.Container{Transaction: mockTransactionApp},
			}

			e := echo.New()
			e.Validator = validator.NewValidator()

			endpoint := "/v1/transaction"

			requestBytes, _ := json.Marshal(cs.InputTransaction)
			req := httptest.NewRequest(http.MethodPost, endpoint, bytes.NewReader(requestBytes)).WithContext(ctx)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath(endpoint)

			err := api.create(c)
			assert.Equal(t, cs.ExpectedErr, err)

			if err == nil {
				expectedResultJSON, err := json.Marshal(dto.Response{Data: transaction})
				assert.NoError(t, err)

				var expectedResult dto.Response
				err = json.Unmarshal(expectedResultJSON, &expectedResult)
				assert.NoError(t, err)

				var currentResult dto.Response
				json.NewDecoder(rec.Body).Decode(&currentResult)

				assert.Equal(t, expectedResult, currentResult)
			}
		})
	}
}

func TestIncreaseBalance(t *testing.T) {
	transaction := dto.IncreaseBalanceUser{
		UserId: "user-id",
		Value:  100.10,
	}

	balance := transaction.Value

	cases := map[string]struct {
		InputTransaction dto.IncreaseBalanceUser
		ExpectedResult   float64
		ExpectedErr      error
		PrepareMock      func(mockTransactionApp *mocks.MockAppTransactionInterface)
	}{
		"deve retornar sucesso": {
			InputTransaction: transaction,
			ExpectedResult:   balance,
			ExpectedErr:      nil,
			PrepareMock: func(mockTransactionApp *mocks.MockAppTransactionInterface) {
				mockTransactionApp.EXPECT().IncreaseBalanceUser(gomock.Any(), gomock.Any()).Times(1).Return(balance, nil)
			},
		},
		"deve retornar erro": {
			InputTransaction: transaction,
			ExpectedResult:   0,
			ExpectedErr:      echo.ErrInternalServerError,
			PrepareMock: func(mockTransactionApp *mocks.MockAppTransactionInterface) {
				mockTransactionApp.EXPECT().IncreaseBalanceUser(gomock.Any(), gomock.Any()).Times(1).Return(float64(0), echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockTransactionApp := mocks.NewMockAppTransactionInterface(ctrl)
			cs.PrepareMock(mockTransactionApp)

			api := handler{
				app: &app.Container{Transaction: mockTransactionApp},
			}

			e := echo.New()
			e.Validator = validator.NewValidator()

			endpoint := "/v1/transaction/increase-balance"

			requestBytes, _ := json.Marshal(cs.InputTransaction)
			req := httptest.NewRequest(http.MethodPut, endpoint, bytes.NewReader(requestBytes)).WithContext(ctx)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath(endpoint)

			err := api.increaseBalance(c)
			assert.Equal(t, cs.ExpectedErr, err)

			if err == nil {
				expectedResultJSON, err := json.Marshal(dto.Response{Data: balance})
				assert.NoError(t, err)

				var expectedResult dto.Response
				err = json.Unmarshal(expectedResultJSON, &expectedResult)
				assert.NoError(t, err)

				var currentResult dto.Response
				json.NewDecoder(rec.Body).Decode(&currentResult)

				assert.Equal(t, expectedResult, currentResult)
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
		PrepareMock    func(mockTransactionApp *mocks.MockAppTransactionInterface)
	}{
		"deve retornar sucesso": {
			ExpectedResult: transactions,
			ExpectedErr:    nil,
			PrepareMock: func(mockTransactionApp *mocks.MockAppTransactionInterface) {
				mockTransactionApp.EXPECT().ReadAll(gomock.Any()).Times(1).Return(transactions, nil)
			},
		},
		"deve retornar erro": {
			ExpectedResult: nil,
			ExpectedErr:    echo.ErrInternalServerError,
			PrepareMock: func(mockTransactionApp *mocks.MockAppTransactionInterface) {
				mockTransactionApp.EXPECT().ReadAll(gomock.Any()).Times(1).Return(nil, echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockTransactionApp := mocks.NewMockAppTransactionInterface(ctrl)
			cs.PrepareMock(mockTransactionApp)

			api := handler{
				app: &app.Container{Transaction: mockTransactionApp},
			}

			e := echo.New()

			endpoint := "/v1/transaction"
			req := httptest.NewRequest(http.MethodGet, endpoint, nil).WithContext(ctx)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath(endpoint)

			err := api.readAll(c)
			assert.Equal(t, cs.ExpectedErr, err)

			if err == nil {
				expectedResultJSON, err := json.Marshal(dto.Response{Data: transactions})
				assert.NoError(t, err)

				var expectedResult dto.Response
				err = json.Unmarshal(expectedResultJSON, &expectedResult)
				assert.NoError(t, err)

				var currentResult dto.Response
				json.NewDecoder(rec.Body).Decode(&currentResult)

				assert.Equal(t, expectedResult, currentResult)
			}
		})
	}
}
