package factories

import (
	"fmt"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbocycle/material/gases"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
	"github.com/Sovianum/turbocycle/common"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
)

type StateGetterType func(node *TypedNode, requiredFields []string) (*pb.NodeState, error)

type StateGetterFactory interface {
	GetStateGetter(nodeType string) (StateGetterType, error)
}

func NewStateGetterFactory() StateGetterFactory {
	return &stateGetterFactory{}
}

type stateGetterFactory struct{}

func (f *stateGetterFactory) GetStateGetter(nodeType string) (StateGetterType, error) {
	if _, ok := stateGetterMap[nodeType]; !ok {
		return nil, fmt.Errorf("not found")
	}
	return stateGetterMap[nodeType], nil
}

var stateGetterMap = map[string]StateGetterType{
	PressureLossNodeType: pressureLossPortStateGetter,
}

func pressureLossPortStateGetter(typedNode *TypedNode, requiredFields[]string) (*pb.NodeState, error) {
	if typedNode.NodeType != PressureLossNodeType {
		return nil, fmt.Errorf("failed to cast %s to %s", typedNode.NodeType, PressureLossNodeType)
	}
	result := newNodeState()

	node := typedNode.Node
	switch node.(type) {
	case constructive.PressureLossNode:
		casted := node.(constructive.PressureLossNode)
		updateFromComplexGasChannel(result, casted)
		return result, nil
	default:
		return nil, common.GetTypeError(PressureLossNodeType, typedNode.Node)
	}
}

func updateFromComplexGasChannel(state *pb.NodeState, c nodes.ComplexGasChannel) {
	updateFromTemperatureChannel(state, c)
	updateFromPressureChannel(state, c)
	updateFromGasChannel(state, c)
	updateFromMassRateChannel(state, c)
}

func updateFromGasChannel(state *pb.NodeState, c nodes.GasChannel) {
	updateFromGasSink(state, c)
	updateFromGasSource(state, c)
}

func updateFromGasSink(state *pb.NodeState, c nodes.GasSink) {
	s, _ := getGasPortState(c.GasInput())
	state.Children[gasOutput] = s
}

func updateFromGasSource(state *pb.NodeState, c nodes.GasSource) {
	s, _ := getGasPortState(c.GasOutput())
	state.Children[gasOutput] = s
}

func updateFromMassRateChannel(state *pb.NodeState, c nodes.MassRateChannel) {
	updateFromMassRateSink(state, c)
	updateFromMassRateSource(state, c)
}

func updateFromMassRateSink(state *pb.NodeState, c nodes.MassRateSink) {
	s, _ := getNumPortState(c.MassRateInput())
	state.Children[massRateOutput] = s
}

func updateFromMassRateSource(state *pb.NodeState, c nodes.MassRateSource) {
	s, _ := getNumPortState(c.MassRateOutput())
	state.Children[massRateOutput] = s
}

func updateFromPressureChannel(state *pb.NodeState, c nodes.PressureChannel) {
	updateFromPressureSink(state, c)
	updateFromPressureSource(state, c)
}

func updateFromPressureSink(state *pb.NodeState, c nodes.PressureSink) {
	s, _ := getNumPortState(c.PressureInput())
	state.Children[pressureOutput] = s
}

func updateFromPressureSource(state *pb.NodeState, c nodes.PressureSource) {
	s, _ := getNumPortState(c.PressureOutput())
	state.Children[pressureOutput] = s
}

func updateFromTemperatureChannel(state *pb.NodeState, c nodes.TemperatureChannel) {
	updateFromTemperatureSink(state, c)
	updateFromTemperatureSource(state, c)
}

func updateFromTemperatureSink(state *pb.NodeState, c nodes.TemperatureSink) {
	s, _ := getNumPortState(c.TemperatureInput())
	state.Children[temperatureOutput] = s
}

func updateFromTemperatureSource(state *pb.NodeState, c nodes.TemperatureSource) {
	s, _ := getNumPortState(c.TemperatureOutput())
	state.Children[temperatureOutput] = s
}

func getGasPortState(port graph.Port) (*pb.NodeState, error) {
	state := port.GetState()
	if state == nil {
		return &pb.NodeState{
			StringValues: map[string]string{
				"gas": "empty",
			},
		}, nil
	}

	val := state.Value()
	switch val.(type) {
	case gases.Gas:
		return &pb.NodeState{
			StringValues: map[string]string{
				"gas": val.(gases.Gas).String(),
			},
		}, nil
	default:
		return nil, fmt.Errorf("failed to cast state to gas")
	}
}

func getNumPortState(port graph.Port) (*pb.NodeState, error) {
	state := port.GetState()
	if state == nil {
		return &pb.NodeState{
			NumValues: map[string]float64{
				"value": 0,
			},
		}, nil
	}

	val := state.Value()
	switch val.(type) {
	case float64:
		return &pb.NodeState{
			NumValues: map[string]float64{
				"value": val.(float64),
			},
		}, nil
	default:
		return nil, fmt.Errorf("failed to cast state to float64")
	}
}

func newNodeState() *pb.NodeState {
	result := new(pb.NodeState)
	result.NumValues = make(map[string]float64)
	result.StringValues = make(map[string]string)
	result.Children = make(map[string]*pb.NodeState)
	return result
}
