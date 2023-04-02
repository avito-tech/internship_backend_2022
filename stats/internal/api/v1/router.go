package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/radovsky1/internship_backend_2022/stats/internal/usecase"
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
	router.GET("/report", rH.getReport)
}
