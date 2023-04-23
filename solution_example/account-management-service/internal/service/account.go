package service

import (
	"account-management-service/internal/entity"
	"account-management-service/internal/repo"
	"account-management-service/internal/repo/repoerrs"
	"context"
)

type AccountService struct {
	accountRepo repo.Account
}

func NewAccountService(accountRepo repo.Account) *AccountService {
	return &AccountService{accountRepo: accountRepo}
}

func (s *AccountService) CreateAccount(ctx context.Context) (int, error) {
	id, err := s.accountRepo.CreateAccount(ctx)
	if err != nil {
		if err == repoerrs.ErrAlreadyExists {
			return 0, ErrAccountAlreadyExists
		}
		return 0, ErrCannotCreateAccount
	}

	return id, nil
}

func (s *AccountService) GetAccountById(ctx context.Context, userId int) (entity.Account, error) {
	return s.accountRepo.GetAccountById(ctx, userId)
}

func (s *AccountService) Deposit(ctx context.Context, input AccountDepositInput) error {
	return s.accountRepo.Deposit(ctx, input.Id, input.Amount)
}

func (s *AccountService) Withdraw(ctx context.Context, input AccountWithdrawInput) error {
	return s.accountRepo.Withdraw(ctx, input.Id, input.Amount)
}

func (s *AccountService) Transfer(ctx context.Context, input AccountTransferInput) error {
	return s.accountRepo.Transfer(ctx, input.From, input.To, input.Amount)
}
