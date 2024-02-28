package postgres

import (
	"fmt"

	"github.com/IskanderSh/hezzl-task/internal/config"
	"github.com/IskanderSh/hezzl-task/internal/lib/error/wrapper"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage(cfg config.Storage) (*Storage, error) {
	const op = "storage.postgres.NewStorage"

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password)
	print(connectionString)

	db, err := sqlx.Open("postgres", connectionString)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	if err = db.Ping(); err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return &Storage{db: db}, nil
}
