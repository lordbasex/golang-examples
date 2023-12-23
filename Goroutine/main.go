package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

var wg sync.WaitGroup
var stopFlag bool
var server *http.Server
var serverRunning bool
var stopCh chan struct{}

func main() {
	// Inicia la goroutine controlada
	wg.Add(1)
	stopCh = make(chan struct{})
	go printHelloWorld()

	// Espera a que la goroutine controlada termine
	// wg.Wait()
	go stopGoroutine()

	// Ejecuta varias goroutines sin control en un bucle infinito
	for i := 1; i <= 3; i++ {
		go printMessage(fmt.Sprintf("Mensaje %d", i))
	}

	// Espera para mantener las goroutines en ejecución
	select {}
}

func printHelloWorld() {
	defer wg.Done()

	// Configura el manejador HTTP
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("index").Parse("<html><body><h1>Hola Mundo</h1></body></html>")
		if err != nil {
			log.Println("Error al analizar la plantilla:", err)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println("Error al ejecutar la plantilla:", err)
		}
	})

	// Configura el servidor web en el puerto 8080
	server = &http.Server{Addr: ":8080"}

	// Canal para recibir el resultado de ListenAndServe
	done := make(chan error, 1)

	// Inicia el servidor web en una goroutine
	go func() {
		log.Println("Servidor web iniciado en http://localhost:8080")
		done <- server.ListenAndServe()
	}()

	// Espera hasta que stopFlag sea verdadero o se reciba una señal en stopCh
	select {
	case <-stopCh:
		// Cerrando el servidor web...
		log.Println("Cerrando el servidor web...")
		err := server.Shutdown(nil)
		if err != nil {
			log.Fatal("Error al cerrar el servidor web:", err)
		}
		serverRunning = false
	}

	// Espera a que ListenAndServe retorne
	err := <-done
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("Error al iniciar el servidor web:", err)
	}
}

func printMessage(message string) {
	for {
		log.Println(message)
		time.Sleep(1 * time.Second)
	}
}

func stopGoroutine() {
	// Espera 10 segundos antes de finalizar la goroutine
	time.Sleep(10 * time.Second)
	log.Println("stopGoroutine")
	close(stopCh) // Cierra el canal para enviar la señal de cierre
}
