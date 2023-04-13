package user

import (
	"context"
	"log"

	"github.com/garoque/backend-code-challenge-snapfi/internal/database"
	"github.com/garoque/backend-code-challenge-snapfi/internal/entity"
)

type AppUserInterface interface {
	Create(ctx context.Context, user entity.User) error
	ReadOneById(ctx context.Context, userId string) (*entity.User, error)
	ReadAll(ctx context.Context) ([]entity.User, error)
}

type appUserImpl struct {
	db *database.Container
}

func NewAppUser(db *database.Container) AppUserInterface {
	return &appUserImpl{db}
}

func (u *appUserImpl) Create(ctx context.Context, user entity.User) error {
	err := u.db.User.Create(ctx, user)
	if err != nil {
		log.Println("Error app.user.Create.db.Create: ", err.Error())
		return err
	}

	return nil
}

func (u *appUserImpl) ReadOneById(ctx context.Context, userId string) (*entity.User, error) {
	user, err := u.db.User.ReadOneById(ctx, userId)
	if err != nil {
		log.Println("Error app.user.ReadOneById.db.ReadOneById: ", err.Error())
		return nil, err
	}

	return user, nil
}

func (u *appUserImpl) ReadAll(ctx context.Context) ([]entity.User, error) {
	users, err := u.db.User.ReadAll(ctx)
	if err != nil {
		log.Println("Error app.user.ReadAll.db.ReadAll: ", err.Error())
		return nil, err
	}

	return users, nil
}
