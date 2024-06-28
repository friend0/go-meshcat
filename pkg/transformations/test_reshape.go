package transformations

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestReshapeMatrix(t *testing.T) {
	// Example original 3x3 matrix
	data := []float64{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	}
	original := mat.NewDense(3, 3, data)

	tests := []struct {
		newRows  int
		newCols  int
		expected *mat.Dense
	}{
		{
			newRows: 2, newCols: 2,
			expected: mat.NewDense(2, 2, []float64{
				1, 2,
				4, 5,
			}),
		},
		{
			newRows: 4, newCols: 4,
			expected: mat.NewDense(4, 4, []float64{
				1, 2, 3, 0,
				4, 5, 6, 0,
				7, 8, 9, 0,
				0, 0, 0, 0,
			}),
		},
		{
			newRows: 1, newCols: 9,
			expected: mat.NewDense(1, 9, []float64{
				1, 2, 3, 4, 5, 6, 7, 8, 9,
			}),
		},
	}

	for _, tt := range tests {
		result := reshapeMatrix(original, tt.newRows, tt.newCols)
		if !mat.EqualApprox(result, tt.expected, 1e-9) {
			t.Errorf("reshapeMatrix(%d, %d) = \n%v, want \n%v", tt.newRows, tt.newCols, mat.Formatted(result), mat.Formatted(tt.expected))
		}
	}
}
