package common

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type NodeStorageTestSuite struct {
	suite.Suite
	storage *mapObjectStorage
}

func (s *NodeStorageTestSuite) SetupTest() {
	s.storage = NewMapObjectStorage().(*mapObjectStorage)
}

func (s *NodeStorageTestSuite) TestAdd_OK() {
	key, value := "key", "value"

	err := s.storage.Add(key, value)

	s.Require().Nil(err)
	s.Equal(s.storage.objectMap[key], value)
}

func (s *NodeStorageTestSuite) TestAdd_Duplicate() {
	key1, value1 := "key", "value1"
	key2, value2 := "key", "value2"

	err1 := s.storage.Add(key1, value1)
	err2 := s.storage.Add(key2, value2)

	s.Require().Nil(err1)
	s.Require().Error(err2)
	s.Equal(s.storage.objectMap[key1], value1)
}

func (s *NodeStorageTestSuite) TestDelete() {
	key1, value1 := "key1", "value1"
	key2, value2 := "key2", "value2"

	s.storage.Add(key1, value1)
	s.storage.Add(key2, value2)

	s.storage.Drop(key1)

	s.Equal(1, len(s.storage.objectMap))
}

func (s *NodeStorageTestSuite) TestGet_OK() {
	key, value := "key", "value"

	err := s.storage.Add(key, value)

	s.Require().Nil(err)

	gotValue, err := s.storage.Get(key)

	s.Require().Nil(err)
	s.Equal(value, gotValue)
}

func (s *NodeStorageTestSuite) TestGet_NotFound() {
	key := "key"

	_, err := s.storage.Get(key)
	s.Require().Error(err)
}

func TestNodeStorageTestSuite(t *testing.T) {
	suite.Run(t, new(NodeStorageTestSuite))
}
