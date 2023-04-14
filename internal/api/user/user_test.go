package user

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/garoque/backend-code-challenge-snapfi/internal/app"
	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	"github.com/garoque/backend-code-challenge-snapfi/internal/mocks"
	"github.com/garoque/backend-code-challenge-snapfi/pkg/validator"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	cases := map[string]struct {
		InputUserDto dto.CreateUser
		ExpectedErr  error
		PrepareMock  func(mockUserApp *mocks.MockAppUserInterface)
	}{
		"deve retornar sucesso": {
			InputUserDto: dto.CreateUser{Name: "Gabriel"},
			ExpectedErr:  nil,
			PrepareMock: func(mockUserApp *mocks.MockAppUserInterface) {
				mockUserApp.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(nil)
			},
		},
		"deve retornar erro": {
			InputUserDto: dto.CreateUser{Name: "Gabriel"},
			ExpectedErr:  echo.ErrInternalServerError,
			PrepareMock: func(mockUserApp *mocks.MockAppUserInterface) {
				mockUserApp.EXPECT().Create(gomock.Any(), gomock.Any()).Times(1).Return(echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockUserApp := mocks.NewMockAppUserInterface(ctrl)
			cs.PrepareMock(mockUserApp)

			api := handler{
				app: &app.Container{User: mockUserApp},
			}

			e := echo.New()
			e.Validator = validator.NewValidator()

			endpoint := "/v1/user"

			req := httptest.NewRequest(http.MethodPost, endpoint, strings.NewReader(`{"name": "Gabriel"}`)).WithContext(ctx)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath(endpoint)

			err := api.create(c)
			assert.Equal(t, cs.ExpectedErr, err)
		})
	}
}

func TestReadOne(t *testing.T) {
	user := entity.NewUser(dto.CreateUser{Name: "Gabriel"})
	userId := user.ID

	cases := map[string]struct {
		InputUserId    string
		ExpectedResult *entity.User
		ExpectedErr    error
		PrepareMock    func(mockUserApp *mocks.MockAppUserInterface)
	}{
		"deve retornar sucesso": {
			InputUserId:    userId,
			ExpectedResult: user,
			ExpectedErr:    nil,
			PrepareMock: func(mockUserApp *mocks.MockAppUserInterface) {
				mockUserApp.EXPECT().ReadOneById(gomock.Any(), userId).Times(1).Return(user, nil)
			},
		},
		"deve retornar erro": {
			InputUserId:    userId,
			ExpectedResult: nil,
			ExpectedErr:    echo.ErrInternalServerError,
			PrepareMock: func(mockUserApp *mocks.MockAppUserInterface) {
				mockUserApp.EXPECT().ReadOneById(gomock.Any(), userId).Times(1).Return(nil, echo.ErrInternalServerError)
			},
		},
		"deve retornar erro: 'The provided ID is empty'": {
			InputUserId:    "",
			ExpectedResult: nil,
			ExpectedErr:    echo.NewHTTPError(echo.ErrBadRequest.Code, "The provided ID is empty"),
			PrepareMock: func(mockUserApp *mocks.MockAppUserInterface) {
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockUserApp := mocks.NewMockAppUserInterface(ctrl)
			cs.PrepareMock(mockUserApp)

			api := handler{
				app: &app.Container{User: mockUserApp},
			}

			e := echo.New()

			endpoint := "/v1/user/:id"
			req := httptest.NewRequest(http.MethodGet, endpoint, nil).WithContext(ctx)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath(endpoint)
			c.SetParamNames("id")
			c.SetParamValues(cs.InputUserId)

			err := api.readOne(c)
			assert.Equal(t, cs.ExpectedErr, err)

			if err == nil {
				expectedResultJSON, err := json.Marshal(dto.Response{Data: user})
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
	user := entity.NewUser(dto.CreateUser{Name: "Gabriel"})
	users := []entity.User{*user}

	cases := map[string]struct {
		ExpectedResult []entity.User
		ExpectedErr    error
		PrepareMock    func(mockUserApp *mocks.MockAppUserInterface)
	}{
		"deve retornar sucesso": {
			ExpectedResult: users,
			ExpectedErr:    nil,
			PrepareMock: func(mockUserApp *mocks.MockAppUserInterface) {
				mockUserApp.EXPECT().ReadAll(gomock.Any()).Times(1).Return(users, nil)
			},
		},
		"deve retornar erro": {
			ExpectedResult: nil,
			ExpectedErr:    echo.ErrInternalServerError,
			PrepareMock: func(mockUserApp *mocks.MockAppUserInterface) {
				mockUserApp.EXPECT().ReadAll(gomock.Any()).Times(1).Return(nil, echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockUserApp := mocks.NewMockAppUserInterface(ctrl)
			cs.PrepareMock(mockUserApp)

			api := handler{
				app: &app.Container{User: mockUserApp},
			}

			e := echo.New()

			endpoint := "/v1/user"
			req := httptest.NewRequest(http.MethodGet, endpoint, nil).WithContext(ctx)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath(endpoint)

			err := api.readAll(c)
			assert.Equal(t, cs.ExpectedErr, err)

			if err == nil {
				expectedResultJSON, err := json.Marshal(dto.Response{Data: users})
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
