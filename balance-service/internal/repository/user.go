package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/domain"
)

func (rw rw) CreateUser(ctx context.Context, u *domain.User) error {
	if u == nil {
		return nil
	}

	if _, err := rw.store.Exec(
		ctx,
		`INSERT INTO 
    			users (user_id) VALUES ($1)`, u.UserID); err != nil {
		return err
	}

	return nil
}

func (rw rw) GetUserByUserID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user := &domain.User{}

	if err := rw.store.QueryRow(
		ctx,
		`SELECT id, user_id FROM users WHERE user_id = $1`, userID).Scan(&user.ID, &user.UserID); err != nil {
		return nil, err
	}

	return user, nil
}
