package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

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
	// go wq.Gather()
	s.Q = wq
}

func (wq *WorkQueue) Add(work Work) {
	wq.Q <- work
}

func (wq *WorkQueue) Get() Work {
	return <-wq.Q
}

func (wq *WorkQueue) Close() {
	close(wq.Q)
	close(wq.Results)
}

func MissionWorker(id int, wq WorkQueue) {
	for {
		work := wq.Get()
		fmt.Println("Worker", id, "started job")
		work.Do(wq.Results)
	}
}

func (wq *WorkQueue) Gather() {
	for r := range wq.Results {
		fmt.Printf("Result: %s\n", r)
	}
}

type MissionWork struct {
	Path      string
	Type      string
	Radius    float64
	Omega     float64
	Waypoints [][]float64
	Conn      *nats.Conn
}

type Waypoint [][]float64

func Linspace(low, high float64, n int) []float64 {
	res := make([]float64, n)
	if n == 0 {
		return res
	}
	for i := 0; i < n; i++ {
		if i == 0 {
			res[i] = low
		} else {
			res[i] = low + float64(i)*(high-low)/float64(n-1)
		}
	}
	return res
}

func Circspace(low, high, radius float64, n int) [][]float64 {
	res := make([][]float64, n)
	if n == 0 {
		return res
	}
	for i := 0; i < n; i++ {
		t := low + float64(i)*(high-low)/float64(n)
		res[i] = []float64{radius * math.Cos(t), radius * math.Sin(t), 1}

	}
	return res
}

type NatsMissionWriter struct {
	Conn *nats.Conn
	Path string
}

func (nmw NatsMissionWriter) Write(p []byte) (n int, err error) {
	fmt.Printf("Publishing %v to %s", string(p), nmw.Path)
	err = nmw.Conn.Publish(nmw.Path, p)
	if err != nil {
		fmt.Println("Error publishing to NATS")
		fmt.Println(err)
		return 0, err
	}
	return len(p), nil
}

func transform_publisher(wp []float64, w io.Writer) error {
	tc := TransformationCommand{
		Translation: wp,
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
	ticker := time.NewTicker(ts)
	wg.Add(1)
	go func(waypoints [][]float64) {
		defer wg.Done()
		for _, wp := range waypoints {
			<-ticker.C
			fmt.Println("tick")
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
	// full_path := strings.Join([]string{"meshcat.transform"}, ".")
	nmw := NatsMissionWriter{
		Path: mw.Path,
		Conn: mw.Conn,
	}
	if mw.Type == "orbit" {
		waypoints := Circspace(0, 2*math.Pi, mw.Radius, 100)
		WaypointIterator(nmw, waypoints, transform_publisher, time.Duration(mw.Omega/100*1e9))
	}
	results <- "Complete"
}
