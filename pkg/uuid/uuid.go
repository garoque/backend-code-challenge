package uuid

import "github.com/google/uuid"

func NewId() string {
	return uuid.New().String()
}
