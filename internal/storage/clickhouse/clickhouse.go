package clickhouse

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/IskanderSh/hezzl-task/internal/config"
	"github.com/IskanderSh/hezzl-task/internal/lib/error/wrapper"
	"github.com/IskanderSh/hezzl-task/internal/models"
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

func (s *LogStorage) NewLogs(logs *[]models.GoodLog) error {
	const op = "storage.clickhouse.NewLogs"

	log := s.log.With(slog.String("op", op))

	scope, err := s.connection.Begin()
	if err != nil {
		return wrapper.Wrap(op, err)
	}

	batch, err := scope.Prepare(insertQuery)
	if err != nil {
		return wrapper.Wrap(op, err)
	}

	for _, value := range *logs {
		_, err := batch.Exec(
			value.ID,
			value.ProjectID,
			value.Name,
			value.Description,
			value.Priority,
			value.Removed,
			value.EventTime,
		)
		if err != nil {
			log.Warn("error when inserting log to clickhouse", value.ID)
		}
	}

	err = scope.Commit()
	if err != nil {
		return wrapper.Wrap(op, err)
	}

	return nil
}
