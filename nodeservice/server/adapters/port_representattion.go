package adapters

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
)

type PortDescription = *pb.PortDescription
type PortExtractor func(node graph.Node) (graph.Port, error)

type PortRepresentation interface {
	ExtractPort(node graph.Node) (graph.Port, error)
	GetDescription() *PortDescription
}

