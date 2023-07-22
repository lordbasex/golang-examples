package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

var clave string = `1234567890` // Clave secreta para el cifrado (32 bytes para AES-256)

func main() {

	// Convertir la clave a bytes y obtener el tamaño adecuado (32 bytes para AES-256)
	claveBytes := []byte(clave)[:32]

	textoCifrado := "7oDlyivEKaTcvLmzOt52l1J7N1ldOsVNQ52fh4yarK4="

	// Desencriptar el texto
	textoDesencriptado, err := desencriptarAES(textoCifrado, claveBytes)
	if err != nil {
		fmt.Println("Error al desencriptar:", err)
		return
	}

	fmt.Println("Texto desencriptado:", textoDesencriptado)
}

// Función para desencriptar un texto cifrado con AES en modo CBC
func desencriptarAES(textoCifrado string, clave []byte) (string, error) {
	// Decodificar el texto cifrado desde base64
	textoCifradoBytes, err := base64.StdEncoding.DecodeString(textoCifrado)
	if err != nil {
		return "", err
	}

	// Crear un bloque AES con la clave proporcionada
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

	// Devolver el texto desencriptado como una cadena de texto
	return string(textoDesencriptado), nil
}
