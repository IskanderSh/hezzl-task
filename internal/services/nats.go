package services

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/IskanderSh/hezzl-task/internal/config"
	"github.com/IskanderSh/hezzl-task/internal/lib/error/wrapper"
	"github.com/IskanderSh/hezzl-task/internal/models"
	"github.com/nats-io/nats.go"
)

type NatsServer struct {
	log        *slog.Logger
	connection *nats.Conn
	subject    string
	logStorage []models.GoodLog
}

func NewNatsServer(log *slog.Logger, cfg config.MessageBroker) (*NatsServer, error) {
	const op = "services.nats.NewNatsServer"

	connectString := fmt.Sprintf("nats://%s:%d", cfg.Host, cfg.Port)

	nc, err := nats.Connect(connectString)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return &NatsServer{log: log, connection: nc, subject: cfg.Subject, logStorage: make([]models.GoodLog, batchSize)}, nil
}

const (
	batchSize = 100
)

func (s *NatsServer) SendLog(log *models.GoodLog) {
	const op = "services.nats.SendLog"

	s.logStorage = append(s.logStorage, *log)

	if len(s.logStorage) == batchSize {
		err := s.publishLogs(&s.logStorage)
		if err != nil {
			s.log.Error("couldn't send logs to storage", err)
		}

		s.logStorage = make([]models.GoodLog, batchSize)
	}
}

func (s *NatsServer) publishLogs(logs *[]models.GoodLog) error {
	const op = "services.nats.PublishLogs"

	data, err := json.Marshal(logs)
	if err != nil {
		return wrapper.Wrap(op, err)
	}

	s.log.Debug("successfully marshal logs")

	if err := s.connection.Publish(s.subject, data); err != nil {
		return wrapper.Wrap(op, err)
	}

	s.log.Debug("successfully publish messages to subject")

	return nil
}
