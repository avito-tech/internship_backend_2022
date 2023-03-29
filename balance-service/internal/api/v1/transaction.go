package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/domain"
	"net/http"
	"strconv"
)

type createTransactionRequest struct {
	UserID    string  `json:"user_id" binding:"required"`
	OrderID   string  `json:"order_id" binding:"required"`
	ServiceID string  `json:"service_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required"`
}

type createTransactionResponse struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (rH RouterHandler) createTransaction(c *gin.Context) {
	var req createTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	orderID, err := uuid.Parse(req.OrderID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	transaction := domain.Transaction{
		UserID:    userID,
		OrderID:   orderID,
		ServiceID: serviceID,
		Amount:    req.Amount,
	}

	if err := rH.ucService.CreateTransaction(c, &transaction); err != nil {
		_ = c.Error(err)
		return
	}

	resp := createTransactionResponse{
		ID:        transaction.ID.String(),
		Status:    string(transaction.Status),
		CreatedAt: transaction.CreatedAt.String(),
		UpdatedAt: transaction.UpdatedAt.String(),
	}

	c.JSON(http.StatusOK, resp)
}

type getTransactionRequest struct {
	ID string `uri:"id" binding:"required"`
}

type getTransactionResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	OrderID   string `json:"order_id"`
	ServiceID string `json:"service_id"`
	Amount    string `json:"amount"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (rH RouterHandler) getTransaction(c *gin.Context) {
	var req getTransactionRequest

	if err := c.ShouldBindUri(&req); err != nil {
		_ = c.Error(err)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	transaction, err := rH.ucService.GetTransactionByID(c, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := getTransactionResponse{
		ID:        transaction.ID.String(),
		UserID:    transaction.UserID.String(),
		OrderID:   transaction.OrderID.String(),
		ServiceID: transaction.ServiceID.String(),
		Amount:    strconv.FormatFloat(transaction.Amount, 'f', 2, 64),
		Status:    string(transaction.Status),
		CreatedAt: transaction.CreatedAt.String(),
		UpdatedAt: transaction.UpdatedAt.String(),
	}

	c.JSON(http.StatusOK, resp)
}

type updateTransactionRequest struct {
	ID     string `json:"id" binding:"required"`
	Status string `json:"status" binding:"required"`
}

type updateTransactionResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	OrderID   string `json:"order_id"`
	ServiceID string `json:"service_id"`
	Amount    string `json:"amount"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func (rH RouterHandler) updateTransaction(c *gin.Context) {
	var req updateTransactionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	id, err := uuid.Parse(req.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = rH.ucService.UpdateTransaction(c, id, domain.TransactionStatus(req.Status))
	if err != nil {
		_ = c.Error(err)
		return
	}

	transaction, err := rH.ucService.GetTransactionByID(c, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := updateTransactionResponse{
		ID:        transaction.ID.String(),
		UserID:    transaction.UserID.String(),
		OrderID:   transaction.OrderID.String(),
		ServiceID: transaction.ServiceID.String(),
		Amount:    strconv.FormatFloat(transaction.Amount, 'f', 2, 64),
		Status:    string(transaction.Status),
		CreatedAt: transaction.CreatedAt.String(),
		UpdatedAt: transaction.UpdatedAt.String(),
	}

	c.JSON(http.StatusOK, resp)
}

type getTransactionsByUserIDRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type getTransactionsByUserIDResponse struct {
	Transactions []getTransactionResponse `json:"transactions"`
}

func (rH RouterHandler) getTransactionsByUserID(c *gin.Context) {
	var req getTransactionsByUserIDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	transactions, err := rH.ucService.GetTransactionsByUserID(c, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	resp := getTransactionsByUserIDResponse{}
	for _, transaction := range transactions {
		t := getTransactionResponse{
			ID:        transaction.ID.String(),
			UserID:    transaction.UserID.String(),
			OrderID:   transaction.OrderID.String(),
			ServiceID: transaction.ServiceID.String(),
			Amount:    strconv.FormatFloat(transaction.Amount, 'f', 2, 64),
			Status:    string(transaction.Status),
			CreatedAt: transaction.CreatedAt.String(),
			UpdatedAt: transaction.UpdatedAt.String(),
		}
		resp.Transactions = append(resp.Transactions, t)
	}

	c.JSON(http.StatusOK, resp)
}
