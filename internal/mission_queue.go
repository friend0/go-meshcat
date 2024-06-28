package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	"github.com/friend0/go-meshcat/pkg/numgo"
	"github.com/nats-io/nats.go"
)

type Waypoint [][]float64

type Work interface {
	Do(results chan string)
}

type WorkQueue struct {
	Q       chan Work
	Results chan string
	NATS    *nats.Conn
}

func (s *Server) InitializeWorkQueue(workers int, queue_size int, conn *nats.Conn) {
	wq := WorkQueue{
		Q:       make(chan Work, queue_size),
		Results: make(chan string, queue_size),
		NATS:    conn,
	}
	for i := range workers {
		go MissionWorker(i, wq)
	}
	go wq.Gather()
	s.Q = wq
}

func (wq *WorkQueue) Add(work Work) error {
	select {
	case wq.Q <- work:
		return nil
	default:
		return errors.New("work channel full")
	}
}

func (wq *WorkQueue) Get() Work {
	return <-wq.Q
}

// todo: just clearing out results as we get missions completed
// Gather clears out the results queue.
// In the future, elements gathered on the results channel can be published to nats or logged
func (wq *WorkQueue) Gather() {
	<-wq.Results
}

func (wq *WorkQueue) Close() {
	close(wq.Q)
	close(wq.Results)
}

func MissionWorker(id int, wq WorkQueue) {
	for work := range wq.Q {
		fmt.Println("Worker", id, "started job")
		work.Do(wq.Results)
	}
}

type MissionWork struct {
	Conn      *nats.Conn
	Path      string
	Type      string
	Waypoints [][]float64
	Radius    float64
	Omega     float64
}

type NatsMissionWriter struct {
	Conn *nats.Conn
	Path string
}

func (nmw NatsMissionWriter) Write(p []byte) (n int, err error) {
	if nmw.Conn == nil {
		return 0, errors.New("NATS connection not initialized")
	}
	err = nmw.Conn.Publish(nmw.Path, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func transform_publisher(wp []float64, w io.Writer) error {
	// calculate the heading from the current waypoint
	heading := math.Atan2(-wp[1], wp[0])

	tc := SetTransformationCommand{
		Command: Command{
			Type: "set_transform",
			Path: "starling",
		},
		TransformationCommand: TransformationCommand{
			Translation: wp,
			Rotation:    []float64{0, 0, heading - 90},
		},
	}
	tranformation_json, err := json.Marshal(tc)
	if err != nil {
		return err
	}
	w.Write(tranformation_json)
	return nil
}

func WaypointIterator(sink io.Writer, waypoints [][]float64, transform_publisher func([]float64, io.Writer) error, ts time.Duration) {
	var wg sync.WaitGroup

	// todo: configurable minimum time step
	if ts < 8*time.Millisecond {
		ts = 8 * time.Millisecond
	}
	ticker := time.NewTicker(ts)
	wg.Add(1)
	go func(waypoints [][]float64) {
		defer wg.Done()
		for _, wp := range waypoints {
			<-ticker.C
			if len(wp) == 4 {
				ticker.Reset(time.Duration(wp[3]) * time.Millisecond)
			}
			transform_publisher(wp[:3], sink)
		}
	}(waypoints)
	wg.Wait()
	ticker.Stop()
}

func (mw MissionWork) Do(results chan string) {
	nmw := NatsMissionWriter{
		Path: mw.Path,
		Conn: mw.Conn,
	}
	if mw.Type == "orbit" {
		duration := 10
		fps := 60
		for {
			waypoints := numgo.Circspace(0, 2*math.Pi, mw.Radius, duration*fps*int(mw.Omega), true)
			WaypointIterator(nmw, waypoints, transform_publisher, 8*time.Millisecond)
		}
	}
	results <- "Complete"
}
