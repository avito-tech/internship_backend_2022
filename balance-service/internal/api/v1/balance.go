package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/domain"
	"net/http"
	"strconv"
)

type createBalanceRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Free     string `json:"free" binding:"required"`
	Reserved string `json:"reserved"`
}

type createBalanceResponse struct {
	ID string `json:"id"`
}

func (rH RouterHandler) createBalance(c *gin.Context) {
	var req createBalanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	free, err := strconv.ParseFloat(req.Free, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}

	reserved, err := strconv.ParseFloat(req.Reserved, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}

	balance := domain.Balance{
		UserID:   userID,
		Free:     free,
		Reserved: reserved,
	}

	if err := rH.ucService.CreateBalance(c, &balance); err != nil {
		_ = c.Error(err)
		return
	}

	user := domain.User{
		UserID: userID,
	}

	if err := rH.ucService.CreateUser(c, &user); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, createBalanceResponse{
		ID: balance.ID.String(),
	})
}

type getBalanceByUserIDRequest struct {
	UserID string `json:"user_id" binding:"required"`
}

type getBalanceByUserIDResponse struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	Free     float64 `json:"free"`
	Reserved float64 `json:"reserved"`
}

func (rH RouterHandler) getBalanceByUserID(c *gin.Context) {
	var req getBalanceByUserIDRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	balance, err := rH.ucService.GetBalanceByUserID(c, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, getBalanceByUserIDResponse{
		ID:       balance.ID.String(),
		UserID:   balance.UserID.String(),
		Free:     balance.Free,
		Reserved: balance.Reserved,
	})
}

type updateBalanceRequest struct {
	ID       string `json:"id" binding:"required"`
	Free     string `json:"free" binding:"required"`
	Reserved string `json:"reserved"`
}

type updateBalanceResponse struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	Free     float64 `json:"free"`
	Reserved float64 `json:"reserved"`
}

func (rH RouterHandler) updateBalance(c *gin.Context) {
	var req updateBalanceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		_ = c.Error(err)
		return
	}

	balanceID, err := uuid.Parse(req.ID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	free, err := strconv.ParseFloat(req.Free, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}

	reserved, err := strconv.ParseFloat(req.Reserved, 64)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if err := rH.ucService.UpdateBalance(c, balanceID, free, reserved); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, updateBalanceResponse{
		ID:       balanceID.String(),
		Free:     free,
		Reserved: reserved,
	})
}
