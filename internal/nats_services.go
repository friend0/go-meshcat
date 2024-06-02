package internal

import (
	"bytes"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
	"github.com/vmihailenco/msgpack/v5"
)

func (s *Server) NATSSubscriptions() error {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	// enc.UseArrayEncodedStructs(true)
	// enc.UseCompactFloats(true)
	// enc.UseCompactInts(true)

	// Subscribe to NATS subject
	_, err := s.NATS.Subscribe("meshcat", func(msg *nats.Msg) {
		log.Printf("Received meshcat message from NATS: %s", string(msg.Data))
		b := Box{}
		obj := b.hydrateObject()
		err := enc.Encode(&SetObject{
			Type:   "set_object",
			Object: obj,
			Path:   "/boxes/",
		})
		if err != nil {
			log.Printf("error sending msg: %v", err)
		}
		// Forward the message to the WebSocket server
		if s.WS != nil {
			err := s.WS.WriteMessage(websocket.BinaryMessage, buf.Bytes())
			if err != nil {
				log.Printf("Error sending message to WebSocket server: %v", err)
			}
		} else {
			fmt.Println("No WS conn")
		}
		buf.Reset()
	})
	if err != nil {
		log.Fatalf("Error subscribing to NATS subject: %v", err)
	}

	_, err = s.NATS.QueueSubscribe("meshcat.url", "MESHCAT_URL_Q", func(msg *nats.Msg) {
		b, err := msgpack.Marshal(&msg)
		fmt.Println("Here")
		if err != nil {
			log.Printf("error sending msg: %v", err)
		}
		if s.WS != nil {
			s.WS.WriteMessage(websocket.BinaryMessage, b)
			_, m, _ := s.WS.ReadMessage()
			log.Printf("read one %v", m)
		} else {
			log.Printf("no ws conn")
		}
	})
	if err != nil {
		log.Fatalf("Error subscribing to NATS subject: %v", err)
	}
	s.NATS.Flush()

	log.Printf("Listening on [%s]", "meshcat")
	return nil
}

var MeshcatCommands = map[string]bool{
	"set_transform": true,
	"set_object":    true,
	"delete":        true,
	"set_property":  true,
	"set_animation": true,
}

func (s *Server) MessageHandler(topic, msg string) error {
	if s.NATS == nil {

	}
	switch topic {
	case "url":
	case "wait":
	case "set_target":
	case "capture_image":
	default:
		_, ok := MeshcatCommands[topic]
		if !ok {
			return fmt.Errorf("%v is not a valid command", topic)
		}

	}
	return nil
}
