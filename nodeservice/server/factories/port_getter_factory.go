package factories

import (
	"github.com/Sovianum/turbocycle/core/graph"
	"fmt"
	"github.com/Sovianum/turbocycle/impl/engine/nodes"
	"github.com/Sovianum/turbocycle/impl/engine/nodes/constructive"
)

type PortGetter func(node *TypedNode, tag string) (graph.Port, error)

type PortGetterFactory interface {
	GetPortGetter(nodeType string) (PortGetter, error)
}

func NewPortGetterFactory() PortGetterFactory {
	return &portGetterFactory{}
}

type portGetterFactory struct {}

func (f *portGetterFactory) GetPortGetter(nodeType string) (PortGetter, error) {
	if getter, ok := portGetterMap[nodeType]; ok {
		return getter, nil
	}
	return nil, fmt.Errorf("getter for port %s not found", nodeType)
}

var portGetterMap = map[string]PortGetter{
	PressureLossNodeType: pressureLossPortGetter,
}

func pressureLossPortGetter(node *TypedNode, tag string) (graph.Port, error) {
	casted := node.Node.(constructive.PressureLossNode)

	if port, err := temperatureChannelPortGetter(casted, tag); err == nil {
		return port, nil
	} else if port, err := pressureChannelPortGetter(casted, tag); err == nil {
		return port, nil
	} else if port, err := gasChannelPortGetter(casted, tag); err == nil {
		return port, nil
	} else if port, err := massRateChannelPortGetter(casted, tag); err == nil {
		return port, nil
	}
	return nil, getNotFoundErr(tag)
}

func powerChannelPortGetter(s nodes.PowerChannel, tag string) (graph.Port, error) {
	if port, err := powerSourcePortGetter(s, tag); err == nil {
		return port, nil
	} else if port, err := powerSinkPortGetter(s, tag); err == nil {
		return port, nil
	}
	return nil, getNotFoundErr(tag)
}

func powerSinkPortGetter(s nodes.PowerSink, tag string) (graph.Port, error) {
	switch tag {
	case nodes.PowerInputTag:
		return s.PowerInput(), nil
	default:
		return nil, getNotFoundErr(tag)
	}
}

func powerSourcePortGetter(s nodes.PowerSource, tag string) (graph.Port, error) {
	switch tag {
	case nodes.PowerOutputTag:
		return s.PowerOutput(), nil
	default:
		return nil, getNotFoundErr(tag)
	}
}

func massRateChannelPortGetter(s nodes.MassRateChannel, tag string) (graph.Port, error) {
	if port, err := massRateSourcePortGetter(s, tag); err == nil {
		return port, nil
	} else if port, err := massRateSinkPortGetter(s, tag); err == nil {
		return port, nil
	}
	return nil, getNotFoundErr(tag)
}

func massRateSinkPortGetter(s nodes.MassRateSink, tag string) (graph.Port, error) {
	switch tag {
	case nodes.MassRateInputTag:
		return s.MassRateInput(), nil
	default:
		return nil, getNotFoundErr(tag)
	}
}

func massRateSourcePortGetter(s nodes.MassRateSource, tag string) (graph.Port, error) {
	switch tag {
	case nodes.MassRateOutputTag:
		return s.MassRateOutput(), nil
	default:
		return nil, getNotFoundErr(tag)
	}
}

func pressureChannelPortGetter(s nodes.PressureChannel, tag string) (graph.Port, error) {
	if port, err := pressureSourcePortGetter(s, tag); err == nil {
		return port, nil
	} else if port, err := pressureSinkPortGetter(s, tag); err == nil {
		return port, nil
	}
	return nil, getNotFoundErr(tag)
}

func pressureSinkPortGetter(s nodes.PressureSink, tag string) (graph.Port, error) {
	switch tag {
	case nodes.PressureInputTag:
		return s.PressureInput(), nil
	default:
		return nil, getNotFoundErr(tag)
	}
}

func pressureSourcePortGetter(s nodes.PressureSource, tag string) (graph.Port, error) {
	switch tag {
	case nodes.PressureOutputTag:
		return s.PressureOutput(), nil
	default:
		return nil, getNotFoundErr(tag)
	}
}

func temperatureChannelPortGetter(s nodes.TemperatureChannel, tag string) (graph.Port, error) {
	if port, err := temperatureSourcePortGetter(s, tag); err == nil {
		return port, nil
	} else if port, err := temperatureSinkPortGetter(s, tag); err == nil {
		return port, nil
	}
	return nil, getNotFoundErr(tag)
}

func temperatureSinkPortGetter(s nodes.TemperatureSink, tag string) (graph.Port, error) {
	switch tag {
	case nodes.TemperatureInputTag:
		return s.TemperatureInput(), nil
	default:
		return nil, getNotFoundErr(tag)
	}
}

func temperatureSourcePortGetter(s nodes.TemperatureSource, tag string) (graph.Port, error) {
	switch tag {
	case nodes.TemperatureOutputTag:
		return s.TemperatureOutput(), nil
	default:
		return nil, getNotFoundErr(tag)
	}
}

func gasChannelPortGetter(s nodes.GasChannel, tag string) (graph.Port, error) {
	if port, err := gasSourcePortGetter(s, tag); err == nil {
		return port, nil
	} else if port, err := gasSinkPortGetter(s, tag); err == nil {
		return port, nil
	}
	return nil, getNotFoundErr(tag)
}

func gasSinkPortGetter(s nodes.GasSink, tag string) (graph.Port, error) {
	switch tag {
	case nodes.GasInputTag:
		return s.GasInput(), nil
	default:
		return nil, getNotFoundErr(tag)
	}
}

func gasSourcePortGetter(s nodes.GasSource, tag string) (graph.Port, error) {
	switch tag {
	case nodes.GasOutputTag:
		return s.GasOutput(), nil
	default:
		return nil, getNotFoundErr(tag)
	}
}

func getNotFoundErr(tag string) error {
	return fmt.Errorf("failed to find port %s", tag)
}
