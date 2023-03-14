package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/domain"
)

func (rw rw) CreateBalance(ctx context.Context, balance *domain.Balance) error {
	if balance == nil {
		return nil
	}

	if _, err := rw.store.Exec(
		ctx,
		`INSERT INTO 
				balances (user_id, free, reserved) VALUES ($1, $2, $3)`, balance.UserID, balance.Free, balance.Reserved); err != nil {
		return err
	}

	return nil
}

func (rw rw) GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*domain.Balance, error) {
	balance := &domain.Balance{}

	if err := rw.store.QueryRow(
		ctx,
		`SELECT id, user_id, free, reserved FROM balances WHERE user_id = $1`, userID).Scan(&balance.ID, &balance.UserID, &balance.Free, &balance.Reserved); err != nil {
		return nil, err
	}

	return balance, nil
}

func (rw rw) UpdateBalance(ctx context.Context, balanceID uuid.UUID, free, reserved float64) error {
	if _, err := rw.store.Exec(
		ctx,
		`UPDATE balances SET free=$2, reserved=$3 WHERE id=$1`, balanceID, free, reserved); err != nil {
		return err
	}

	return nil
}

func (rw rw) Reserve(ctx context.Context, balanceID uuid.UUID, amount float64) error {
	if _, err := rw.store.Exec(
		ctx,
		`UPDATE balances SET reserved=reserved+$2, free=free-$2 WHERE id=$1`, balanceID, amount); err != nil {
		return err
	}

	return nil
}
