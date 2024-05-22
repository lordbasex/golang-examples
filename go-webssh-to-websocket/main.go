package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
)

type TerminalSize struct {
	Cols int `json:"cols"`
	Rows int `json:"rows"`
}

func sshHandler(c *websocket.Conn) {
	user := Config("SSH_USER")
	pass := Config("SSH_PASS")
	host := Config("SSH_HOST")
	port := Config("SSH_PORT")

	hostport := fmt.Sprintf("%s:%s", host, port)

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshConn, err := ssh.Dial("tcp", hostport, config)
	if err != nil {
		log.Println("SSH dial error:", err)
		return
	}
	defer sshConn.Close()

	session, err := sshConn.NewSession()
	if err != nil {
		log.Println("SSH session error:", err)
		return
	}
	defer session.Close()

	sshOut, err := session.StdoutPipe()
	if err != nil {
		log.Println("STDOUT pipe error:", err)
		return
	}

	sshIn, err := session.StdinPipe()
	if err != nil {
		log.Println("STDIN pipe error:", err)
		return
	}

	var termSize TerminalSize
	if err := c.ReadJSON(&termSize); err != nil {
		log.Println("Read terminal size error:", err)
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 8192,
		ssh.TTY_OP_OSPEED: 8192,
		ssh.IEXTEN:        0,
	}

	if err := session.RequestPty("xterm", termSize.Cols, termSize.Rows, modes); err != nil {
		log.Println("Request PTY error:", err)
		return
	}

	if err := session.Shell(); err != nil {
		log.Println("Start shell error:", err)
		return
	}

	go func() {
		defer session.Close()
		buf := make([]byte, 1024)
		for {
			n, err := sshOut.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Println("Read from SSH stdout error:", err)
				}
				return
			}
			if n > 0 {
				err = c.WriteMessage(websocket.BinaryMessage, buf[:n])
				if err != nil {
					log.Println("Write to WebSocket error:", err)
					return
				}
			}
		}
	}()

	for {
		messageType, p, err := c.ReadMessage()
		if err != nil {
			if err != io.EOF {
				log.Println("Read from WebSocket error:", err)
			}
			return
		}
		if messageType == websocket.TextMessage {
			var newSize TerminalSize
			if err := json.Unmarshal(p, &newSize); err == nil {
				session.WindowChange(newSize.Rows, newSize.Cols)
			} else {
				_, err = sshIn.Write(p)
				if err != nil {
					log.Println("Write to SSH stdin error:", err)
					return
				}
			}
		}
	}
}

// Config retrieves the value of an environment variable by its key.
func Config(key string) string {
	// Check if the .env file exists
	if _, err := os.Stat(".env"); err == nil {
		// If the file exists, load environment variables from it
		err := godotenv.Load(".env")
		if err != nil {
			// Print an error message if loading the .env file fails
			fmt.Println("Error loading .env file:", err)
		}
	}

	// Return the value of the specified environment variable
	return os.Getenv(key)
}

func checkForVariables() error {
	if Config("SSH_USER") == "" {
		return fmt.Errorf("SSH_USER is not set")
	}
	if Config("SSH_PASS") == "" {
		return fmt.Errorf("SSH_PASS is not set")
	}
	if Config("SSH_HOST") == "" {
		return fmt.Errorf("SSH_HOST is not set")
	}
	if Config("SSH_PORT") == "" {
		return fmt.Errorf("SSH_PORT is not set")
	}
	return nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	errOf := checkForVariables()
	if errOf != nil {
		log.Fatal(errOf)
	}

	app := fiber.New()

	// Logger middleware
	app.Use(logger.New())

	// Route for the root path
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./public/index.html")
	})

	// Middleware to verify WebSocket upgrade
	app.Use("/ssh", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ssh", websocket.New(func(c *websocket.Conn) {
		sshHandler(c)
	}))

	// Serve static files from the "public" directory
	app.Static("/", "./public")

	fmt.Println("Starting server on :8280")
	log.Fatal(app.Listen(":8280"))
}
