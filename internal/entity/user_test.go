package entity

import (
	"testing"

	"github.com/garoque/backend-code-challenge-snapfi/internal/api/dto"
	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user := NewUser(dto.CreateUser{Name: "Gabriel"})
	assert.NotNil(t, user)
	assert.NotNil(t, user.Mutex)
	assert.Equal(t, "Gabriel", user.Name)
}
