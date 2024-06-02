package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/friend0/go-meshcat/internal"
)

func run(ctx context.Context) (err error) {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	s, err := internal.NewServer()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Router.Shutdown(ctx)
	defer s.NATS.Close()
	if err := s.Router.Start(":8080"); err != http.ErrServerClosed {
		return err
	}
	runtime.Goexit()
	return nil
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := run(ctx); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
}
