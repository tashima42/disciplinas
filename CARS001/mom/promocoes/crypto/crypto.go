// Package crypto generates, verifies and signs using rsa
package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

// GenerateKeyPair generates an rsa keypair and returns the private key, the public key or an error
func GenerateKeyPair() ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, errors.New("Error generating RSA key: " + err.Error())
	}

	der, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, errors.New("Error marshalling RSA private key: " + err.Error())
	}

	pubder, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, errors.New("Error marshalling RSA public key: " + err.Error())
	}

	private := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	})

	public := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubder,
	})

	return private, public, nil
}

func ParsePrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	privateKeyFile, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	privateBlock, _ := pem.Decode(privateKeyFile)
	privateKey, err := x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey.(*rsa.PrivateKey), nil
}

func ParsePublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}
	publicBlock, _ := pem.Decode(publicKeyFile)
	publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey.(*rsa.PublicKey), nil
}

func Sign(privateKey *rsa.PrivateKey, b []byte) ([]byte, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write(b)
	if err != nil {
		return nil, err
	}
	msgHashSum := msgHash.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, msgHashSum)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func Verify(publicKey *rsa.PublicKey, message []byte, signature []byte) error {
	msgHash := sha256.New()
	_, err := msgHash.Write(message)
	if err != nil {
		return err
	}
	msgHashSum := msgHash.Sum(nil)

	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, msgHashSum, signature)
}
