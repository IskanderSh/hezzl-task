package postgres

import (
	"fmt"
	"log/slog"

	"github.com/IskanderSh/hezzl-task/internal/config"
	"github.com/IskanderSh/hezzl-task/internal/lib/error/wrapper"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	log *slog.Logger
	db  *sqlx.DB
}

func NewStorage(log *slog.Logger, cfg config.Storage) (*Storage, error) {
	const op = "storage.postgres.NewStorage"

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password)
	log.Info(fmt.Sprintf("connection string for postgres: %s", connectionString))

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return &Storage{log: log, db: db}, nil
}
