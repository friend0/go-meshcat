package transformations

import "gonum.org/v1/gonum/mat"

func reshapeMatrix(original *mat.Dense, newRows, newCols int) *mat.Dense {
	// Get the original dimensions
	rows, cols := original.Dims()

	// Create a new matrix with the desired dimensions
	newMatrix := mat.NewDense(newRows, newCols, nil)

	// Copy elements from the original matrix to the new matrix
	for i := 0; i < newRows; i++ {
		for j := 0; j < newCols; j++ {
			if i < rows && j < cols {
				newMatrix.Set(i, j, original.At(i, j))
			} else {
				// Fill extra cells with zero if the new matrix is larger
				newMatrix.Set(i, j, 0)
			}
		}
	}

	return newMatrix
}
