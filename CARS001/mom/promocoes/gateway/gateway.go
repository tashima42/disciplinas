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

type gateway struct {
	rq                *rabbitmq.RabbitMQ
	privateKey        *rsa.PrivateKey
	promocaoPublicKey *rsa.PublicKey
	scanner           *bufio.Scanner
	mu                *sync.Mutex
	promocoes         map[string]promocao.Promocao
}

func NewGateway(rabbitMqURL, gatewayPrivateKeyPath, promocaoPublicKeyPath string) (*gateway, error) {
	scanner := bufio.NewScanner(os.Stdin)

	rq, err := rabbitmq.NewRabbitMQ(rabbitMqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	if err := rq.DeclareExchangePromocoes(); err != nil {
		return nil, err
	}

	privateKey, err := crypto.ParsePrivateKey(gatewayPrivateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	promocaoPublicKey, err := crypto.ParsePublicKey(promocaoPublicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse promocao public key: %w", err)
	}

	return &gateway{
		rq:                rq,
		privateKey:        privateKey,
		promocaoPublicKey: promocaoPublicKey,
		scanner:           scanner,
		mu:                &sync.Mutex{},
		promocoes:         map[string]promocao.Promocao{},
	}, nil
}

func (g *gateway) Run() error {
	promocoesVerificadasQueue, err := g.rq.Channel().QueueDeclare("promocoes_verificadas_gateway", true, false, true, false, nil)
	if err != nil {
		return errors.New("failed to declare queue: " + err.Error())
	}

	if err := g.rq.Channel().QueueBind(promocoesVerificadasQueue.Name, "promocao.publicada", "promocoes", false, nil); err != nil {
		return errors.New("failed to bind queue to promocoes exchange: " + err.Error())
	}

	promocoesVerificadasChan, err := g.rq.Channel().Consume(
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

	go g.listenForVerifiedPromos(promocoesVerificadasChan)

	for {
		fmt.Print("\033[H\033[2J")

		fmt.Println("=== Gateway ===")
		fmt.Println("1 - Cadastrar promoção")
		fmt.Println("2 - Votar em uma promoção")
		fmt.Println("3 - Listar promoções verificadas")
		fmt.Println("4 - Sair")
		fmt.Print("\n> ")

		if !g.scanner.Scan() {
			break
		}
		input := strings.TrimSpace(g.scanner.Text())

		switch input {
		// cadastrar
		case "1":
			if err := g.cadastrarPromocaoMenu(); err != nil {
				fmt.Println("falha ao cadastrar promoção: " + err.Error())
			}
		case "2":
			if err := g.votarMenu(); err != nil {
				fmt.Println("falha ao votar na promoção: " + err.Error())
			}
		case "3":
			if err := g.listPromosMenu(); err != nil {
				fmt.Println("falha ao listar promocoes: " + err.Error())
			}
		case "4":
			fmt.Println("Tchau Tchau!")
			return nil
		default:
			fmt.Println("Opção inválida")
		}

		fmt.Println("Pressione [Enter] para continuar...")
		g.scanner.Scan()
	}

	return nil
}

func (g *gateway) listPromosMenu() error {
	fmt.Print("\033[H\033[2J")
	fmt.Println("=== Listar promoções verificadas ===")
	g.mu.Lock()
	for _, promo := range g.promocoes {
		fmt.Printf("%s - Titulo: %s, Categoria: %s\n", promo.ID, promo.Titulo, promo.Categoria)
	}
	g.mu.Unlock()
	return nil
}

func (g *gateway) listenForVerifiedPromos(promocoesVerificadasChan <-chan amqp091.Delivery) {
	for msg := range promocoesVerificadasChan {
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
			if err := crypto.Verify(g.promocaoPublicKey, msg.Body, s); err != nil {
				slog.Error("failed to verify message: " + err.Error())
				return
			}
			var promocao promocao.Promocao
			if err := json.Unmarshal(msg.Body, &promocao); err != nil {
				slog.Error("failed to unmarshal json: " + err.Error())
				return
			}

			g.mu.Lock()
			defer g.mu.Unlock()

			g.promocoes[promocao.ID] = promocao
		}(msg)
	}
}

func (g *gateway) cadastrarPromocaoMenu() error {
	fmt.Print("\033[H\033[2J")
	fmt.Println("=== Cadastrar promoção ===")

	fmt.Print("Título: ")
	if !g.scanner.Scan() {
		return errors.New("failed to scan")
	}
	titulo := strings.TrimSpace(g.scanner.Text())
	fmt.Print("Categoria: ")
	if !g.scanner.Scan() {
		return errors.New("failed to scan")
	}
	categoria := strings.TrimSpace(g.scanner.Text())
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = " cadastrando..."
	s.Start()
	if err := g.cadastrar(titulo, categoria); err != nil {
		return err
	}
	s.Stop()
	return nil
}

func (g *gateway) votarMenu() error {
	fmt.Print("\033[H\033[2J")
	fmt.Println("=== Votar em uma promoção ===")

	fmt.Print("id: ")
	if !g.scanner.Scan() {
		return errors.New("failed to scan")
	}
	id := strings.TrimSpace(g.scanner.Text())

	g.mu.Lock()
	defer g.mu.Unlock()

	promo, found := g.promocoes[id]
	if !found {
		return errors.New("promoção não encontrada, id: " + id)
	}

	return g.votar(promo)
}

func (g *gateway) cadastrar(titulo, categoria string) error {
	promo := promocao.NewPromocao(titulo, categoria)
	promoBody, err := json.Marshal(promo)
	if err != nil {
		return fmt.Errorf("failed to marshal promo to json: %w", err)
	}

	promoSignature, err := crypto.Sign(g.privateKey, promoBody)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	if err := g.rq.Channel().Publish("promocoes", "promocao.recebida", false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        promoBody,
		Headers:     amqp091.Table{"signature": promoSignature},
	}); err != nil {
		return fmt.Errorf("failed to publish message to exchange: %w", err)
	}

	return nil
}

func (g *gateway) votar(promo promocao.Promocao) error {
	promoBody, err := json.Marshal(promo)
	if err != nil {
		return fmt.Errorf("failed to marshal promo to json: %w", err)
	}

	promoSignature, err := crypto.Sign(g.privateKey, promoBody)
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	if err := g.rq.Channel().Publish("promocoes", "promocao.voto", false, false, amqp091.Publishing{
		ContentType: "application/json",
		Body:        promoBody,
		Headers:     amqp091.Table{"signature": promoSignature},
	}); err != nil {
		return fmt.Errorf("failed to publish message to exchange: %w", err)
	}
	return nil
}
