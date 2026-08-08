package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	_ "image/jpeg" // Registra el decodificador JPG nativo
	"image/png"    // Registra el decodificador y codificador PNG nativo
	"os"
)

func main() {
	// 1. Ruta de la imagen local para el test
	nombreArchivo := "test.jpg"

	// 2. Abrir el archivo
	file, err := os.Open(nombreArchivo)
	if err != nil {
		fmt.Printf("Error: No se pudo abrir '%s'. Asegúrate de que existe en la carpeta.\n Error: %v\n", nombreArchivo, err)
		return
	}
	defer file.Close()

	// 3. Decodificar la imagen (detecta si es JPG o PNG automáticamente)
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Error al decodificar la imagen: %v\n", err)
		return
	}

	// 4. Codificar la imagen a PNG en un buffer de memoria
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		fmt.Printf("Error al codificar a PNG: %v\n", err)
		return
	}

	// 5. Convertir los bytes del PNG a formato Base64
	imagenBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	// 6. Enviar la secuencia de escape gráfica a la terminal
	// \x1b_G inicia el protocolo de gráficos de Kitty
	// a=T (transmitir y mostrar), f=100 (formato PNG)
	fmt.Print("\x1b_Ga=T,f=100;")
	fmt.Print(imagenBase64)
	fmt.Print("\x1b\\\n") // \x1b\\ cierra la secuencia de comandos

	fmt.Println("\n[Test completado] Si ves la imagen real arriba, tu terminal soporta gráficos nativos.")
}
