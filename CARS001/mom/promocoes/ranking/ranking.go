// Package ranking recebe e computa votos
package ranking

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/rabbitmq/amqp091-go"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/crypto"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/promocao"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/rabbitmq"
)

const hotDealThreshold = 2

type ranking struct {
	rq               *rabbitmq.RabbitMQ
	privateKey       *rsa.PrivateKey
	gatewayPublicKey *rsa.PublicKey
	promoVotes       map[string]*PromoVotes
	mu               *sync.Mutex
}

type PromoVotes struct {
	Promocao promocao.Promocao
	Votes    int
}

func NewRanking(rabbitMqURL, rankingPrivateKeyPath, gatewayPublicKeyPath string) (*ranking, error) {
	rq, err := rabbitmq.NewRabbitMQ(rabbitMqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	if err := rq.DeclareExchangePromocoes(); err != nil {
		return nil, err
	}

	privateKey, err := crypto.ParsePrivateKey(rankingPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	gatewayPublicKey, err := crypto.ParsePublicKey(gatewayPublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gateway public key: %w", err)
	}

	return &ranking{
		rq:               rq,
		privateKey:       privateKey,
		gatewayPublicKey: gatewayPublicKey,
		promoVotes:       map[string]*PromoVotes{},
		mu:               &sync.Mutex{},
	}, nil
}

func (r *ranking) Run() error {
	votosQueue, err := r.rq.Channel().QueueDeclare("votos", true, false, true, false, nil)
	if err != nil {
		return errors.New("failed to declare queue: " + err.Error())
	}

	if err := r.rq.Channel().QueueBind(votosQueue.Name, "promocao.voto", "promocoes", false, nil); err != nil {
		return errors.New("failed to bind queue to promocoes exchange: " + err.Error())
	}

	votosChan, err := r.rq.Channel().Consume(
		votosQueue.Name,
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

	for msg := range votosChan {
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
			if err := crypto.Verify(r.gatewayPublicKey, msg.Body, s); err != nil {
				slog.Error("failed to verify message: " + err.Error())
				return
			}
			var promocao promocao.Promocao
			if err := json.Unmarshal(msg.Body, &promocao); err != nil {
				slog.Error("failed to unmarshal json: " + err.Error())
				return
			}

			slog.Info(fmt.Sprintf("Voto: %+v\n", promocao))

			r.mu.Lock()
			defer r.mu.Unlock()

			if _, found := r.promoVotes[promocao.ID]; !found {
				r.promoVotes[promocao.ID] = &PromoVotes{promocao, 0}
			}

			r.promoVotes[promocao.ID].Votes += 1

			if err := r.calcHotDeals(promocao.ID); err != nil {
				slog.Error("failed to calc hot deal: " + err.Error())
				return
			}
		}(msg)
	}

	return nil
}

func (r *ranking) calcHotDeals(id string) error {
	pv, found := r.promoVotes[id]
	if !found {
		return errors.New("promo vote not found: " + id)
	}

	if pv.Votes < hotDealThreshold {
		return nil
	}

	promoBody, err := json.Marshal(pv.Promocao)
	if err != nil {
		return fmt.Errorf("failed to marshal promo to json: %w", err)
	}

	promoSignature, err := crypto.Sign(r.privateKey, promoBody)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	if err := r.rq.Channel().Publish("promocoes", "promocao.destaque", false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        promoBody,
		Headers:     amqp091.Table{"signature": promoSignature},
	}); err != nil {
		return fmt.Errorf("failed to publish message to exchange: %w", err)
	}

	slog.Info("Hot deal published: " + pv.Promocao.ID)

	return nil
}
