// Package gateway
package gateway

import (
	"bufio"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/briandowns/spinner"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/crypto"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/promocao"
	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/rabbitmq"
)

func Run(rabbitMqURL, gatewayPrivateKeyPath, promocaoPublicKeyPath string) error {
	scanner := bufio.NewScanner(os.Stdin)

	rq, err := rabbitmq.NewRabbitMQ(rabbitMqURL)
	if err != nil {
		return fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	if err := rq.DeclareExchangePromocoes(); err != nil {
		return err
	}

	privateKey, err := crypto.ParsePrivateKey(gatewayPrivateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	promocaoPublicKey, err := crypto.ParsePublicKey(promocaoPublicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to parse gateway public key: %w", err)
	}

	mu := sync.Mutex{}
	promocoes := []promocao.Promocao{}

	q, err := rq.Channel().QueueDeclare("promocoes_verificadas", true, false, true, false, nil)
	if err != nil {
		return errors.New("failed to declare queue: " + err.Error())
	}

	if err := rq.Channel().QueueBind(q.Name, "promocao.publicada", "promocoes", false, nil); err != nil {
		return errors.New("failed to bind promocoes_verificadas queue to promocoes exchange: " + err.Error())
	}

	msgs, err := rq.Channel().Consume(
		q.Name,
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

	go func(promocaoPublicKey *rsa.PublicKey, mu *sync.Mutex) {
		for msg := range msgs {
			go func(msg amqp091.Delivery, promocaoPublicKey *rsa.PublicKey, mu *sync.Mutex) {
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
				if err := crypto.Verify(promocaoPublicKey, msg.Body, s); err != nil {
					slog.Error("failed to verify message: " + err.Error())
					return
				}
				var promocao promocao.Promocao
				if err := json.Unmarshal(msg.Body, &promocao); err != nil {
					slog.Error("failed to unmarshal json: " + err.Error())
					return
				}

				mu.Lock()
				defer mu.Unlock()

				promocoes = append(promocoes, promocao)
			}(msg, promocaoPublicKey, mu)
		}
	}(promocaoPublicKey, &mu)

loop:
	for {
		fmt.Print("\033[H\033[2J")

		fmt.Println("=== Gateway ===")
		fmt.Println("1 - Cadastrar promoção")
		fmt.Println("2 - Votar em uma promoção")
		fmt.Println("3 - Listar promoções verificadas")
		fmt.Println("4 - Sair")
		fmt.Print("\n> ")

		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		switch input {
		// cadastrar
		case "1":
			fmt.Print("\033[H\033[2J")
			fmt.Println("=== Cadastrar promoção ===")

			fmt.Print("Título: ")
			if !scanner.Scan() {
				break loop
			}
			titulo := strings.TrimSpace(scanner.Text())
			fmt.Print("Categoria: ")
			if !scanner.Scan() {
				break loop
			}
			categoria := strings.TrimSpace(scanner.Text())
			s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
			s.Suffix = " cadastrando..."
			s.Start()
			if err := cadastrar(rq, privateKey, titulo, categoria); err != nil {
				return err
			}
			s.Stop()
		case "2":
			fmt.Print("\033[H\033[2J")
			fmt.Println("=== Votar em uma promoção ===")

			fmt.Print("id: ")
			if !scanner.Scan() {
				break loop
			}
			id := strings.TrimSpace(scanner.Text())
			if err := votar(rq, privateKey, id); err != nil {
				return err
			}
		case "3":
			fmt.Print("\033[H\033[2J")
			fmt.Println("=== Listar promoções verificadas ===")
			mu.Lock()
			for _, promo := range promocoes {
				fmt.Printf("%s - Titulo: %s, Categoria: %s\n", promo.ID, promo.Titulo, promo.Categoria)
			}
			mu.Unlock()
		case "4":
			fmt.Println("Tchau Tchau!")
			return nil
		default:
			fmt.Println("Opção inválida")
		}

		fmt.Println("Pressione [Enter] para continuar...")
		scanner.Scan()
	}

	return nil
}

func cadastrar(rq *rabbitmq.RabbitMQ, privateKey *rsa.PrivateKey, titulo, categoria string) error {
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

func votar(rq *rabbitmq.RabbitMQ, privateKey *rsa.PrivateKey, id string) error {
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
