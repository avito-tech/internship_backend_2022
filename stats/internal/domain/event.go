package domain

import "github.com/google/uuid"

type NewTransaction struct {
	UserID        uuid.UUID
	TransactionID uuid.UUID
}
