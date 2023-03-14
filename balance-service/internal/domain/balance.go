package domain

import "github.com/google/uuid"

type Balance struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Free     float64
	Reserved float64
}
