package mocks

import (
	"github.com/Sovianum/turbonetwork/nodeservice/server/adapters"
	"github.com/Sovianum/turbonetwork/pb"
	"fmt"
)

type Pair struct {
	First interface{}
	Second interface{}
}

func NewNodeStorageMock() *NodeStorageMock {
	return &NodeStorageMock{
		addResponses:make([]Pair, 0),
		getResponses:make([]Pair, 0),
		dropResponses:make([]error, 0),
	}
}

type NodeStorageMock struct {
	addCnt int
	getCnt int
	dropCnt int

	addResponses []Pair
	getResponses []Pair
	dropResponses []error
}

func (m *NodeStorageMock) ExpectAddResponse(id *pb.NodeIdentifier, err error) *NodeStorageMock {
	m.addResponses = append(m.addResponses, Pair{id, err})
	return m
}

func (m *NodeStorageMock) ExpectGetResponse(node *adapters.TypedNode, err error) *NodeStorageMock {
	m.getResponses = append(m.getResponses, Pair{node, err})
	return m
}

func (m *NodeStorageMock) ExpectDropResponse(err error) *NodeStorageMock {
	m.dropResponses = append(m.dropResponses, err)
	return m
}

func (m *NodeStorageMock) Add(node *adapters.TypedNode) (*pb.NodeIdentifier, error) {
	if m.addCnt >= len(m.addResponses) {
		return nil, fmt.Errorf("unexpected add request")
	}
	r := m.addResponses[m.addCnt]
	m.addCnt++

	var err error
	if r.Second != nil {
		err = r.Second.(error)
	}

	return r.First.(*pb.NodeIdentifier), err
}

func (m *NodeStorageMock) Get(id *pb.NodeIdentifier) (*adapters.TypedNode, error) {
	if m.getCnt >= len(m.getResponses) {
		return nil, fmt.Errorf("unexpected get request")
	}
	r := m.getResponses[m.getCnt]
	m.getCnt++

	var err error
	if r.Second != nil {
		err = r.Second.(error)
	}

	return r.First.(*adapters.TypedNode), err
}

func (m *NodeStorageMock) Drop(id *pb.NodeIdentifier) error {
	if m.dropCnt >= len(m.dropResponses) {
		return fmt.Errorf("unexpected drop request")
	}
	return m.dropResponses[m.addCnt]
}

