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
	// enc.UseArrayEncodedStructs(true)
	// enc.UseCompactFloats(true)
	// enc.UseCompactInts(true)

	// todo: track subscriptions so we can cleanup
	_, err := s.urlSubscription()
	if err != nil {
		return err
	}
	_, err = s.setObjectSubscription()
	if err != nil {
		return err
	}

	_, err = s.setTransform()
	if err != nil {
		return err
	}

	_, err = s.delete()
	if err != nil {
		return err
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

func (s *Server) urlSubscription() (*nats.Subscription, error) {

	sub, err := s.NATS.QueueSubscribe("meshcat.url", "MESHCAT_URL_Q", func(msg *nats.Msg) {
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
	return sub, err
}

// SetObject handler
func (s *Server) setObjectSubscription() (*nats.Subscription, error) {
	sub, err := s.NATS.Subscribe("meshcat.objects", func(msg *nats.Msg) {
		var buf bytes.Buffer
		enc := msgpack.NewEncoder(&buf)
		log.Printf("Received meshcat message from NATS: %s", string(msg.Data))
		b := NewBox(1, 1, 1)
		/* 		_, err := NewStarling(1.0, 1.0, 1.0)
		   		if err != nil {
		   			log.Printf("error during starling geometry creation")
		   		} */
		// s.Logger.Info(fmt.Sprintf("Received meshcat message from NATS: %d", len(b)))
		obj := Objectify(b)
		err := enc.Encode(SetObject{
			Object: obj,
			Command: Command{
				Type: "set_object",
				Path: "meshcat/objects",
			}})

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
	return sub, err
}

func (s *Server) setTransform() (*nats.Subscription, error) {
	return nil, nil
}

func (s *Server) delete() (*nats.Subscription, error) {
	return nil, nil
}
