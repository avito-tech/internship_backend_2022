package main

import (
	"context"
	"errors"
	config "github.com/radovsky1/internship_backend_2022/stats/cfg"
	"github.com/radovsky1/internship_backend_2022/stats/internal/amqp"
	"github.com/radovsky1/internship_backend_2022/stats/internal/api"
	"github.com/radovsky1/internship_backend_2022/stats/internal/repository"
	"github.com/radovsky1/internship_backend_2022/stats/internal/usecase"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Println("failed to read config:", err.Error())

		return
	}

	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	dbPool, err := repository.Connect(cfg)
	if err != nil {
		logger.Error(err.Error())
	}

	defer func() {
		if dbPool != nil {
			dbPool.Close()
		}
	}()

	uc := usecase.New(
		repository.New(dbPool),
		logger,
	)

	httpApiRouter := api.New(logger, dbPool, uc)

	httpServer := http.Server{
		Addr:              cfg.HTTPAPI.Addr,
		Handler:           httpApiRouter,
		ReadHeaderTimeout: time.Minute * 1,
	}

	consumer, err := amqp.BuildConsumer(cfg)
	if err != nil {
		logger.Error(err.Error())
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("http server: failed to listen and serve", zap.Error(err))

			quit <- syscall.SIGTERM
		}
	}()

	logger.Info("http server: ready to accept requests")

	go func() {
		if err := consumer.StartConsume(logger, uc); err != nil {
			logger.Error("feed worker: failed to start consume", zap.Error(err))
			quit <- syscall.SIGTERM
		}
	}()

	logger.Info("feed worker: ready to accept msg")

	<-quit

	ctxGIN, cancelGIN := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelGIN()

	if err := httpServer.Shutdown(ctxGIN); err != nil {
		logger.Error("http server: forced to shutdown", zap.Error(err))
	}

	logger.Info("http server: successfully stopped")

	ctxConsumer, cancelConsumer := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelConsumer()

	if err := consumer.Close(ctxConsumer); err != nil {
		logger.Error("http server: forced to shutdown", zap.Error(err))
	}

	logger.Info("feed worker: successfully stopped")
}
