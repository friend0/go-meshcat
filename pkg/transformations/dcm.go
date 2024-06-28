package transformations

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

func EulerToDCM(phi, theta, psi float64) (*mat.Dense, error) {
	c1 := math.Cos(phi)
	c2 := math.Cos(theta)
	c3 := math.Cos(psi)
	s1 := math.Sin(phi)
	s2 := math.Sin(theta)
	s3 := math.Sin(psi)
	dcm := []float64{
		c2 * c3, c2 * s3, -s2,
		c3*s2*s1 - s3*c1, s3*s2*s1 + c3*c1, c2 * s1,
		c2*s2*c1 + s3*s1, s3*s2*c1 - c3*s1, c2 * c1,
	}

	return mat.NewDense(3, 3, dcm), nil
}

func EulerToDCM4(phi, theta, psi float64) (*mat.Dense, error) {
	r3, err := EulerToDCM(phi, theta, psi)
	return reshapeMatrix(r3, 4, 4), err
}
