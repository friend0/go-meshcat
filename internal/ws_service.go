package internal

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var writeWait = 10 * time.Second

const pongWait = 60 * time.Second

type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Messages to broadcast to clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
  hub := &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
  go hub.run()
  return hub
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
      fmt.Println("got a message for broadcast")
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) Write(message []byte) error {
  fmt.Println("Got here...")
  h.broadcast <- message
  fmt.Println("got here too")
  return nil
}

// Define the WebSocket upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 8192,
	WriteBufferPool: &sync.Pool{},
	CheckOrigin: func(r *http.Request) bool {
		return true
		// return origin == "http://localhost:8080"
	},
}

// func wsprocess(ws *websocket.Conn) {
// 	// defer ws.Close()
// 	for {
// 		// Read message from WebSocket
// 		messageType, message, err := ws.ReadMessage()
// 		if err != nil {
// 			log.Println("Read:", err, messageType)
// 			break
// 		}
// 		log.Printf("Received: %s", message)

// 		// Echo the message back to the client
// 		err = ws.WriteMessage(messageType, message)
// 		if err != nil {
// 			log.Println("Write:", err)
// 			break
// 		}
// 	}
// }

// func (s *Server) WsHandler() echo.HandlerFunc {
// 	// Define a handler for WebSocket connections
// 	return func(c echo.Context) error {
// 		// Upgrade the HTTP connection to a WebSocket connection
// 		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
// 		if err != nil {
// 			return err
// 		}
// 		ws.SetWriteDeadline(time.Now().Add(pongWait))
// 		ws.SetPongHandler(func(string) error { ws.SetWriteDeadline(time.Now().Add(pongWait)); fmt.Println("pongers"); return nil })
// 		ws.SetPingHandler(func(appData string) error {
// 			ws.SetWriteDeadline(time.Now().Add(writeWait))
// 			ws.WriteMessage(websocket.PingMessage, nil)
// 			return nil
// 		})
// 		ws.SetReadLimit(75 * 1024 * 1024)
// 		s.WS = ws
// 		go wsprocess(s.WS)
// 		return nil
// 	}
// }
