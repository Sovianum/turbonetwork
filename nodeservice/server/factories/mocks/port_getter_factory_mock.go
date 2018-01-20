package mocks

import (
	"fmt"
	"github.com/Sovianum/turbonetwork/nodeservice/server/factories"
)

func NewPortGetterFactoryMock() *PortGetterFactoryMock {
	return &PortGetterFactoryMock{
		errList:    make([]error, 0),
		getterList: make([]factories.PortGetter, 0),
	}
}

type PortGetterFactoryMock struct {
	cnt int

	errList    []error
	getterList []factories.PortGetter
}

func (m *PortGetterFactoryMock) ExpectResponse(getter factories.PortGetter, err error) *PortGetterFactoryMock {
	m.errList = append(m.errList, err)
	m.getterList = append(m.getterList, getter)

	return m
}

func (m *PortGetterFactoryMock) GetPortGetter(nodeType string) (factories.PortGetter, error) {
	if m.cnt >= len(m.errList) {
		return nil, fmt.Errorf("unexpected call to getter factory")
	}

	m.cnt++
	return m.getterList[m.cnt-1], m.errList[m.cnt-1]
}
