// Package notificacao
package notificacao

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/crypto"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/promocao"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/rabbitmq"
)

type notificacao struct {
	rq                *rabbitmq.RabbitMQ
	promocaoPublicKey *rsa.PublicKey
	rankingPublicKey  *rsa.PublicKey
}

func NewNotificacao(rabbitMqURL, promocaoPublicKeyPath, rankingPublicKeyPath string) (*notificacao, error) {
	rq, err := rabbitmq.NewRabbitMQ(rabbitMqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	if err := rq.DeclareExchangePromocoes(); err != nil {
		return nil, err
	}

	promocaoPublicKey, err := crypto.ParsePublicKey(promocaoPublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse promocao public key: %w", err)
	}

	rankingPublicKey, err := crypto.ParsePublicKey(rankingPublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ranking public key: %w", err)
	}

	return &notificacao{
		rq:                rq,
		promocaoPublicKey: promocaoPublicKey,
		rankingPublicKey:  rankingPublicKey,
	}, nil
}

func (n *notificacao) Run() error {
	promocoesVerificadasQueue, err := n.rq.Channel().QueueDeclare("promocoes_verificadas_notificacao", true, false, true, false, nil)
	if err != nil {
		return errors.New("failed to declare queue: " + err.Error())
	}

	if err := n.rq.Channel().QueueBind(promocoesVerificadasQueue.Name, "promocao.publicada", "promocoes", false, nil); err != nil {
		return errors.New("failed to bind promocoes_verificadas queue to promocoes exchange: " + err.Error())
	}

	promocoesVerificadasChan, err := n.rq.Channel().Consume(
		promocoesVerificadasQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.New("failed to register consumer: " + err.Error())
	}

	promocoesDestaqueQueue, err := n.rq.Channel().QueueDeclare("promocoes_destaque", true, false, true, false, nil)
	if err != nil {
		return errors.New("failed to declare queue: " + err.Error())
	}

	if err := n.rq.Channel().QueueBind(promocoesDestaqueQueue.Name, "promocao.destaque", "promocoes", false, nil); err != nil {
		return errors.New("failed to bind promocoes_destaque queue to promocoes exchange: " + err.Error())
	}

	promocoesDestaqueChan, err := n.rq.Channel().Consume(
		promocoesDestaqueQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return errors.New("failed to register consumer: " + err.Error())
	}

	go n.notifyPromoCategory(promocoesVerificadasChan, "promocao", n.promocaoPublicKey)
	go n.notifyPromoCategory(promocoesDestaqueChan, "hot deal!", n.rankingPublicKey)

	// keep running indefinitely
	select {}
}

func (n *notificacao) notifyPromoCategory(promosChan <-chan amqp091.Delivery, prefix string, publicKey *rsa.PublicKey) {
	for msg := range promosChan {
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
			if err := crypto.Verify(publicKey, msg.Body, s); err != nil {
				slog.Error("failed to verify message: " + err.Error())
				return
			}
			var promocao promocao.Promocao
			if err := json.Unmarshal(msg.Body, &promocao); err != nil {
				slog.Error("failed to unmarshal json: " + err.Error())
				return
			}

			if err := n.rq.Channel().Publish("promocoes", "promocao."+promocao.Categoria, false, false, amqp091.Publishing{
				ContentType: "text/plain",
				Body:        fmt.Appendf(nil, "%s - [%s]: %s", prefix, promocao.ID, promocao.Titulo),
			}); err != nil {
				slog.Error("failed to publish message to exchange: " + err.Error())
			}
		}(msg)
	}
}
