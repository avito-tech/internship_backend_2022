package api

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	v1 "github.com/radovsky1/internship_backend_2022/stats/internal/api/v1"
	"github.com/radovsky1/internship_backend_2022/stats/internal/api/v1/formatter"
	"github.com/radovsky1/internship_backend_2022/stats/internal/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"time"
)

func New(
	log *zap.Logger,
	db *pgxpool.Pool,
	ucService usecase.ServiceUsecase,
) *gin.Engine {
	router := newGINRouter(log)

	registerSwagger(router)
	registerHealthCheck(router, db)

	routerHandler := v1.NewRouter(ucService, log)

	api := router.Group("/api/v1")
	routerHandler.RegisterRoutes(api)

	return router
}

func newGINRouter(log *zap.Logger) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	router.Use(ginzap.RecoveryWithZap(log, true))
	router.Use(ginzap.Ginzap(log, time.RFC3339, true))

	router.Use(formatter.HandleErrors)

	return router
}

func registerHealthCheck(router *gin.Engine, db *pgxpool.Pool) {
	handler := newHealthCheckHandler(db)

	router.GET("/health", handler.Health)
	router.GET("/ready", handler.Ready)
}

func registerSwagger(router *gin.Engine) {
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
