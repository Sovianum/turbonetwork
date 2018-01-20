package mocks

import (
	"fmt"
	"github.com/Sovianum/turbonetwork/nodeservice/server/factories"
)

func NewUpdaterFactoryMock() *UpdaterFactoryMock {
	return &UpdaterFactoryMock{
		updaterList: make([]factories.UpdaterType, 0),
		errList:    make([]error, 0),
	}
}

type UpdaterFactoryMock struct {
	cnt int

	updaterList []factories.UpdaterType
	errList     []error
}

func (m *UpdaterFactoryMock) ExpectResponse(updater factories.UpdaterType, err error) *UpdaterFactoryMock {
	m.updaterList = append(m.updaterList, updater)
	m.errList = append(m.errList, err)

	return m
}

func (m *UpdaterFactoryMock) GetUpdater(nodeType string) (factories.UpdaterType, error) {
	if m.cnt >= len(m.errList) {
		return nil, fmt.Errorf("unexpected call to constructor factory")
	}

	m.cnt++
	return m.updaterList[m.cnt-1], m.errList[m.cnt-1]
}
