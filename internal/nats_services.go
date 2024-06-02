package internal

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

func (s *Server) NATSSubscriptions() {
	// Subscribe to NATS subject
	_, err := s.NATS.Subscribe("meshcat", func(msg *nats.Msg) {
		log.Printf("Received meshcat message from NATS: %s", string(msg.Data))

		// Forward the message to the WebSocket server
		if s.WS != nil {
			err := s.WS.WriteMessage(websocket.TextMessage, msg.Data)
			if err != nil {
				log.Printf("Error sending message to WebSocket server: %v", err)
			}
		} else {
			fmt.Println("No WS conn")
		}
	})
	s.NATS.Flush()
	if err != nil {
		log.Fatalf("Error subscribing to NATS subject: %v", err)
	}

	log.Printf("Listening on [%s]", "meshcat")
	log.SetFlags(log.LstdFlags)
}
