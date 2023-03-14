package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/amqp"
	domain "github.com/radovsky1/internship_backend_2022/balance-service/internal/domain"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/repository"
	"go.uber.org/zap"
)

type User interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByUserID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
}

type Balance interface {
	CreateBalance(ctx context.Context, balance *domain.Balance) error
	GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*domain.Balance, error)
	UpdateBalance(ctx context.Context, balanceID uuid.UUID, free, reserved float64) error
}

type Transfer interface {
	CreateTransfer(ctx context.Context, transfer *domain.Transfer) error
	GetTransferByID(ctx context.Context, transferID uuid.UUID) (*domain.Transfer, error)
	GetTransfersByFromUserID(ctx context.Context, fromUserID uuid.UUID) ([]*domain.Transfer, error)
	GetTransfersByToUserID(ctx context.Context, toUserID uuid.UUID) ([]*domain.Transfer, error)
}

type Transaction interface {
	CreateTransaction(ctx context.Context, transaction *domain.Transaction) error
	GetTransactionByID(ctx context.Context, transactionID uuid.UUID) (*domain.Transaction, error)
	UpdateTransaction(ctx context.Context, transactionID uuid.UUID, status domain.TransactionStatus) error
	GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error)
}

type ServiceUsecase interface {
	User
	Balance
	Transfer
	Transaction
}

type uc struct {
	serviceRepo repository.ServiceRepository
	queue       *amqp.Publisher
	logger      *zap.Logger
}

func New(repo repository.ServiceRepository, q *amqp.Publisher, logger *zap.Logger) ServiceUsecase {
	return &uc{
		serviceRepo: repo,
		queue:       q,
		logger:      logger,
	}
}
