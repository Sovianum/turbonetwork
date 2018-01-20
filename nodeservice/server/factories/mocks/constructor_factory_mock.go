package mocks

import (
	"github.com/Sovianum/turbonetwork/nodeservice/server/factories"
	"fmt"
)


func NewConstructorFactoryMock() *ConstructorFactoryMock {
	return &ConstructorFactoryMock{
		errList: make([]error, 0),
		constructorList: make([]factories.ConstructorType, 0),
	}
}

type ConstructorFactoryMock struct {
	cnt int

	errList []error
	constructorList []factories.ConstructorType
}

func (m *ConstructorFactoryMock) ExpectResponse(c factories.ConstructorType, err error) *ConstructorFactoryMock {
	m.errList = append(m.errList, err)
	m.constructorList = append(m.constructorList, c)

	return m
}

func (m *ConstructorFactoryMock) GetConstructor(nodeType string) (factories.ConstructorType, error) {
	if m.cnt >= len(m.errList) {
		return nil, fmt.Errorf("unexpected call to constructor factory")
	}

	m.cnt++
	return m.constructorList[m.cnt-1], m.errList[m.cnt-1]
}

