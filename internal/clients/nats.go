package clients

import (
	"fmt"
	"log/slog"

	"github.com/IskanderSh/hezzl-task/internal/config"
	"github.com/IskanderSh/hezzl-task/internal/lib/error/wrapper"
	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	log        *slog.Logger
	connection *nats.Conn
	provider   LogsProvider
}

type LogsProvider interface {
}

func NewNatsClient(log *slog.Logger, cfg config.MessageBroker, provider LogsProvider) (*NatsClient, error) {
	connectString := fmt.Sprintf("nats://%s:%d", cfg.Host, cfg.Port)

	nc, err := nats.Connect(connectString)
	if err != nil {
		return nil, err
	}

	return &NatsClient{log: log, connection: nc, provider: provider}, nil
}

func (nc *NatsClient) SubscribeSubjects(subject string) error {
	const op = "clients.nats.SubscribeSubjects"

	_, err := nc.connection.Subscribe(subject, nc.ReceiveLog)
	if err != nil {
		return wrapper.Wrap(op, err)
	}

	return nil
}

func (nc *NatsClient) ReceiveLog(m *nats.Msg) {
	const op = "clients.nats.ReceiveLog"

	log := nc.log.With(slog.String("op", op))

	log.Debug(fmt.Sprintf("log's subject received a message: %s", string(m.Data)))
}
