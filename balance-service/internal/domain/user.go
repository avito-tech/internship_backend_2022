package domain

import "github.com/google/uuid"

type User struct {
	ID     uuid.UUID
	UserID uuid.UUID
}
