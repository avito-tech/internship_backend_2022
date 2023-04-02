package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/stats/internal/domain"
)

func (uc uc) GetReportByServiceID(ctx context.Context, serviceID uuid.UUID, month int, year int) ([]*domain.ReportByServiceID, error) {
	report, err := uc.serviceRepo.GetReportByServiceID(ctx, serviceID, month, year)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (uc uc) TransactionHandler(ctx context.Context, transaction domain.Transaction) error {
	t := &domain.TransactionByServiceID{
		ID:        transaction.ID,
		ServiceID: transaction.ServiceID,
		Amount:    transaction.Amount,
		CreatedAt: transaction.CreatedAt,
	}

	if err := uc.serviceRepo.CreateTransaction(ctx, t); err != nil {
		return err
	}

	return nil
}
