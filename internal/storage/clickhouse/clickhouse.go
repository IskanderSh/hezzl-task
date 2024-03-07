package clickhouse

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/IskanderSh/hezzl-task/internal/config"
	"github.com/IskanderSh/hezzl-task/internal/lib/error/wrapper"
	"github.com/IskanderSh/hezzl-task/internal/models"
)

type LogStorage struct {
	log        *slog.Logger
	connection driver.Conn
}

func NewLogStorage(ctx context.Context, log *slog.Logger, cfg config.LogStorage) (*LogStorage, error) {
	const op = "storage.clickhouse.NewLogStorage"

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	if err := conn.Ping(ctx); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Debug(fmt.Sprintf("Exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace))
		}
		return nil, wrapper.Wrap(op, err)
	}

	return &LogStorage{
		log:        log,
		connection: conn,
	}, nil
}

func (s *LogStorage) NewLogs(ctx context.Context, logs *[]models.GoodLog) error {
	const op = "storage.clickhouse.NewLogs"

	log := s.log.With(slog.String("op", op))

	batch, err := s.connection.PrepareBatch(ctx, insertQuery)
	if err != nil {
		return wrapper.Wrap(op, err)
	}

	for _, value := range *logs {
		err := batch.Append(
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

	err = batch.Send()
	if err != nil {
		return wrapper.Wrap(op, err)
	}

	return nil
}
