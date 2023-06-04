package main

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Client struct {
	ID     string
	Socket *websocket.Conn
}

type SocketPayload struct {
	Headers map[string]string `json:"headers"`
	Content string            `json:"content"`
}

type Server struct {
	Clients map[string]*Client
	mutex   sync.Mutex
}

func (s *Server) AddClient(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Clients[client.ID] = client
}

func (s *Server) RemoveClient(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.Clients, client.ID)
}

// SetIntToString es una función que toma un entero como entrada y devuelve una cadena como resultado.
// Utiliza la función Itoa de la biblioteca de strconv para convertir el entero a una cadena.
func SetIntToString(i int) string {
	r := strconv.Itoa(i)
	return r
}

func reader(s *Server, client *Client) {
	for {
		messageType, p, err := client.Socket.ReadMessage()
		if err != nil {
			log.Println(err)
			s.RemoveClient(client)
			return
		}
		log.Println(string(p))

		if err := client.Socket.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			s.RemoveClient(client)
			return
		}
	}
}

func sender(s *Server, client *Client, onDone func()) {
	accelerator := 0
	ticker := time.NewTicker(1 * time.Second) // Intervalo de 1 segundo
	defer ticker.Stop()

	for range ticker.C {
		if accelerator > 100 {
			accelerator = 0
		}
		log.Printf("Accelerator: %d", accelerator)
		err := client.Socket.WriteMessage(websocket.TextMessage, []byte(SetIntToString(accelerator)))
		if err != nil {
			log.Println(err)
			s.RemoveClient(client)
			onDone()
			return
		}
		accelerator++
	}

	// Cuando el sender ha terminado su trabajo, llamamos a la función onDone
	//onDone()
}

func wsEndpoint(s *Server, w http.ResponseWriter, r *http.Request) {

	// Obtener el ID del cliente (puede ser generado o provisto por el cliente)
	userID := r.URL.Query().Get("userID")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// Leer y analizar el mensaje JSON
	var payload SocketPayload
	err = ws.ReadJSON(&payload)
	if err != nil {
		log.Println(err)
	}

	// Obtener los encabezados y el contenido del mensaje
	authHeader := payload.Headers["Authorization"]
	icid := payload.Headers["icid"]
	content := payload.Content

	log.Printf("authHeader: %s", authHeader)
	log.Printf("icid: %s", icid)
	log.Printf("content: %s", content)

	client := &Client{
		ID:     userID,
		Socket: ws,
	}
	s.AddClient(client)

	log.Println("Cliente conectado:", client.ID)

	err = ws.WriteMessage(1, []byte("pong"))
	if err != nil {
		log.Println(err)
	}

	go sender(s, client, func() {
		// Esta función se ejecutará cuando la goroutine de sender finalice
		log.Println("La goroutine de sender se ha destruido")
		log.Print(client)
	})

	reader(s, client)
}

func setupRoutes(s *Server) {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsEndpoint(s, w, r)
	})
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	server := &Server{
		Clients: make(map[string]*Client),
	}
	log.Println("START")
	setupRoutes(server)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
