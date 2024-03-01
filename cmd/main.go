package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/IskanderSh/hezzl-task/internal/app"
	"github.com/IskanderSh/hezzl-task/internal/config"
)

const (
	envLocal = "local"
	envProd  = "prod"

	debugLvl = "DEBUG"
	infoLvl  = "INFO"
	warnLvl  = "WARN"
	errorLvl = "ERROR"
)

func main() {
	// load config file
	cfg := config.MustLoad()

	// init logger
	log := setupLogger(cfg)
	log.Info("logger initialized successfully")

	// init app
	application := app.NewServer(log, cfg)

	// start app
	go func() {
		if err := application.HTTPServer.ListenAndServe(); err != nil {
			panic(err)
		}
	}()
	log.Info("application started successfully")

	// graceful shutdown
}

func setupLogger(cfg *config.Config) *slog.Logger {
	var log *slog.Logger

	switch cfg.Env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: getLogLevel(cfg.LogLevel)}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: getLogLevel(cfg.LogLevel)}))
	}

	return log
}

func getLogLevel(lvl string) slog.Level {
	var res slog.Level

	switch strings.ToUpper(lvl) {
	case debugLvl:
		res = slog.LevelDebug
	case infoLvl:
		res = slog.LevelInfo
	case warnLvl:
		res = slog.LevelWarn
	case errorLvl:
		res = slog.LevelError
	}

	return res
}
