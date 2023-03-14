package usecase

import (
	"context"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/domain"
)

func (uc uc) CreateUser(ctx context.Context, user *domain.User) error {
	if user == nil {
		return nil
	}

	if err := uc.serviceRepo.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (uc uc) GetUserByUserID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := uc.serviceRepo.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
