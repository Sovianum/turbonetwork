package mocks

import (
	"fmt"
	"github.com/Sovianum/turbonetwork/nodeservice/server/factories"
)

func NewStateGetterFactoryMock() *StateGetterFactoryMock {
	return &StateGetterFactoryMock{
		getterList: make([]factories.StateGetterType, 0),
		errList:    make([]error, 0),
	}
}

type StateGetterFactoryMock struct {
	cnt int

	getterList []factories.StateGetterType
	errList    []error
}

func (m *StateGetterFactoryMock) ExpectResponse(getter factories.StateGetterType, err error) *StateGetterFactoryMock {
	m.errList = append(m.errList, err)
	m.getterList = append(m.getterList, getter)

	return m
}

func (m *StateGetterFactoryMock) GetStateGetter(nodeType string) (factories.StateGetterType, error) {
	if m.cnt >= len(m.errList) {
		return nil, fmt.Errorf("unexpected call to constructor factory")
	}

	m.cnt++
	return m.getterList[m.cnt-1], m.errList[m.cnt-1]
}
