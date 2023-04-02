package domain

import (
	"github.com/google/uuid"
	"time"
)

type TransactionStatus string

const (
	Processed TransactionStatus = "processed"
	Accepted  TransactionStatus = "accepted"
	Rejected  TransactionStatus = "rejected"
)

type Transaction struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ServiceID uuid.UUID
	OrderID   uuid.UUID
	Amount    float64
	Status    TransactionStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionByServiceID struct {
	ID        uuid.UUID
	ServiceID uuid.UUID
	Amount    float64
	CreatedAt time.Time
}
