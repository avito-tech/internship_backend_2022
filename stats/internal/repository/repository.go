package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/radovsky1/internship_backend_2022/stats/internal/domain"
)

type rw struct {
	store *pgxpool.Pool
}

type Transaction interface {
	CreateTransaction(ctx context.Context, transaction *domain.TransactionByServiceID) error
	GetTransactionByID(ctx context.Context, transactionID uuid.UUID) (*domain.TransactionByServiceID, error)
	GetReportByServiceID(ctx context.Context, serviceID uuid.UUID, month int, year int) ([]*domain.ReportByServiceID, error)
}

type ServiceRepository interface {
	Transaction
}

func New(dbPool *pgxpool.Pool) ServiceRepository {
	return rw{
		store: dbPool,
	}
}
