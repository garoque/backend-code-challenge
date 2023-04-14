package user

import (
	"context"
	"testing"

	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/garoque/backend-code-challenge-snapfi/internal/database"
	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
	"github.com/garoque/backend-code-challenge-snapfi/internal/mocks"
	"github.com/garoque/backend-code-challenge-snapfi/pkg/uuid"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/labstack/echo/v4"
)

func TestCreate(t *testing.T) {
	user := *entity.NewUser(dto.CreateUser{Name: "Gabriel"})

	cases := map[string]struct {
		InputUser   entity.User
		ExpectedErr error
		PrepareMock func(mockUserDb *mocks.MockDabataseUserInterface)
	}{
		"deve retornar sucesso": {
			InputUser:   user,
			ExpectedErr: nil,
			PrepareMock: func(mockUserDb *mocks.MockDabataseUserInterface) {
				mockUserDb.EXPECT().Create(gomock.Any(), user).Times(1).Return(nil)
			},
		},
		"deve retornar erro": {
			InputUser:   user,
			ExpectedErr: echo.ErrInternalServerError,
			PrepareMock: func(mockUserDb *mocks.MockDabataseUserInterface) {
				mockUserDb.EXPECT().Create(gomock.Any(), user).Times(1).Return(echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockUserDb := mocks.NewMockDabataseUserInterface(ctrl)
			cs.PrepareMock(mockUserDb)

			app := NewAppUser(&database.Container{User: mockUserDb})

			err := app.Create(ctx, cs.InputUser)
			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestReadOneById(t *testing.T) {
	user := &entity.User{
		ID:   uuid.NewId(),
		Name: "Gabriel",
	}
	userId := user.ID

	cases := map[string]struct {
		InputUserId    string
		ExpectedResult *entity.User
		ExpectedErr    error
		PrepareMock    func(mockUserDb *mocks.MockDabataseUserInterface)
	}{
		"deve retornar sucesso": {
			InputUserId:    userId,
			ExpectedResult: user,
			ExpectedErr:    nil,
			PrepareMock: func(mockUserDb *mocks.MockDabataseUserInterface) {
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), userId).Times(1).Return(user, nil)
			},
		},
		"deve retornar erro": {
			InputUserId:    userId,
			ExpectedResult: nil,
			ExpectedErr:    echo.ErrInternalServerError,
			PrepareMock: func(mockUserDb *mocks.MockDabataseUserInterface) {
				mockUserDb.EXPECT().ReadOneById(gomock.Any(), userId).Times(1).Return(nil, echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockUserDb := mocks.NewMockDabataseUserInterface(ctrl)
			cs.PrepareMock(mockUserDb)

			app := NewAppUser(&database.Container{User: mockUserDb})

			user, err := app.ReadOneById(ctx, cs.InputUserId)
			if diff := cmp.Diff(user, cs.ExpectedResult); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestReadAll(t *testing.T) {
	users := []entity.User{{
		ID:   uuid.NewId(),
		Name: "Gabriel",
	}}

	cases := map[string]struct {
		ExpectedResult []entity.User
		ExpectedErr    error
		PrepareMock    func(mockUserDb *mocks.MockDabataseUserInterface)
	}{
		"deve retornar sucesso": {
			ExpectedResult: users,
			ExpectedErr:    nil,
			PrepareMock: func(mockUserDb *mocks.MockDabataseUserInterface) {
				mockUserDb.EXPECT().ReadAll(gomock.Any()).Times(1).Return(users, nil)
			},
		},
		"deve retornar erro": {
			ExpectedResult: nil,
			ExpectedErr:    echo.ErrInternalServerError,
			PrepareMock: func(mockUserDb *mocks.MockDabataseUserInterface) {
				mockUserDb.EXPECT().ReadAll(gomock.Any()).Times(1).Return(nil, echo.ErrInternalServerError)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockUserDb := mocks.NewMockDabataseUserInterface(ctrl)
			cs.PrepareMock(mockUserDb)

			app := NewAppUser(&database.Container{User: mockUserDb})

			user, err := app.ReadAll(ctx)
			if diff := cmp.Diff(user, cs.ExpectedResult); diff != "" {
				t.Error(diff)
			}

			if diff := cmp.Diff(err, cs.ExpectedErr); diff != "" {
				t.Error(diff)
			}
		})
	}
}
