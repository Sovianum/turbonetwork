package mocks

import (
	"fmt"
	"github.com/Sovianum/turbonetwork/nodeservice/adapters"
)

// NewNodeAdapterFactoryMock constructs empty NodeAdapterFactoryMock
func NewNodeAdapterFactoryMock() *NodeAdapterFactoryMock {
	return &NodeAdapterFactoryMock{
		errList:     make([]error, 0),
		adapterList: make([]adapters.NodeAdapter, 0),
	}
}

// NodeAdapterFactoryMock mocks NodeAdapterFactory interface
type NodeAdapterFactoryMock struct {
	cnt         int
	errList     []error
	adapterList []adapters.NodeAdapter
}

// GetAdapter returns adapters in order of expectations
func (m *NodeAdapterFactoryMock) GetAdapter(nodeType string) (adapters.NodeAdapter, error) {
	if m.cnt >= len(m.errList) {
		return nil, fmt.Errorf("unexpected call to adapter factory")
	}

	m.cnt++
	return m.adapterList[m.cnt-1], m.errList[m.cnt-1]
}

// ExpectResponse saves adapters and errors to mock
func (m *NodeAdapterFactoryMock) ExpectResponse(a adapters.NodeAdapter, err error) *NodeAdapterFactoryMock {
	m.errList = append(m.errList, err)
	m.adapterList = append(m.adapterList, a)

	return m
}
