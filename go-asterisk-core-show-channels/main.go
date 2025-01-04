package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ivahaev/amigo"
)

// struct CoreShowChannelEvent
type CoreShowChannelEvent struct {
	Channel         string `json:"Channel"`
	Context         string `json:"Context"`
	Exten           string `json:"Exten"`
	Priority        string `json:"Priority"`
	ChannelState    string `json:"ChannelState"`
	Application     string `json:"Application"`
	ApplicationData string `json:"ApplicationData"`
	CallerIDNum     string `json:"CallerIDNum"`
	AccountCode     string `json:"AccountCode"`
	PeerAccount     string `json:"PeerAccount"`
	Duration        string `json:"Duration"`
	BridgeId        string `json:"BridgeId"`
}

func main() {
	version := "1.0.0"
	fmt.Printf("go-asterisk-core-show-channels v%s Copyright (C) 2012 - 2025, CNSoluciones \n", version)
	fmt.Println("Created by Federico Pereira <fpereira@cnsoluciones.com>")
	fmt.Println()

	// Definir los flags para los argumentos de la línea de comandos
	host := flag.String("host", "", "Host for AMI (default from ENV or config)")
	port := flag.String("port", "", "Port for AMI (default from ENV or config)")
	username := flag.String("username", "", "Username for AMI (default from ENV or config)")
	password := flag.String("password", "", "Password for AMI (default from ENV or config)")

	// Parsear los argumentos
	flag.Parse()

	// Configuración de AMI con valores por defecto, o usando variables de entorno o argumentos
	settings := &amigo.Settings{
		Host:     getEnvOrDefault("AMI_HOST", *host, "127.0.0.1"),
		Port:     getEnvOrDefault("AMI_PORT", *port, "5038"),
		Username: getEnvOrDefault("AMI_USERNAME", *username, "admin"),
		Password: getEnvOrDefault("AMI_PASSWORD", *password, "password"),
	}

	// Inicialización y conexión
	a := amigo.New(settings)
	var mu sync.Mutex

	// Función para manejar la conexión al AMI
	connectToAMI(a)

	// Manejadores de eventos
	a.On("connect", func(message string) {
		fmt.Println("Connected:", message)

		// Crear un ticker para ejecutar CoreShowChannels cada 1 segundo
		ticker := time.NewTicker(1 * time.Second) // Ejecutar cada 1 segundo
		go func() {
			for range ticker.C {
				mu.Lock()
				// Enviar la acción "Command" con el comando "core show channels concise"
				result, err := a.Action(map[string]string{
					"Action":  "Command",
					"Command": "core show channels concise",
				})

				mu.Unlock()

				if err != nil {
					fmt.Println("Error sending CoreShowChannels:", err)
				} else {
					// Extraer el Output
					output := result["Output"]

					if output == "" {
						// Descarta si el Output está vacío
						continue
					}

					// Verificar si el output contiene el separador "!"
					if !strings.Contains(output, "!") {
						// Si no contiene "!", no hacer nada
						continue
					}

					// Separar la cadena por el separador "!"
					parts := strings.Split(output, "!")

					// Imprimir la cantidad de partes que se obtuvo
					// Solo procesar si tiene al menos 14 partes

					if len(parts) >= 14 {

						/*
							doc: https://github.com/asterisk/asterisk/blob/master/main/cli.c#L1151C43-L1152C69
						*/

						coreShowChannel := CoreShowChannelEvent{
							Channel:         parts[0],
							Context:         parts[1],
							Exten:           parts[2],
							Priority:        parts[3],
							ChannelState:    parts[4],
							Application:     parts[5],
							ApplicationData: parts[6],
							CallerIDNum:     parts[7],
							AccountCode:     parts[8],
							PeerAccount:     parts[9],
							Duration:        parts[11],
							BridgeId:        parts[13],
						}

						// Imprimir los datos en el formato solicitado
						fmt.Println("--------------->>")
						fmt.Println("Channel:", coreShowChannel.Channel)
						fmt.Println("    Context:", coreShowChannel.Context)
						fmt.Println("    Exten:", coreShowChannel.Exten)
						fmt.Println("    Priority:", coreShowChannel.Priority)
						fmt.Println("    ChannelState:", coreShowChannel.ChannelState)
						fmt.Println("    Application:", coreShowChannel.Application)
						fmt.Println("    ApplicationData:", coreShowChannel.ApplicationData)
						fmt.Println("    CallerIDNum:", coreShowChannel.CallerIDNum)
						fmt.Println("    AccountCode:", coreShowChannel.AccountCode)
						fmt.Println("    PeerAccount:", coreShowChannel.PeerAccount)
						fmt.Println("    Duration:", coreShowChannel.Duration)
						fmt.Println("    BridgeId:", coreShowChannel.BridgeId)
						fmt.Println()
						fmt.Println("    Dump:", output)
						fmt.Println("<<---------------")

					} else {
						// Si no hay suficientes partes, no hacer nada
						continue
					}
				}
			}
		}()
	})

	a.On("error", func(message string) {
		fmt.Println("Connection error:", message)
		reconnectAMI(a, &mu) // Intentar reconectar automáticamente
	})

	// Canal de eventos
	c := make(chan map[string]string, 100)
	a.SetEventChannel(c)

	// Mantener el programa en ejecución
	select {}
}

// Función para obtener valor de una variable de entorno o devolver el valor por defecto
func getEnvOrDefault(envVar, flagValue, defaultValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if e := os.Getenv(envVar); len(e) > 0 {
		return e
	}
	return defaultValue
}

// connectToAMI establece una conexión inicial al AMI
func connectToAMI(a *amigo.Amigo) {
	a.Connect() // Llamada sin capturar valor
	fmt.Println("Attempting to connect to AMI...")
}

// reconnectAMI intenta reconectar al AMI en caso de error
func reconnectAMI(a *amigo.Amigo, mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
	for {
		time.Sleep(5 * time.Second)
		fmt.Println("Attempting to reconnect to AMI...")
		a.Connect() // Llamada sin capturar valor
		if a.Connected() {
			fmt.Println("Reconnected successfully")
			break
		}
		fmt.Println("Reconnect failed, retrying...")
	}
}
