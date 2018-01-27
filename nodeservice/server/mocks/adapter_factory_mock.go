package mocks

import (
	"fmt"
	"github.com/Sovianum/turbonetwork/nodeservice/server/adapters"
)

func NewNodeAdapterFactoryMock() *NodeAdapterFactoryMock {
	return &NodeAdapterFactoryMock{
		errList:     make([]error, 0),
		adapterList: make([]adapters.NodeAdapter, 0),
	}
}

type NodeAdapterFactoryMock struct {
	cnt         int
	errList     []error
	adapterList []adapters.NodeAdapter
}

func (m *NodeAdapterFactoryMock) GetAdapter(nodeType string) (adapters.NodeAdapter, error) {
	if m.cnt >= len(m.errList) {
		return nil, fmt.Errorf("unexpected call to adapter factory")
	}

	m.cnt++
	return m.adapterList[m.cnt-1], m.errList[m.cnt-1]
}

func (m *NodeAdapterFactoryMock) ExpectResponse(a adapters.NodeAdapter, err error) *NodeAdapterFactoryMock {
	m.errList = append(m.errList, err)
	m.adapterList = append(m.adapterList, a)

	return m
}
