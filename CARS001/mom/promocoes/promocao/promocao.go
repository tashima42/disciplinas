// Package promocao encapsulates logic behind voting and publishing deals
package promocao

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/crypto"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/rabbitmq"
)

type Promocao struct {
	ID        string `json:"id"`
	Titulo    string `json:"titulo"`
	Categoria string `json:"categoria"`
}

func NewPromocao(titulo, categoria string) Promocao {
	return Promocao{ID: uuid.New().String(), Titulo: titulo, Categoria: categoria}
}

type verificador struct {
	rq               *rabbitmq.RabbitMQ
	privateKey       *rsa.PrivateKey
	gatewayPublicKey *rsa.PublicKey
}

func NewVerificador(rabbitMqURL, gatewayPublicKeyPath, promocaoPrivateKeyPath string) (*verificador, error) {
	privateKey, err := crypto.ParsePrivateKey(promocaoPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse promocao private key: %w", err)
	}
	gatewayPublicKey, err := crypto.ParsePublicKey(gatewayPublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gateway public key: %w", err)
	}
	rq, err := rabbitmq.NewRabbitMQ(rabbitMqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	if err := rq.DeclareExchangePromocoes(); err != nil {
		return nil, err
	}

	return &verificador{
		rq:               rq,
		privateKey:       privateKey,
		gatewayPublicKey: gatewayPublicKey,
	}, nil
}

func (v *verificador) Run() error {
	q, err := v.rq.Channel().QueueDeclare(
		"fila_promocao",
		true,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Esse serviço recebe eventos indicando que novas promoções foram recebidas.
	if err := v.rq.Channel().QueueBind(q.Name, "promocao.recebida", "promocoes", false, nil); err != nil {
		return fmt.Errorf("failed to bind fila_promocao queue to promocoes exchange: %w", err)
	}

	msgs, err := v.rq.Channel().Consume(
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

	// Ao receber um evento, o serviço deve inicialmente validar a assinatura digital da mensagem para garantir sua
	// autenticidade e integridade. Quando um evento de promoção é validado, o serviço registra
	// a promoção, assina e publica um novo evento informando que a promoção foi
	// disponibilizada no sistema.
	// O microsserviço Promocao consome os eventos promocao.recebida, assina digitalmente e publica eventos promocao.publicada.
	for msg := range msgs {
		go func(msg amqp091.Delivery) {
			signature, found := msg.Headers["signature"]
			if !found {
				slog.Error("signature header not found")
				return
			}
			s, ok := signature.([]byte)
			if !ok {
				slog.Error("failed to transform signature to bytes")
				return
			}
			if err := crypto.Verify(v.gatewayPublicKey, msg.Body, s); err != nil {
				slog.Error("failed to verify message: " + err.Error())
				return
			}
			var promocao Promocao
			if err := json.Unmarshal(msg.Body, &promocao); err != nil {
				slog.Error("failed to unmarshal json: " + err.Error())
				return
			}

			slog.Info(fmt.Sprintf("Promocao: %+v\n", promocao))

			promoSignature, err := crypto.Sign(v.privateKey, msg.Body)
			if err != nil {
				slog.Error("failed to sign message: " + err.Error())
				return
			}

			if err := v.rq.Channel().Publish("promocoes", "promocao.publicada", false, false, amqp091.Publishing{
				ContentType: "application/json",
				Body:        msg.Body,
				Headers:     amqp091.Table{"signature": promoSignature},
			}); err != nil {
				slog.Error("failed to publish message to exchange: " + err.Error())
				return
			}

			if err := msg.Ack(false); err != nil {
				slog.Error("failed to ack message: " + err.Error())
				return
			}
		}(msg)
	}
	return nil
}
