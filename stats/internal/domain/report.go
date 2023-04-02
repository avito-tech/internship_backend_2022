package domain

import "github.com/google/uuid"

type ReportByServiceID struct {
	ServiceID uuid.UUID
	Amount    float64
}
