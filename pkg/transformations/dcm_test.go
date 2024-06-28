package transformations

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/mat"
)

// isClose checks if two floating-point numbers are close to each other within a specified tolerance
func isClose(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}

func matricesAreClose(a, b *mat.Dense, tol float64) bool {
	ra, ca := a.Dims()
	rb, cb := b.Dims()
	if ra != rb || ca != cb {
		return false
	}
	for i := 0; i < ra; i++ {
		for j := 0; j < ca; j++ {
			if !isClose(a.At(i, j), b.At(i, j), tol) {
				return false
			}
		}
	}
	return true
}

func TestEulerToDCM(t *testing.T) {
	// Define test cases
	testCases := []struct {
		name     string
		input    []float64
		expected []float64
	}{
		{
			name:     "Zero angles",
			input:    []float64{0, 0, 0},
			expected: []float64{1, 0, 0, 0, 1, 0, 0, 0, 1},
		},
		{
			name:     "90 degrees rotation about X",
			input:    []float64{math.Pi / 2, 0, 0},
			expected: []float64{1, 0, 0, 0, 0, 1, 0, -1, 0},
		},
		// Add more test cases as needed
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the function with the input and compare the result with the expected output
			result, err := EulerToDCM(tc.input[0], tc.input[1], tc.input[2])
			if err != nil {
				t.Errorf("EulerToDCM(%v, %v, %v) returned an error: %v", tc.input[0], tc.input[1], tc.input[2], err)
			}
			expected := mat.NewDense(3, 3, tc.expected)
			if !matricesAreClose(result, expected, 1e-9) {
				t.Errorf("EulerToDCM(%v, %v, %v) = %v; want %v", tc.input[0], tc.input[1], tc.input[2], result, tc.expected)
			}
		})
	}
}
