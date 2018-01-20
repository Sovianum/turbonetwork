package factories

import (
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"fmt"
	"github.com/Sovianum/turbocycle/core/graph"
)

type UpdaterType func(node graph.Node, data *pb.RequestData) error

type UpdaterFactory interface {
	GetUpdater(nodeType string) (UpdaterType, error)
}

func NewUpdaterFactory() UpdaterFactory {
	return &updaterFactory{}
}

type updaterFactory struct {}

func (f *updaterFactory) GetUpdater(nodeType string) (UpdaterType, error) {
	if _, ok := constructorMap[nodeType]; !ok {
		return nil, fmt.Errorf("not found")
	}
	return updaterMap[nodeType], nil
}

var updaterMap = map[string]UpdaterType{}
