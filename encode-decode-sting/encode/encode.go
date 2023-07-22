package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

var clave string = `1234567890`

func main() {

	// Texto para encriptar
	textoOriginal := "hola mundo...."

	// Convertir la clave a bytes y obtener el tamaño adecuado (16 bytes para AES-256)
	claveBytes := []byte(clave)[:32]

	// Encriptar el texto
	textoCifrado, err := encriptarAES(textoOriginal, claveBytes)
	if err != nil {
		fmt.Println("Error al encriptar:", err)
		return
	}

	fmt.Println("Texto cifrado:", textoCifrado)

	// Desencriptar el texto
	textoDesencriptado, err := desencriptarAES(textoCifrado, claveBytes)
	if err != nil {
		fmt.Println("Error al desencriptar:", err)
		return
	}

	fmt.Println("Texto desencriptado:", textoDesencriptado)
}

func encriptarAES(texto string, clave []byte) (string, error) {
	// Convertir el texto a bytes
	textoBytes := []byte(texto)

	// Añadir relleno PKCS7 al texto para que tenga un tamaño múltiplo del bloque AES
	blockSize := aes.BlockSize
	padding := blockSize - (len(textoBytes) % blockSize)
	padText := append(textoBytes, bytesRepeating(byte(padding), padding)...)

	// Crear un bloque AES
	block, err := aes.NewCipher(clave)
	if err != nil {
		return "", err
	}

	// Generar un IV (Vector de Inicialización) aleatorio
	iv := make([]byte, blockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Modo de cifrado CBC
	cifrador := cipher.NewCBCEncrypter(block, iv)

	// Crear un slice para el texto cifrado que incluya el IV
	textoCifrado := make([]byte, blockSize+len(padText))
	copy(textoCifrado[:blockSize], iv)

	// Encriptar el texto y almacenar el resultado en el slice de textoCifrado
	cifrador.CryptBlocks(textoCifrado[blockSize:], padText)

	// Convertir a base64 para que sea legible
	return base64.StdEncoding.EncodeToString(textoCifrado), nil
}

func desencriptarAES(textoCifrado string, clave []byte) (string, error) {
	// Decodificar el texto cifrado desde base64
	textoCifradoBytes, err := base64.StdEncoding.DecodeString(textoCifrado)
	if err != nil {
		return "", err
	}

	// Crear un bloque AES
	block, err := aes.NewCipher(clave)
	if err != nil {
		return "", err
	}

	// Extraer el IV (primeros 16 bytes del texto cifrado)
	iv := textoCifradoBytes[:aes.BlockSize]

	// Modo de cifrado CBC
	cifrador := cipher.NewCBCDecrypter(block, iv)

	// Desencriptar el texto y almacenar el resultado en un slice nuevo
	textoDesencriptado := make([]byte, len(textoCifradoBytes)-aes.BlockSize)
	cifrador.CryptBlocks(textoDesencriptado, textoCifradoBytes[aes.BlockSize:])

	// Eliminar el relleno PKCS7
	padding := int(textoDesencriptado[len(textoDesencriptado)-1])
	textoDesencriptado = textoDesencriptado[:len(textoDesencriptado)-padding]

	return string(textoDesencriptado), nil
}

// Función auxiliar para repetir un byte varias veces
func bytesRepeating(b byte, n int) []byte {
	bytes := make([]byte, n)
	for i := range bytes {
		bytes[i] = b
	}
	return bytes
}
