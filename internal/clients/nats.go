package clients

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/IskanderSh/hezzl-task/internal/config"
	"github.com/IskanderSh/hezzl-task/internal/lib/error/wrapper"
	"github.com/IskanderSh/hezzl-task/internal/models"
	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	log        *slog.Logger
	connection *nats.Conn
	subject    string
	provider   LogsProvider
}

type LogsProvider interface {
	NewLogs(logs *[]models.GoodLog) error
}

func NewNatsClient(log *slog.Logger, cfg config.MessageBroker, provider LogsProvider) (*NatsClient, error) {
	const op = "clients.nats.NewNatsClient"

	connectString := fmt.Sprintf("nats://%s:%d", cfg.Host, cfg.Port)

	nc, err := nats.Connect(connectString)
	if err != nil {
		return nil, wrapper.Wrap(op, err)
	}

	return &NatsClient{log: log, connection: nc, subject: cfg.Subject, provider: provider}, nil
}

func (nc *NatsClient) SubscribeSubjects() error {
	const op = "clients.nats.SubscribeSubjects"

	_, err := nc.connection.Subscribe(nc.subject, nc.ReceiveLog)
	if err != nil {
		return wrapper.Wrap(op, err)
	}

	return nil
}

func (nc *NatsClient) ReceiveLog(m *nats.Msg) {
	const op = "clients.nats.ReceiveLog"

	log := nc.log.With(slog.String("op", op))

	log.Debug(fmt.Sprintf("log's subject received a message: %s", string(m.Data)))

	var logs []models.GoodLog

	if err := json.Unmarshal(m.Data, &logs); err != nil {
		log.Error("couldn't convert data to logs struct")
	}

	if err := nc.provider.NewLogs(&logs); err != nil {
		log.Error("couldn't save logs in logs storage", err)
	}
}
