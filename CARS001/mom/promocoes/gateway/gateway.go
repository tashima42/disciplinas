// Package gateway
package gateway

import (
	"fmt"

	"github.com/tashima42/disciplinas/CARS001/mom/promocoes/crypto"
)

func Cadastrar() error {
	privateKey, err := crypto.ParsePrivateKey("./gateway/gateway.key")
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}
	publicKey, err := crypto.ParsePublicKey("./gateway/gateway.pub")
	if err != nil {
		return fmt.Errorf("failed to parse public key: %w", err)
	}

	message := "hello, world!"

	signature, err := crypto.Sign(privateKey, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to sign message: %w", err)
	}

	if err := crypto.Verify(publicKey, []byte(message), signature); err != nil {
		return fmt.Errorf("failed to verify message: %w", err)
	}

	return nil
}
