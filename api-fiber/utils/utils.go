package utils

import (
	"api-fiber/config"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"

	"github.com/gofiber/fiber/v2"
)

// GenarateRSAKeys generates a pair of RSA private and public keys and returns them as strings in PEM format.
func GenarateRSAKeys() (string, string) {
	// Generate RSA private key with 2014 bits
	privateKey, err := rsa.GenerateKey(rand.Reader, 2014)
	if err != nil {
		log.Print(err)
	}

	// Marshal the private key to DER format
	privateKeyDer := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDer,
	}
	// Encode the private key to PEM format
	privateKeyPem := string(pem.EncodeToMemory(&privateKeyBlock))

	// Extract the public key from the private key
	publicKey := privateKey.PublicKey
	publicKeyDer, err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		log.Print(err)
	}

	// Create a PEM block for the public key
	publicKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}
	// Encode the public key to PEM format
	publicKeyPem := string(pem.EncodeToMemory(&publicKeyBlock))

	// Return the generated private and public keys in PEM format
	return privateKeyPem, publicKeyPem
}

// StartServer func for starting a simple server.
func StartServer(a *fiber.App) {
	// Run server.
	if err := a.Listen(config.Config("SERVER_URL")); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
