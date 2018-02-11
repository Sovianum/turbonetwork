package repr

import (
	"github.com/Sovianum/turbonetwork/pb"
	"math"
)

func validateConnMatrix(cm *intMatrix, validatorFunc func(x1, x2 int) int) bool {
	rows, cols := cm.Dims()

	for i := 0; i != rows; i++ {
		for j := i + 1; j != cols; j++ {
			if validatorFunc(cm.At(i, j), cm.At(j, i)) != 0 {
				return false
			}
		}
	}
	return true
}

func defaultValidator(x1, x2 int) int {
	n := int(pb.NodeDescription_AttachedPortDescription_NEUTRAL)
	i := int(pb.NodeDescription_AttachedPortDescription_INPUT)
	o := int(pb.NodeDescription_AttachedPortDescription_OUTPUT)

	converter := func(x int) int {
		switch x {
		case n:
			return 0
		case i:
			return -1
		case o:
			return 1
		default:
			return math.MaxInt8
		}
	}

	if x1 == n || x2 == n {
		return 0
	}
	return converter(x1) + converter(x2)
}
