package usecase

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/amqp"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/domain"
	"time"
)

func (uc uc) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	if transaction == nil {
		return nil
	}

	if err := uc.serviceRepo.CreateTransaction(ctx, transaction); err != nil {
		return err
	}

	balance, err := uc.serviceRepo.GetBalanceByUserID(ctx, transaction.UserID)
	if err != nil {
		return err
	}

	err = uc.serviceRepo.Reserve(ctx, balance.ID, transaction.Amount)
	if err != nil {
		return err
	}

	return nil
}

func (uc uc) GetTransactionByID(ctx context.Context, transactionID uuid.UUID) (*domain.Transaction, error) {
	transaction, err := uc.serviceRepo.GetTransactionByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (uc uc) UpdateTransaction(ctx context.Context, transactionID uuid.UUID, status domain.TransactionStatus) error {
	if err := uc.serviceRepo.UpdateTransaction(ctx, transactionID, status); err != nil {
		return err
	}

	if status == domain.Accepted {
		t, err := uc.serviceRepo.GetTransactionByID(ctx, transactionID)
		if err != nil {
			return err
		}
		go func() {
			err := uc.queue.Push(&amqp.Message{
				Data: domain.NewTransaction{
					TransactionID: transactionID,
					UserID:        t.UserID,
					Amount:        t.Amount,
					ServiceID:     t.ServiceID,
					OrderID:       t.OrderID,
				},
				Timestamp: time.Now(),
				Key:       amqp.TransactionEvent,
			})
			if err != nil {
				uc.logger.Error(fmt.Sprintf("error push msg new post: %+v", err))
			}
		}()
	}

	return nil
}

func (uc uc) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error) {
	transactions, err := uc.serviceRepo.GetTransactionsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}
