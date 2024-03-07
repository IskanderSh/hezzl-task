package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IskanderSh/hezzl-task/internal/clients"
	"github.com/IskanderSh/hezzl-task/internal/config"
	"github.com/IskanderSh/hezzl-task/internal/handlers"
	"github.com/IskanderSh/hezzl-task/internal/services"
	redis "github.com/IskanderSh/hezzl-task/internal/storage/cache"
	"github.com/IskanderSh/hezzl-task/internal/storage/postgres"

	//redis "github.com/IskanderSh/hezzl-task/internal/storage/cache"
	"github.com/IskanderSh/hezzl-task/internal/storage/clickhouse"
)

type Server struct {
	HTTPServer *http.Server
}

func NewServer(log *slog.Logger, cfg *config.Config) *Server {
	//Storages
	storage, err := postgres.NewStorage(log, cfg.Storage)
	if err != nil {
		panic(err)
	}
	log.Info("successfully create connection to storage")

	cache := redis.NewCache(cfg.Cache)
	log.Info("successfully create connection to cache")

	ctx := context.Background()
	logStorage, err := clickhouse.NewLogStorage(ctx, log, cfg.LogStorage)
	if err != nil {
		panic(err)
	}
	log.Info("successfully create connection to log storage")

	// Clients
	brokerClient, err := clients.NewNatsClient(log, cfg.MessageBroker, logStorage)
	if err != nil {
		panic(err)
	}
	log.Info("successfully create connection to message broker")

	// Subscribe to subject
	go func() {
		if err := brokerClient.SubscribeSubjects(); err != nil {
			log.Error(err.Error())
		}
	}()

	// Services
	brokerServer, err := services.NewNatsServer(log, cfg.MessageBroker)
	if err != nil {
		panic(err)
	}

	goodService := services.NewGoodService(log, storage, cache, brokerServer)

	// Handlers
	handler := handlers.NewGoodHandler(log, goodService)

	// Router
	router := handler.InitRoutes()

	// HTTPServer
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Application.Port),
		Handler: router,
	}

	return &Server{HTTPServer: httpServer}
}
