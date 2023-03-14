package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	domain2 "github.com/radovsky1/internship_backend_2022/balance-service/internal/domain"
)

type rw struct {
	store *pgxpool.Pool
}

type User interface {
	CreateUser(ctx context.Context, user *domain2.User) error
	GetUserByUserID(ctx context.Context, userID uuid.UUID) (*domain2.User, error)
}

type Balance interface {
	CreateBalance(ctx context.Context, balance *domain2.Balance) error
	GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*domain2.Balance, error)
	UpdateBalance(ctx context.Context, balanceID uuid.UUID, free, reserved float64) error
	Reserve(ctx context.Context, balanceID uuid.UUID, amount float64) error
}

type Transfer interface {
	CreateTransfer(ctx context.Context, transfer *domain2.Transfer) error
	GetTransferByID(ctx context.Context, transferID uuid.UUID) (*domain2.Transfer, error)
	GetTransfersByFromUserID(ctx context.Context, fromUserID uuid.UUID) ([]*domain2.Transfer, error)
	GetTransfersByToUserID(ctx context.Context, toUserID uuid.UUID) ([]*domain2.Transfer, error)
}

type Transaction interface {
	CreateTransaction(ctx context.Context, transaction *domain2.Transaction) error
	GetTransactionByID(ctx context.Context, transactionID uuid.UUID) (*domain2.Transaction, error)
	UpdateTransaction(ctx context.Context, transactionID uuid.UUID, status domain2.TransactionStatus) error
	GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain2.Transaction, error)
}

// go:generate mockery --name ServiceRepository
type ServiceRepository interface {
	User
	Balance
	Transfer
	Transaction
}

func New(dbPool *pgxpool.Pool) ServiceRepository {
	return rw{
		store: dbPool,
	}
}
