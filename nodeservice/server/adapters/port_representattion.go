package adapters

import "github.com/Sovianum/turbocycle/core/graph"

type PortExtractor func(node graph.Node) (graph.Port, error)

type PortRepresentation interface {
	PortDescription
	ExtractPort(node graph.Node) (graph.Port, error)
	GetDescription() PortDescription
}

func NewSinglePortRepresentation(name string, extractor PortExtractor) PortRepresentation {
	return &portRepresentation{
		portDescription: newSinglePortDescription(name),
		extractor:       extractor,
	}
}

func NewMultiPortRepresentation(name string, extractor PortExtractor) PortRepresentation {
	return &portRepresentation{
		portDescription: newMultiPortDescription(name),
		extractor:       extractor,
	}
}

type portRepresentation struct {
	*portDescription
	extractor PortExtractor
}

func (r *portRepresentation) GetDescription() PortDescription {
	return r.portDescription
}

func (r *portRepresentation) ExtractPort(node graph.Node) (graph.Port, error) {
	return r.extractor(node)
}

