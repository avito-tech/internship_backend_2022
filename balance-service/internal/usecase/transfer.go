package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/domain"
)

func (uc uc) CreateTransfer(ctx context.Context, transfer *domain.Transfer) error {
	if transfer == nil {
		return nil
	}

	if err := uc.serviceRepo.CreateTransfer(ctx, transfer); err != nil {
		return err
	}

	return nil
}

func (uc uc) GetTransferByID(ctx context.Context, transferID uuid.UUID) (*domain.Transfer, error) {
	transfer, err := uc.serviceRepo.GetTransferByID(ctx, transferID)
	if err != nil {
		return nil, err
	}

	return transfer, nil
}

func (uc uc) GetTransfersByFromUserID(ctx context.Context, fromUserID uuid.UUID) ([]*domain.Transfer, error) {
	transfers, err := uc.serviceRepo.GetTransfersByFromUserID(ctx, fromUserID)
	if err != nil {
		return nil, err
	}

	return transfers, nil
}

func (uc uc) GetTransfersByToUserID(ctx context.Context, toUserID uuid.UUID) ([]*domain.Transfer, error) {
	transfers, err := uc.serviceRepo.GetTransfersByToUserID(ctx, toUserID)
	if err != nil {
		return nil, err
	}

	return transfers, nil
}
