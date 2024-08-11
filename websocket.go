package main

import (
	"log"
	app "lrucache/internal"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Server struct {
	clients      map[*WebSocketClient]bool
	cacheService app.Cache
}

func newWebsocketServer(cacheService app.Cache) *Server {
	server := Server{
		make(map[*WebSocketClient]bool),
		cacheService,
	}

	return &server
}

type WebSocketClient struct {
	conn  *websocket.Conn
	mutex sync.Mutex
}

func (server *Server) websocketHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("websocket connection")
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("websocket connection established")
	defer connection.Close()

	client := &WebSocketClient{conn: connection}
	server.clients[client] = true

	go func() {
		for range app.BroadcastChannel {
			cache := server.cacheService.GetCacheState()
			for client := range server.clients {
				client.mutex.Lock() // Lock before writing
				if err := client.conn.WriteJSON(cache); err != nil {
					log.Printf("WebSocket error: %v", err)
					client.conn.Close()
					delete(server.clients, client)
				}
				client.mutex.Unlock() // Unlock after writing
			}
		}
	}()

	// Block until the connection is closed
	for {
		_, _, err := connection.ReadMessage()
		if err != nil {
			log.Println("WebSocket connection closed:", err)
			break
		}
	}

}
