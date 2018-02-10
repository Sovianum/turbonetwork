package configuration

import "fmt"

type GraphMultiMatrix interface {
	GetConnMatrix(selectors []int) (IntMatrix, error)
}

func NewGraphMultiMatrix(nodeMatrices []IntMatrix, nodeCnt int) (GraphMultiMatrix, error) {
	for i, nm := range nodeMatrices {
		_, c := nm.Dims()
		if c != nodeCnt {
			return nil, fmt.Errorf("columns[%d] != nodeCnt = %d", i, nodeCnt)
		}
	}
	return &graphMultiMatrix{
		nodeMatrices: nodeMatrices,
	}, nil
}

type graphMultiMatrix struct {
	nodeMatrices []IntMatrix
}

func (gmm *graphMultiMatrix) GetConnMatrix(selectors []int) (IntMatrix, error) {
	if len(selectors) != len(gmm.nodeMatrices) {
		return nil, fmt.Errorf("len of selectors does not match len of node matrices")
	}
	result := NewIntMatrix(len(gmm.nodeMatrices), len(gmm.nodeMatrices))

	for i, nm := range gmm.nodeMatrices {
		for j := 0; j != len(gmm.nodeMatrices); j++ {
			err := result.Set(nm.At(selectors[i], j), i, j)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}
