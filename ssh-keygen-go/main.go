package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

func generateSSHKeyPair() (privateKeyBytes []byte, publicKeyBytes []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, err
	}

	privateKeyBytes = encodePrivateKeyToPEM(privateKey)
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	publicKeyBytes = ssh.MarshalAuthorizedKey(publicKey)

	return privateKeyBytes, publicKeyBytes, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return pem.EncodeToMemory(block)
}

func main() {
	privateKeyBytes, publicKeyBytes, err := generateSSHKeyPair()
	if err != nil {
		log.Fatal("Error generating SSH key pair:", err)
	}

	fmt.Println("Private Key:")
	fmt.Println(string(privateKeyBytes))
	fmt.Println("Public Key:")
	fmt.Println(string(publicKeyBytes))
}
