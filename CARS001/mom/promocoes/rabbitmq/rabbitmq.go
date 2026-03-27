// Package rabbitmq connects and handles connections to a rabbitmq broker
package rabbitmq

import (
	"errors"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQ(url string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, errors.New("failed to connecto to RabbitMQ: " + err.Error())
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.New("failed to open a channel: " + err.Error())
	}
	return &RabbitMQ{conn: conn, channel: ch}, nil
}

func (r *RabbitMQ) Channel() *amqp.Channel {
	return r.channel
}

func (r *RabbitMQ) DeclareExchangePromocoes() error {
	if err := r.Channel().ExchangeDeclare("promocoes", "topic", true, false, false, false, nil); err != nil {
		return fmt.Errorf("failed to declare exchange promocoes: %w", err)
	}
	return nil
}
