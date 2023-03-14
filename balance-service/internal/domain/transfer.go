package domain

import (
	"github.com/google/uuid"
	"time"
)

type Transfer struct {
	ID         uuid.UUID
	FromUserID uuid.UUID
	ToUserID   uuid.UUID
	Amount     float64
	CreatedAt  time.Time
}
