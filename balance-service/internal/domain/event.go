package domain

import "github.com/google/uuid"

type NewTransaction struct {
	TransactionID uuid.UUID
	UserID        uuid.UUID
	Amount        float64
	ServiceID     uuid.UUID
	OrderID       uuid.UUID
}
