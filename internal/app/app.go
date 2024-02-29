package app

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/IskanderSh/hezzl-task/internal/config"
	"github.com/IskanderSh/hezzl-task/internal/handlers"
	"github.com/IskanderSh/hezzl-task/internal/services"
	redis "github.com/IskanderSh/hezzl-task/internal/storage/cache"
	"github.com/IskanderSh/hezzl-task/internal/storage/postgres"
)

type Server struct {
	HTTPServer *http.Server
}

func NewServer(log *slog.Logger, cfg *config.Config) *Server {
	// Storages
	storage, err := postgres.NewStorage(cfg.Storage)
	if err != nil {
		panic(err)
	}

	cache := redis.NewCache(cfg.Cache)

	// Services
	goodService := services.NewGoodService(log, storage, cache)

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
