package server

import (
	"github.com/stretchr/testify/suite"
	"github.com/Sovianum/turbocycle/core/graph"
	"github.com/Sovianum/turbonetwork/nodeservice/server/factories"
	"testing"
)

type NodeStorageTestSuite struct {
	suite.Suite
	storage *mapNodeStorage
}

func (s *NodeStorageTestSuite) SetupTest() {
	s.storage = NewMapNodeStorage().(*mapNodeStorage)
}

func (s *NodeStorageTestSuite) TestAdd() {
	id, err := s.storage.Add(factories.NewTypedNode(
		graph.NewTestNode(0, 0, true, nil),
		"test",
	))

	s.Require().Nil(err)
	s.Equal(id.NodeType, "test")
	s.EqualValues(id.Id, 1)
}

func (s *NodeStorageTestSuite) TestDelete() {
	id, err := s.storage.Add(factories.NewTypedNode(
		graph.NewTestNode(0, 0, true, nil),
		"test",
	))

	s.Require().Equal(len(s.storage.nodeMap), 1)
	s.Require().Nil(err)

	err = s.storage.Drop(id)
	s.Equal(len(s.storage.nodeMap), 0)
}

func (s *NodeStorageTestSuite) TestGet() {
	inputNode := graph.NewTestNode(0, 0, true, nil)

	id, err := s.storage.Add(factories.NewTypedNode(inputNode, "test"))

	s.Require().Equal(len(s.storage.nodeMap), 1)
	s.Require().Nil(err)

	node, err := s.storage.Get(id)
	s.Equal(node.NodeType, "test")
	s.Equal(node.Node, inputNode)
}

func TestNodeStorageTestSuite(t *testing.T) {
	suite.Run(t, new(NodeStorageTestSuite))
}
