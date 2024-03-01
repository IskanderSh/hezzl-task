package clickhouse

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/IskanderSh/hezzl-task/internal/config"
	"github.com/IskanderSh/hezzl-task/internal/lib/error/wrapper"
)

type LogStorage struct {
	log        *slog.Logger
	connection *sql.DB
}

func NewLogStorage(log *slog.Logger, cfg config.LogStorage) (*LogStorage, error) {
	const op = "storage.clickhouse.NewLogStorage"

	connectionString := fmt.Sprintf("tcp://%s:%d", cfg.Host, cfg.Port)

	conn, err := sql.Open("clickhouse", connectionString)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, wrapper.Wrap(op, errors.New("couldn't ping clickhouse"))
	}

	return &LogStorage{
		log:        log,
		connection: conn,
	}, nil
}
