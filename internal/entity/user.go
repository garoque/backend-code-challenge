package entity

import (
	"sync"
	"time"

	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/garoque/backend-code-challenge-snapfi/pkg/uuid"
)

type User struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Balance   float64     `json:"balance"`
	Mutex     *sync.Mutex `json:"-"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt *time.Time  `json:"updatedAt,omitempty" db:"updated_at"`
}

func NewUser(user dto.CreateUser) *User {
	return &User{
		ID:    uuid.NewId(),
		Name:  user.Name,
		Mutex: &sync.Mutex{},
	}
}
