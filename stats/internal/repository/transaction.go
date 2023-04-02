package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/stats/internal/domain"
)

func (rw rw) CreateTransaction(ctx context.Context, transaction *domain.TransactionByServiceID) error {
	if transaction == nil {
		return nil
	}

	if _, err := rw.store.Exec(ctx,
		`INSERT INTO transactions (id, service_id, amount, created_at) VALUES ($1, $2, $3, $4)`,
		transaction.ID, transaction.ServiceID, transaction.Amount, transaction.CreatedAt); err != nil {
		return err
	}

	return nil
}

func (rw rw) GetTransactionByID(ctx context.Context, transactionID uuid.UUID) (*domain.TransactionByServiceID, error) {
	transaction := &domain.TransactionByServiceID{}

	if err := rw.store.QueryRow(ctx,
		`SELECT id, service_id, amount, created_at FROM transactions WHERE id = $1`,
		transactionID).
		Scan(&transaction.ID, &transaction.ServiceID, &transaction.Amount,
			&transaction.CreatedAt); err != nil {
		return nil, err
	}

	return transaction, nil
}

func (rw rw) GetReportByServiceID(ctx context.Context, serviceID uuid.UUID, month int, year int) ([]*domain.ReportByServiceID, error) {
	var reports []*domain.ReportByServiceID

	rows, err := rw.store.Query(ctx,
		`SELECT service_id, SUM(amount) FROM transactions WHERE service_id = $1 AND EXTRACT(MONTH FROM created_at) = $2 AND EXTRACT(YEAR FROM created_at) = $3 GROUP BY service_id`,
		serviceID, month, year)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		report := &domain.ReportByServiceID{}

		if err := rows.Scan(&report.ServiceID, &report.Amount); err != nil {
			return nil, err
		}

		reports = append(reports, report)
	}

	return reports, nil
}
