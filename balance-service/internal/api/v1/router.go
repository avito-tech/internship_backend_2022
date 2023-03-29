package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/usecase"
	"go.uber.org/zap"
)

type RouterHandler struct {
	ucService usecase.ServiceUsecase
	logger    *zap.Logger
}

func NewRouter(uc usecase.ServiceUsecase, logger *zap.Logger) *RouterHandler {
	return &RouterHandler{
		ucService: uc,
		logger:    logger,
	}
}

func (rH RouterHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/balance", rH.createBalance)
	router.GET("/balance/:user_id", rH.getBalanceByUserID)
	router.PUT("/balance/:balance_id", rH.updateBalance)

	//router.POST("/transfer", rH.createTransfer)
	//router.GET("/transfer/:transfer_id", rH.getTransferByID)
	//router.GET("/transfer/from/:from_user_id", rH.getTransfersByFromUserID)
	//router.GET("/transfer/to/:to_user_id", rH.getTransfersByToUserID)

	router.POST("/transaction", rH.createTransaction)
	router.GET("/transaction/:transaction_id", rH.getTransaction)
	router.GET("/transaction/user/:user_id", rH.getTransactionsByUserID)
	router.PUT("/transaction/:transaction_id", rH.updateTransaction)
}
