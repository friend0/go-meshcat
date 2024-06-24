package internal

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

type Server struct {
	Router *echo.Echo
	NATS   *nats.Conn
	Hub    *Hub
	Logger *slog.Logger
	Q      WorkQueue
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
	s.InitializeWorkQueue(10, 100, nc)
	s.Hub = NewHub()
	s.Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	s.Routes()
	s.NATSSubscriptions()
	return s, nil
}

func nats_connect() (*nats.Conn, error) {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 1 * time.Minute // Maximum total retry time
	var nc *nats.Conn
	nats_url, err := Getenv("NATS_URL", nats.DefaultURL)
	err = backoff.Retry(
		func() (err error) {
			nc, err = nats.Connect(nats_url)
			if err != nil {
				slog.Debug(fmt.Sprintf("failed to connect to NATS %v", err))
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
