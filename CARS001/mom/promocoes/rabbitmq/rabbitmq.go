// Package rabbitmq connects and handles connections to a rabbitmq broker
package rabbitmq

import (
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	queue   amqp.Queue
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
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, errors.New("failed to declare a queue: " + err.Error())
	}

	return &RabbitMQ{conn: conn, queue: q, channel: ch}, nil
}
