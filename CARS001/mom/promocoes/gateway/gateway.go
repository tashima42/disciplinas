// Package gateway
package gateway

import (
	"encoding/json"
	"fmt"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/crypto"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/promocao"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/rabbitmq"
)

func Cadastrar(rabbitMqURL, gatewayPrivateKeyPath, titulo, categoria string) error {
	privateKey, err := crypto.ParsePrivateKey(gatewayPrivateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	rq, err := rabbitmq.NewRabbitMQ(rabbitMqURL)
	if err != nil {
		return fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	if err := rq.DeclareExchangePromocoes(); err != nil {
		return err
	}

	promo := promocao.NewPromocao(titulo, categoria)
	promoBody, err := json.Marshal(promo)
	if err != nil {
		return fmt.Errorf("failed to marshal promo to json: %w", err)
	}

	promoSignature, err := crypto.Sign(privateKey, promoBody)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	if err := rq.Channel().Publish("promocoes", "promocao.recebida", false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        promoBody,
		Headers:     amqp091.Table{"signature": promoSignature},
	}); err != nil {
		return fmt.Errorf("failed to publish message to exchange: %w", err)
	}

	return nil
}

func Votar(rabbitMqURL, gatewayPrivateKeyPath, id string) error {
	privateKey, err := crypto.ParsePrivateKey(gatewayPrivateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	rq, err := rabbitmq.NewRabbitMQ(rabbitMqURL)
	if err != nil {
		return fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	if err := rq.DeclareExchangePromocoes(); err != nil {
		return err
	}

	promoSignature, err := crypto.Sign(privateKey, []byte(id))
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	if err := rq.Channel().Publish("promocoes", "promocao.voto", false, false, amqp091.Publishing{
		ContentType: "text/plain",
		Body:        []byte(id),
		Headers:     amqp091.Table{"signature": promoSignature},
	}); err != nil {
		return fmt.Errorf("failed to publish message to exchange: %w", err)
	}

	return nil
}
