package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		log.Println(c.Locals("allowed"))
		log.Println(c.Params("id"))
		log.Println(c.Query("v"))
		log.Println(c.Cookies("session"))

		// Lee el mensaje inicial del cliente
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", msg)

		// Inicia un goroutine para enviar mensajes aleatorios cada 2 segundos
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			// Utiliza un bucle for range para esperar la señal del temporizador
			for range ticker.C {
				// Genera un número aleatorio y lo convierte en un mensaje
				randomNumber := rand.Intn(100)
				message := []byte(fmt.Sprintf("Número aleatorio: %d", randomNumber))
				// Envía el mensaje al cliente
				if err := c.WriteMessage(mt, message); err != nil {
					log.Println("write:", err)
					return
				}
			}
		}()

		// Continúa leyendo mensajes del cliente
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", msg)

			// Envía de vuelta el mismo mensaje al cliente
			if err = c.WriteMessage(mt, msg); err != nil {
				log.Println("write:", err)
				break
			}
		}
	}))

	log.Fatal(app.Listen(":3000"))
}
