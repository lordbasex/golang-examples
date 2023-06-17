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

type ClientQueue struct {
	queues     map[string]map[*websocket.Conn]bool
	queueLock  sync.RWMutex
	senderMap  map[string]*senderState
	senderLock sync.Mutex
}

type senderState struct {
	done chan struct{}
	once sync.Once
}

func (cq *ClientQueue) AddClient(queueID string, client *Client) {
	cq.queueLock.Lock()
	defer cq.queueLock.Unlock()

	clients := cq.queues[queueID]
	if clients == nil {
		clients = make(map[*websocket.Conn]bool)
		cq.queues[queueID] = clients
	}
	clients[client.Socket] = true
}

func (cq *ClientQueue) RemoveClient(queueID string, client *Client) {
	cq.queueLock.Lock()
	defer cq.queueLock.Unlock()

	clients := cq.queues[queueID]
	if clients != nil {
		delete(clients, client.Socket)
		if len(clients) == 0 {
			cq.StopSender(queueID)
		}
	}
}

func (cq *ClientQueue) GetClients(queueID string) map[*websocket.Conn]bool {
	cq.queueLock.RLock()
	defer cq.queueLock.RUnlock()

	return cq.queues[queueID]
}

func (cq *ClientQueue) Broadcast(queueID string, message []byte) {
	clients := cq.GetClients(queueID)
	for client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println(err)
		}
	}
}

func (cq *ClientQueue) StartSender(queueID string) {
	cq.senderLock.Lock()
	defer cq.senderLock.Unlock()

	if cq.senderMap[queueID] != nil {
		return
	}

	state := &senderState{
		done: make(chan struct{}),
	}

	cq.senderMap[queueID] = state

	go func() {
		accelerator := 0
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if accelerator > 100 {
					accelerator = 0
				}
				log.Printf("Accelerator: %d", accelerator)

				cq.Broadcast(queueID, []byte(SetIntToString(accelerator)))

				accelerator++
			case <-state.done:
				return
			}
		}
	}()
}

func (cq *ClientQueue) StopSender(queueID string) {
	cq.senderLock.Lock()
	defer cq.senderLock.Unlock()

	state := cq.senderMap[queueID]
	if state == nil {
		return
	}

	state.once.Do(func() {
		close(state.done)
		delete(cq.senderMap, queueID)
	})
}

type Server struct {
	ClientQueue *ClientQueue
}

func SetIntToString(i int) string {
	return strconv.Itoa(i)
}

func reader(s *Server, queueID string, client *Client) {
	for {
		_, _, err := client.Socket.ReadMessage()
		if err != nil {
			log.Println(err)
			s.ClientQueue.RemoveClient(queueID, client)
			return
		}
	}
}

func wsEndpoint(s *Server, w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	queueID := r.URL.Query().Get("queueID")

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	var payload SocketPayload
	err = ws.ReadJSON(&payload)
	if err != nil {
		log.Println(err)
	}

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
	s.ClientQueue.AddClient(queueID, client)

	log.Println("Cliente conectado:", client.ID)

	err = ws.WriteMessage(websocket.TextMessage, []byte("pong"))
	if err != nil {
		log.Println(err)
	}

	if len(s.ClientQueue.GetClients(queueID)) == 1 {
		s.ClientQueue.StartSender(queueID)
	}

	reader(s, queueID, client)
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
		ClientQueue: &ClientQueue{
			queues:    make(map[string]map[*websocket.Conn]bool),
			senderMap: make(map[string]*senderState),
		},
	}
	log.Println("START")
	setupRoutes(server)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
