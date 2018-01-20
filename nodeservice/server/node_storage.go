package server

import (
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"sync"
	"fmt"
	"github.com/Sovianum/turbonetwork/nodeservice/server/factories"
)


type NodeStorage interface {
	Add(node *factories.TypedNode) (*pb.NodeIdentifier, error)
	Get(id *pb.NodeIdentifier) (*factories.TypedNode, error)
	Drop(id *pb.NodeIdentifier) error
}

func NewMapNodeStorage() NodeStorage {
	return &mapNodeStorage{
		idCnt:1,
		mapLock:sync.Mutex{},
		nodeMap:make(map[pb.NodeIdentifier]*factories.TypedNode),
	}
}

type mapNodeStorage struct {
	mapLock sync.Mutex
	idCnt int32
	nodeMap map[pb.NodeIdentifier]*factories.TypedNode
}

func (s *mapNodeStorage) Add(node *factories.TypedNode) (*pb.NodeIdentifier, error) {
	s.mapLock.Lock()
	defer s.mapLock.Unlock()

	id := &pb.NodeIdentifier{Id:s.idCnt, NodeType:node.NodeType}

	s.nodeMap[*id] = node
	s.idCnt++

	return id, nil
}

func (s *mapNodeStorage) Get(id *pb.NodeIdentifier) (*factories.TypedNode, error) {
	s.mapLock.Lock()
	defer s.mapLock.Unlock()

	if _, ok := s.nodeMap[*id]; !ok {
		return nil, fmt.Errorf("not found node with id %d", id.Id)
	}
	return s.nodeMap[*id], nil
}

func (s *mapNodeStorage) Drop(id *pb.NodeIdentifier) error {
	s.mapLock.Lock()
	defer s.mapLock.Unlock()

	delete(s.nodeMap, *id)
	return nil
}



