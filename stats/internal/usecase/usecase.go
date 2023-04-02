package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/stats/internal/domain"
	"github.com/radovsky1/internship_backend_2022/stats/internal/repository"
	"go.uber.org/zap"
)

type Transaction interface {
	GetReportByServiceID(ctx context.Context, serviceID uuid.UUID, month int, year int) ([]*domain.ReportByServiceID, error)
	TransactionHandler(ctx context.Context, transaction domain.Transaction) error
}

type ServiceUsecase interface {
	Transaction
}

type uc struct {
	serviceRepo repository.ServiceRepository
	logger      *zap.Logger
}

func New(repo repository.ServiceRepository, logger *zap.Logger) ServiceUsecase {
	return uc{
		serviceRepo: repo,
		logger:      logger,
	}
}
