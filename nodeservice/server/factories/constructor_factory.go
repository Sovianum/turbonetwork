package factories

import (
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/core/graph"
	"fmt"
)

type ConstructorType func(data *pb.RequestData) (graph.Node, error)

type ConstructorFactory interface {
	GetConstructor(nodeType string) (ConstructorType, error)
}

func NewConstructorFactory() ConstructorFactory {
	return &constructorFactory{}
}

type constructorFactory struct {}

func (f *constructorFactory) GetConstructor(nodeType string) (ConstructorType, error) {
	if _, ok := constructorMap[nodeType]; !ok {
		return nil, fmt.Errorf("not found")
	}
	return constructorMap[nodeType], nil
}

var constructorMap = map[string]ConstructorType{
	PressureLossNodeType: pressureLossNodeConstructor,
}

func pressureLossNodeConstructor(data *pb.RequestData) (graph.Node, error) {
	return constructive.NewPressureLossNode(data.DKwargs["sigma"]), nil
}
