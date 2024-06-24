package internal

import (
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestLinspace(t *testing.T) {
	tests := []struct {
		start, end float64
		num        int
		expected   []float64
	}{
		{0, 1, 5, []float64{0, 0.25, 0.5, 0.75, 1}},
		{0, 10, 5, []float64{0, 2.5, 5, 7.5, 10}},
		{5, 5, 1, []float64{5}},
		{5, 5, 0, []float64{}},
		{0, 1, 1, []float64{0}},
		{1, 4, 4, []float64{1, 2, 3, 4}},
		{2, 2, 3, []float64{2, 2, 2}},
	}

	for _, test := range tests {
		got := Linspace(test.start, test.end, test.num)
		if !reflect.DeepEqual(got, test.expected) {
			t.Errorf("Linspace(%v, %v, %v) = %v; want %v", test.start, test.end, test.num, got, test.expected)
		}
	}
}

func TestWorkQueue(t *testing.T) {
  s := Server{}
  s.InitializeWorkQueue(2, 10, nil)
	s.Q.Add(MissionWork{Type: "orbit", Waypoints: [][]float64{{0, 0}, {1, 1}}})
	s.Q.Add(MissionWork{Type: "orbit", Waypoints: [][]float64{{0, 0}, {1, 1}}})
	s.Q.Add(MissionWork{Type: "orbit", Waypoints: [][]float64{{0, 0}, {1, 1}}})
	for i := 0; i < 3; i++ {
		result := <-s.Q.Results
		fmt.Println(result)
	}
	s.Q.Close()
}

func mock_publisher(wp []float64, w io.Writer) error {

	_, err := fmt.Fprintf(w, "Waypoint: %v\n", wp)
	return err
}

func TestWaypointIterator(t *testing.T) {
	wp := Circspace(0, 2*math.Pi, 1, 10)
	WaypointIterator(os.Stderr, wp, mock_publisher, 1*time.Millisecond)
}
