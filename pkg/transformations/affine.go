package transformations

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

func NewTransformation(translation, rotation []float64, matrix4 *[]float64) (transformation_matrix []float64, err error) {
	// check if Matrix4 is already specified, in which case, just return the result
	if matrix4 != nil && len(*matrix4) == 16 {
		allZero := true
		for _, v := range *matrix4 {
			if v != 0.0 {
				allZero = false
				break
			}
		}
		if !allZero {
			return transformation_matrix, nil
		}
	}

	// determine rotation type
	var rotation4 *mat.Dense
	if len(rotation) == 3 {
		rotation4, err = EulerToDCM4(rotation[0], rotation[1], rotation[2])
		if err != nil {
			return transformation_matrix, fmt.Errorf("unable to convert euler angles to DCM: %v", err)
		}
	} else if len(rotation) == 4 {
		// euler to quaternion
		rotation, err = EulerToQuaternion(([3]float64)(rotation))
		if err != nil {
			return transformation_matrix, fmt.Errorf("unable to convert euler angles to quaternion: %v", err)
		}
	} else {
		rotation4 = mat.NewDense(4, 4, []float64{1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 1})
	}
	if len(translation) == 3 {
		rotation4.Set(0, 3, translation[0])
		rotation4.Set(1, 3, translation[1])
		rotation4.Set(2, 3, translation[2])
		rotation4.Set(3, 3, 1)
	}
	rotation4 = mat.DenseCopyOf(rotation4.T())

	// todo: handle scaling matrix. Not a high priority for now
	return append(append(append(mat.Row(nil, 0, rotation4), mat.Row(nil, 1, rotation4)...), mat.Row(nil, 2, rotation4)...), mat.Row(nil, 3, rotation4)...), nil
}
