package internal

import (
	"context"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

type Server struct {
	Router *echo.Echo
	NATS   *nats.Conn
	WS     *websocket.Conn
}

func NewServer(ctx context.Context) (*Server, error) {
	r := echo.New()
	nc, err := nats_connect()
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

func nats_connect() (*nats.Conn, error) {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 1 * time.Minute // Maximum total retry time
	var nc *nats.Conn
	err := backoff.Retry(
		func() (err error) {
			nc, err = nats.Connect(viper.GetString("NATS_URL"))
			if err != nil {
				return errors.Wrap(err, "failed to connect to NATS")
			}
			return nil
		}, expBackoff)
	if err != nil {
		return nc, errors.Wrap(err, "failed to connect to NATS")
	} else {
		return nc, nil
	}
}
