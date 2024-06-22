// Package quaternions provides functions for common quaternion manipulations.
// Quaternions are a number system that extends the complex numbers.
// They are used in computer graphics, control theory, signal processing, and physics.
//
// This package includes functions for quaternion addition, subtraction, multiplication, division,
// conjugation, norm, and so on. Moreover, the primary convern of this module is the use of quaternions in specifying rotations
// of rigid bodies.

package transformations

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/mat"
)

type Quaternion []float64

// Dims, At and T minimally satisfy the mat.Matrix interface.
func (q Quaternion) Dims() (r, c int)    { return 1, len(q) }
func (q Quaternion) At(_, j int) float64 { return q[j] }
func (q Quaternion) T() mat.Matrix       { return QuaternionT(q) }

func (v Quaternion) RawVector() blas64.Vector {
	return blas64.Vector{N: len(v), Data: v, Inc: 1}
}

type QuaternionT []float64

// Dims, At and T minimally satisfy the mat.Matrix interface.
func (q QuaternionT) Dims() (r, c int)    { return 1, len(q) }
func (q QuaternionT) At(_, j int) float64 { return q[j] }
func (q QuaternionT) T() mat.Matrix       { return Quaternion(q) }

func (v QuaternionT) RawVector() blas64.Vector {
	return blas64.Vector{N: len(v), Data: v, Inc: 1}
}

func IsNormal(q Quaternion) bool {
	norm := math.Sqrt(q[0]*q[0] + q[1]*q[1] + q[2]*q[2] + q[3]*q[3])
	return math.Abs(1-norm) < 1e-9
}

// NewQuaternion creates a new quaternion from the given real and imaginary parts.
func NewQuaternionFromReIm(real float64, imag []float64) Quaternion {
	return Quaternion(append(imag, real))
}

func NewQuaternionFromSlice(q []float64) Quaternion {
	qv := Quaternion(q)
	if !IsNormal(qv) {
		panic("quaternion is not normal")
	}
	return qv
}

// quaternionToRotationMatrix converts a quaternion to a rotation matrix.
// The quaternion should be normalized, i.e., its norm (the square root of the sum of the squares of its components) should be 1.
// If the quaternion is not normalized, the function returns an error.
//
// The function takes a quaternion q represented as an array of 4 float64 values: [x, y, z, w].
// It returns a 3x3 rotation matrix represented as a *mat.Dense and an error.
//
// The rotation matrix is calculated using the formula:
// [
//
//	1 - 2*y*y - 2*z*z, 2*x*y - 2*z*w, 2*x*z + 2*y*w,
//	2*x*y + 2*z*w, 1 - 2*x*x - 2*z*z, 2*y*z - 2*x*w,
//	2*x*z - 2*y*w, 2*y*z + 2*x*w, 1 - 2*x*x - 2*y*y,
//
// ]
//
// Example usage:
//
//	q := [4]float64{0, 0, 0, 1}
//	matrix, err := quaternionToRotationMatrix(q)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(matrix)
func QuaternionToRotationMatrix(q []float64) (mat.Dense, error) {

	if !IsNormal(q) {
		return mat.Dense{}, fmt.Errorf("quaternion is not normalized")
	}
	x, y, z, w := q[0], q[1], q[2], q[3]

	// Compute the yaw (z), pitch (y), and roll (x) from the quaternion
	yaw := math.Atan2(2*(w*z+x*y), 1-2*(y*y+z*z))
	pitch := math.Asin(2 * (w*y - z*x))
	roll := math.Atan2(2*(w*x+y*z), 1-2*(x*x+y*y))

	// Compute the rotation matrix using the ZYX sequence
	cy, sy := math.Cos(yaw), math.Sin(yaw)
	cp, sp := math.Cos(pitch), math.Sin(pitch)
	cr, sr := math.Cos(roll), math.Sin(roll)

	return *mat.NewDense(3, 3, []float64{
		cy * cp, cy*sp*sr - sy*cr, cy*sp*cr + sy*sr,
		sy * cp, sy*sp*sr + cy*cr, sy*sp*cr - cy*sr,
		-sp, cp * sr, cp * cr,
	}), nil
}

// eulerToQuaternion converts Euler angles to a quaternion.
// The Euler angles are represented as an array of 3 float64 values: [roll, pitch, yaw].
// The function uses the aerospace sequence of rotations: ZYX, applied in order from right to left.
//
// The quaternion is calculated using the formula:
// q = [c1*c2*c3 + s1*s2*s3, s1*c2*c3 - c1*s2*s3, c1*s2*c3 + s1*c2*s3, c1*c2*s3 - s1*s2*c3]
// where c1 = cos(roll/2), s1 = sin(roll/2), c2 = cos(pitch/2), s2 = sin(pitch/2), c3 = cos(yaw/2), s3 = sin(yaw/2)
//
// Example usage:
//
//	e := [3]float64{0, 0, math.Pi/2}
//	q := eulerToQuaternion(e)
//	fmt.Println(q) // Outputs: [0 0 1 0]
func EulerToQuaternion(e [3]float64) (Quaternion, error) {
	c1, s1 := math.Cos(e[0]/2), math.Sin(e[0]/2)
	c2, s2 := math.Cos(e[1]/2), math.Sin(e[1]/2)
	c3, s3 := math.Cos(e[2]/2), math.Sin(e[2]/2)

	return Quaternion([]float64{
		c1*c2*c3 + s1*s2*s3,
		s1*c2*c3 - c1*s2*s3,
		c1*s2*c3 + s1*c2*s3,
		c1*c2*s3 - s1*s2*c3,
	}), nil
}
