package service

import (
	"account-management-service/internal/entity"
	"account-management-service/internal/repo"
	"account-management-service/internal/webapi"
	"account-management-service/pkg/hasher"
	"context"
	"time"
)

type AuthCreateUserInput struct {
	Username string
	Password string
}

type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type Auth interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput) (int, error)
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(token string) (int, error)
}

type AccountDepositInput struct {
	Id     int
	Amount int
}

type AccountWithdrawInput struct {
	Id     int
	Amount int
}

type AccountTransferInput struct {
	From   int
	To     int
	Amount int
}

type Account interface {
	CreateAccount(ctx context.Context) (int, error)
	GetAccountById(ctx context.Context, userId int) (entity.Account, error)
	Deposit(ctx context.Context, input AccountDepositInput) error
	Withdraw(ctx context.Context, input AccountWithdrawInput) error
	Transfer(ctx context.Context, input AccountTransferInput) error
}

type Product interface {
	CreateProduct(ctx context.Context, name string) (int, error)
	GetProductById(ctx context.Context, id int) (entity.Product, error)
}

type ReservationCreateInput struct {
	AccountId int
	ProductId int
	OrderId   int
	Amount    int
}

type Reservation interface {
	CreateReservation(ctx context.Context, input ReservationCreateInput) (int, error)
	RefundReservationByOrderId(ctx context.Context, orderId int) error
	RevenueReservationByOrderId(ctx context.Context, orderId int) error
}

type OperationHistoryInput struct {
	AccountId int
	SortType  string
	Offset    int
	Limit     int
}

type OperationHistoryOutput struct {
	Amount      int       `json:"amount"`
	Operation   string    `json:"operation"`
	Time        time.Time `json:"time"`
	Product     string    `json:"product,omitempty"`
	Order       *int      `json:"order,omitempty"`
	Description string    `json:"description,omitempty"`
}

type Operation interface {
	OperationHistory(ctx context.Context, input OperationHistoryInput) ([]OperationHistoryOutput, error)
	MakeReportLink(ctx context.Context, month, year int) (string, error)
	MakeReportFile(ctx context.Context, month, year int) ([]byte, error)
}

type Services struct {
	Auth        Auth
	Account     Account
	Product     Product
	Reservation Reservation
	Operation   Operation
}

type ServicesDependencies struct {
	Repos  *repo.Repositories
	GDrive webapi.GDrive
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth:        NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
		Account:     NewAccountService(deps.Repos.Account),
		Product:     NewProductService(deps.Repos.Product),
		Reservation: NewReservationService(deps.Repos.Reservation),
		Operation:   NewOperationService(deps.Repos.Operation, deps.Repos.Product, deps.GDrive),
	}
}
