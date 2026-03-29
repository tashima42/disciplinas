// Package consumer
package consumer

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/rabbitmq"
)

type consumer struct {
	rq         *rabbitmq.RabbitMQ
	categories []string
	id         string
}

func NewConsumer(rabbitMqURL string, categories []string) (*consumer, error) {
	if len(categories) == 0 {
		return nil, errors.New("expected at least one category")
	}
	rq, err := rabbitmq.NewRabbitMQ(rabbitMqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	if err := rq.DeclareExchangePromocoes(); err != nil {
		return nil, err
	}

	id := uuid.New().String()

	return &consumer{rq: rq, id: id, categories: categories}, nil
}

func (c *consumer) Run() error {
	q, err := c.rq.Channel().QueueDeclare(
		"fila_cliente_"+c.id,
		true,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	for _, category := range c.categories {
		if err := c.rq.Channel().QueueBind(q.Name, "promocao."+category, "promocoes", false, nil); err != nil {
			return fmt.Errorf("failed to bind queue to promocoes exchange: %w", err)
		}
	}

	msgs, err := c.rq.Channel().Consume(
		q.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	for msg := range msgs {
		go func(msg amqp091.Delivery) {
			slog.Info(string(msg.Body))
		}(msg)
	}

	return nil
}
