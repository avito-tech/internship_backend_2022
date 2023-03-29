package main

import (
	"context"
	"errors"
	"github.com/radovsky1/internship_backend_2022/balance-service/cfg"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/amqp"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/api"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/repository"
	"github.com/radovsky1/internship_backend_2022/balance-service/internal/usecase"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c, err := cfg.Read()
	if err != nil {
		log.Println("failed to read config:", err.Error())

		return
	}

	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()

	dbPool, err := repository.Connect(c)
	if err != nil {
		logger.Error(err.Error())
	}

	defer func() {
		if dbPool != nil {
			dbPool.Close()
		}
	}()

	publisher, err := amqp.BuildPublisher(c)
	if err != nil {
		logger.Error(err.Error())
	}

	defer func() {
		if publisher != nil {
			_ = publisher.Close()
		}
	}()

	uc := usecase.New(
		repository.New(dbPool),
		publisher,
		logger,
	)

	httpAPIRouter := api.New(
		logger,
		dbPool,
		uc,
	)

	httpServer := http.Server{
		Addr:              c.HttpAPI.Addr,
		Handler:           httpAPIRouter,
		ReadHeaderTimeout: time.Minute * 1,
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

	<-quit

	ctxGIN, cancelGIN := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelGIN()

	if err := httpServer.Shutdown(ctxGIN); err != nil {
		logger.Error("http server: forced to shutdown", zap.Error(err))
	}

	logger.Info("http server: successfully stopped")
}
