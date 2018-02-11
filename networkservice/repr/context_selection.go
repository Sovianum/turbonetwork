package repr

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/networkservice/configuration"
	"github.com/Sovianum/turbonetwork/pb"
	"fmt"
)

func newContextSelector(nodes []RepresentationNode) *contextSelector {
	result := new(contextSelector)
	result.nodes = nodes
	result.portIndex = result.makePortIndex(nodes)
	result.connMatrix = configuration.NewIntMatrix(len(result.portIndex), len(result.portIndex))

	for i := 0; i != len(result.portIndex); i++ {
		for j := 0; j != len(result.portIndex); j++ {
			result.connMatrix.Set(
				int(pb.NodeDescription_AttachedPortDescription_NEUTRAL),
				i, j,
			)
		}
	}

	result.connConfigs = make([][]map[graph.Port]connType, len(nodes))

	for i, node := range nodes {
		result.connConfigs[i] = node.GetConnectionLines()
	}
	return result
}

type contextSelector struct {
	portIndex   map[graph.Port]int
	connMatrix  configuration.IntMatrix
	connConfigs [][]map[graph.Port]connType
	nodes       []RepresentationNode
}

func (cs *contextSelector) configure() error {
	validConfigurations := cs.findValidConfigurations()
	if l := len(validConfigurations); l == 0 {
		return fmt.Errorf("valid configs not found")
	} else if l > 1 {
		return fmt.Errorf("found multiple configs \n%v", validConfigurations)
	}
	config := validConfigurations[0]

	for i, node := range cs.nodes {
		if err := node.SelectState(config[i]); err != nil {
			return err
		}
	}
	return nil
}

func (cs *contextSelector) findValidConfigurations() [][]int {
	limits := make([]int, len(cs.connConfigs))
	for i, c := range cs.connConfigs {
		limits[i] = len(c)
	}

	variantSelectors := configuration.GetVariants(limits)
	var result [][]int
	for _, selector := range variantSelectors {
		variant := make([]map[graph.Port]connType, len(selector))
		for i, id := range selector {
			variant[i] = cs.connConfigs[i][id]
		}
		if cs.checkGraphMatrix(variant) {
			result = append(result, selector)
		}
	}
	return result
}

func (cs *contextSelector) checkGraphMatrix(variant []map[graph.Port]connType) bool {
	cs.initConnMatrix()
	for _, line := range variant {
		cs.updateConnMatrix(line)
	}
	return configuration.ValidateConnMatrix(cs.connMatrix, configuration.DefaultValidator)
}

func (cs *contextSelector) updateConnMatrix(connLine map[graph.Port]connType) {
	for from, connType := range connLine {
		to := from.GetLinkPort()
		cs.connMatrix.Set(
			int(connType), cs.portIndex[from], cs.portIndex[to],
		)
	}
}

func (cs *contextSelector) initConnMatrix() {
	for i := 0; i != len(cs.portIndex); i++ {
		for j := 0; j != len(cs.portIndex); j++ {
			cs.connMatrix.Set(int(pb.NodeDescription_AttachedPortDescription_NEUTRAL), i, j)
		}
	}
}

func (*contextSelector) makePortIndex(nodes []RepresentationNode) map[graph.Port]int {
	result := make(map[graph.Port]int)
	cnt := 0
	for _, node := range nodes {
		for _, port := range node.GetPorts() {
			if _, ok := result[port]; !ok {
				result[port] = cnt
				cnt++
			}
		}
	}
	return result
}
