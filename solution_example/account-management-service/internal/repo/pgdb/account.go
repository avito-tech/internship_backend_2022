package pgdb

import (
	"account-management-service/internal/entity"
	"account-management-service/internal/repo/repoerrs"
	"account-management-service/pkg/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	log "github.com/sirupsen/logrus"
)

type AccountRepo struct {
	*postgres.Postgres
}

func NewAccountRepo(pg *postgres.Postgres) *AccountRepo {
	return &AccountRepo{pg}
}

func (r *AccountRepo) CreateAccount(ctx context.Context) (int, error) {
	sql, args, _ := r.Builder.
		Insert("accounts").
		Values(squirrel.Expr("DEFAULT")).
		Suffix("RETURNING id").
		ToSql()

	var id int
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		log.Debugf("err: %v", err)
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repoerrs.ErrAlreadyExists
			}
		}
		return 0, fmt.Errorf("AccountRepo.CreateAccount - r.Pool.QueryRow: %v", err)
	}

	return id, nil
}

func (r *AccountRepo) GetAccountById(ctx context.Context, id int) (entity.Account, error) {
	sql, args, _ := r.Builder.
		Select("*").
		From("accounts").
		Where("id = ?", id).
		ToSql()

	var account entity.Account
	err := r.Pool.QueryRow(ctx, sql, args...).Scan(
		&account.Id,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Account{}, repoerrs.ErrNotFound
		}
		return entity.Account{}, fmt.Errorf("AccountRepo.GetAccountById - r.Pool.QueryRow: %v", err)
	}

	return account, nil
}

func (r *AccountRepo) Deposit(ctx context.Context, id, amount int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Deposit - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	sql, args, _ := r.Builder.
		Update("accounts").
		Set("balance", squirrel.Expr("balance + ?", amount)).
		Where("id = ?", id).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Deposit - tx.Exec: %v", err)
	}

	sql, args, _ = r.Builder.
		Insert("operations").
		Columns("account_id", "amount", "operation_type").
		Values(id, amount, entity.OperationTypeDeposit).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Deposit - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Deposit - tx.Commit: %v", err)
	}

	return nil
}

func (r *AccountRepo) Withdraw(ctx context.Context, id, amount int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Withdraw - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// check if account has enough balance to withdraw
	sql, args, _ := r.Builder.
		Select("balance").
		From("accounts").
		Where("id = ?", id).
		ToSql()

	var balance int
	err = tx.QueryRow(ctx, sql, args...).Scan(&balance)
	if err != nil {
		return fmt.Errorf("AccountRepo.Withdraw - tx.QueryRow: %v", err)
	}

	if balance < amount {
		return repoerrs.ErrNotEnoughBalance
	}

	sql, args, _ = r.Builder.
		Update("accounts").
		Set("balance", squirrel.Expr("balance - ?", amount)).
		Where("id = ?", id).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Withdraw - tx.Exec: %v", err)
	}

	sql, args, _ = r.Builder.
		Insert("operations").
		Columns("account_id", "amount", "operation_type").
		Values(id, amount, entity.OperationTypeWithdraw).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Withdraw - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Withdraw - tx.Commit: %v", err)
	}

	return nil
}

func (r *AccountRepo) Transfer(ctx context.Context, from, to, amount int) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - r.Pool.Begin: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	// check if account 'from' has enough balance to transfer
	sql, args, _ := r.Builder.
		Select("balance").
		From("accounts").
		Where("id = ?", from).
		ToSql()

	var balance int
	err = tx.QueryRow(ctx, sql, args...).Scan(&balance)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - tx.QueryRow: %v", err)
	}

	if balance < amount {
		return repoerrs.ErrNotEnoughBalance
	}

	sql, args, _ = r.Builder.
		Update("accounts").
		Set("balance", squirrel.Expr("balance - ?", amount)).
		Where("id = ?", from).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - tx.Exec: %v", err)
	}

	sql, args, _ = r.Builder.
		Update("accounts").
		Set("balance", squirrel.Expr("balance + ?", amount)).
		Where("id = ?", to).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - tx.Exec: %v", err)
	}

	sql, args, _ = r.Builder.
		Insert("operations").
		Columns("account_id", "amount", "operation_type").
		Values(from, amount, entity.OperationTypeTransferFrom).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - tx.Exec: %v", err)
	}

	sql, args, _ = r.Builder.
		Insert("operations").
		Columns("account_id", "amount", "operation_type").
		Values(to, amount, entity.OperationTypeTransferTo).
		ToSql()

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - tx.Exec: %v", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("AccountRepo.Transfer - tx.Commit: %v", err)
	}

	return nil
}
