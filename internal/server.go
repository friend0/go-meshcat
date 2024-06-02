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
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		return nil, err
	}

	s := &Server{
		Router: r,
		NATS:   nc,
	}
	s.Routes()
	s.NATSSubscriptions()
	return s, nil
}
