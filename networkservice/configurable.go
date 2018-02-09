package networkservice

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
)

type graphState struct {
	PortMap portMap
	NodeSet nodeSet
}

type portMap map[graph.Port]map[pb.NodeDescription_AttachedPortDescription_PortType]bool

type nodeSet map[graph.Node]bool
