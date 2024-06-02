package internal

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
)

type Server struct {
	Router *echo.Echo
	NATS   *nats.Conn
	WS     *websocket.Conn
}

func NewServer() (*Server, error) {
	r := echo.New()
	nc, err := nats.Connect("127.0.0.1:4222")
	if err != nil {
		return nil, err
	}
	nc.Publish("meshcat", []byte("TESTER"))
	// defer nc.Close()

	// ws_url := "ws://localhost"
	// // Connect to WebSocket server
	// wc, _, err := websocket.DefaultDialer.Dial(ws_url, nil)
	// if err != nil {
	// 	log.Fatalf("Error connecting to WebSocket server: %v", err)
	// }
	// defer wc.Close()

	s := &Server{
		Router: r,
		NATS:   nc,
		// WS:     wc,
	}
	s.Routes()
	go s.NATSSubscriptions()
	return s, nil
}
