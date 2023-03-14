package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/domain"
)

func (rw rw) CreateTransfer(ctx context.Context, t *domain.Transfer) error {
	if t == nil {
		return nil
	}

	tx, err := rw.store.Begin(ctx)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx,
		`UPDATE balances SET free=free-$2 WHERE id=$1`, t.FromUserID, t.Amount)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}

	_, err = tx.Exec(ctx,
		`UPDATE balances SET free=free+$2 WHERE id=$1`, t.ToUserID, t.Amount)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}

	_, err = tx.Exec(ctx,
		`INSERT transfers (id, from_user_id, to_user_id, amount) VALUES ($1, $2, $3, $4)`,
		t.ID, t.FromUserID, t.ToUserID, t.Amount)
	if err != nil {
		err := tx.Rollback(ctx)
		if err != nil {
			return err
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (rw rw) GetTransferByID(ctx context.Context, transferID uuid.UUID) (*domain.Transfer, error) {
	transfer := &domain.Transfer{}

	if err := rw.store.QueryRow(
		ctx,
		`SELECT id, from_user_id, to_user_id, amount, created_at FROM transfers WHERE id = $1`, transferID).
		Scan(&transfer.ID, &transfer.FromUserID, &transfer.ToUserID, &transfer.Amount, &transfer.CreatedAt); err != nil {
		return nil, err
	}

	return transfer, nil
}

func (rw rw) GetTransfersByFromUserID(ctx context.Context, fromUserID uuid.UUID) ([]*domain.Transfer, error) {
	transfers := []*domain.Transfer{}

	rows, err := rw.store.Query(
		ctx,
		`SELECT id, from_user_id, to_user_id, amount, created_at FROM transfers WHERE from_user_id = $1`,
		fromUserID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		transfer := &domain.Transfer{}

		if err := rows.Scan(&transfer.ID, &transfer.FromUserID, &transfer.ToUserID, &transfer.Amount, &transfer.CreatedAt); err != nil {
			return nil, err
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}

func (rw rw) GetTransfersByToUserID(ctx context.Context, toUserID uuid.UUID) ([]*domain.Transfer, error) {
	var transfers []*domain.Transfer

	rows, err := rw.store.Query(
		ctx,
		`SELECT id, from_user_id, to_user_id, amount, created_at FROM transfers WHERE to_user_id = $1`,
		toUserID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		transfer := &domain.Transfer{}

		if err := rows.Scan(&transfer.ID, &transfer.FromUserID, &transfer.ToUserID, &transfer.Amount, &transfer.CreatedAt); err != nil {
			return nil, err
		}

		transfers = append(transfers, transfer)
	}

	return transfers, nil
}
