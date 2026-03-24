package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("expected receive or send")
	}

	arg := os.Args[1]

	if arg != "send" && arg != "receive" {
		log.Fatal("expected receive or send")
	}

	conn, q, ch, err := start("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("failed to start RabbitMQ: " + err.Error())
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}
		if err := ch.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	switch arg {
	case "send":
		if err := send(q, ch); err != nil {
			log.Fatal("failed to publish a message: " + err.Error())
		}
	case "receive":
		if err := receive(q, ch); err != nil {
			log.Fatal("failed to receive a message: " + err.Error())
		}
	default:
		log.Fatal("unknown state")
	}
}

func start(url string) (*amqp.Connection, amqp.Queue, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, amqp.Queue{}, nil, errors.New("failed to connecto to RabbitMQ: " + err.Error())
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, nil, errors.New("failed to open a channel: " + err.Error())
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
		return nil, amqp.Queue{}, nil, errors.New("failed to declare a queue: " + err.Error())
	}
	return conn, q, ch, nil
}

func send(q amqp.Queue, ch *amqp.Channel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err := ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		return err
	}
	log.Printf(" [x] Sent %s\n", body)
	return nil
}

func receive(q amqp.Queue, ch *amqp.Channel) error {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return errors.New("failed to register a consumer: " + err.Error())
	}

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan // Blocks until you hit Ctrl+C or kill the process
	return nil
}
