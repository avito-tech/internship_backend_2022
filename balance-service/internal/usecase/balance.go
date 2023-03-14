package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/domain"
)

func (uc uc) CreateBalance(ctx context.Context, balance *domain.Balance) error {
	if balance == nil {
		return nil
	}

	if err := uc.serviceRepo.CreateBalance(ctx, balance); err != nil {
		return err
	}

	return nil
}

func (uc uc) GetBalanceByUserID(ctx context.Context, userID uuid.UUID) (*domain.Balance, error) {
	balance, err := uc.serviceRepo.GetBalanceByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

func (uc uc) UpdateBalance(ctx context.Context, balanceID uuid.UUID, free, reserved float64) error {
	if err := uc.serviceRepo.UpdateBalance(ctx, balanceID, free, reserved); err != nil {
		return err
	}

	return nil
}
