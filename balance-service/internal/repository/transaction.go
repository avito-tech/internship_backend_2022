package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/domain"
)

func (rw rw) CreateTransaction(ctx context.Context, transaction *domain.Transaction) error {
	if transaction == nil {
		return nil
	}

	if _, err := rw.store.Exec(ctx,
		`INSERT INTO transactions (user_id, service_id, order_id, amount, status) VALUES ($1, $2, $3, $4, $5)`,
		transaction.UserID, transaction.ServiceID, transaction.OrderID, transaction.Amount, transaction.Status); err != nil {
		return err
	}

	return nil
}

func (rw rw) GetTransactionByID(ctx context.Context, transactionID uuid.UUID) (*domain.Transaction, error) {
	transaction := &domain.Transaction{}

	if err := rw.store.QueryRow(ctx,
		`SELECT id, user_id, service_id, order_id, amount, status, created_at, updated_at FROM transactions WHERE id = $1`,
		transactionID).
		Scan(&transaction.ID, &transaction.UserID, &transaction.ServiceID,
			&transaction.OrderID, &transaction.Amount, &transaction.Status,
			&transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (rw rw) UpdateTransaction(ctx context.Context, transactionID uuid.UUID, status domain.TransactionStatus) error {
	if _, err := rw.store.Exec(ctx,
		`UPDATE transactions SET status = $1, updated_at = now() WHERE id = $2`,
		status, transactionID); err != nil {
		return err
	}

	return nil
}

func (rw rw) GetTransactionsByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Transaction, error) {
	var transactions []*domain.Transaction

	rows, err := rw.store.Query(ctx,
		`SELECT id, user_id, service_id, order_id, amount, status, created_at, updated_at FROM transactions WHERE user_id = $1`,
		userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		transaction := &domain.Transaction{}

		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ServiceID,
			&transaction.OrderID, &transaction.Amount, &transaction.Status,
			&transaction.CreatedAt, &transaction.UpdatedAt); err != nil {
			return nil, err
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}
