package server

import (
	"github.com/Sovianum/turbonetwork/common"
	"github.com/Sovianum/turbonetwork/nodeservice/pb"
	"github.com/Sovianum/turbonetwork/nodeservice/server/adapters"
	"sync"
)

type NodeStorage interface {
	Add(node *adapters.TypedNode) (*pb.NodeIdentifier, error)
	Get(id *pb.NodeIdentifier) (*adapters.TypedNode, error)
	Drop(id *pb.NodeIdentifier) error
}

func NewMapNodeStorage() NodeStorage {
	return &mapNodeStorage{
		idCnt:         1,
		mapLock:       sync.Mutex{},
		objectStorage: common.NewMapObjectStorage(),
	}
}

type mapNodeStorage struct {
	objectStorage common.ObjectStorage
	mapLock       sync.Mutex
	idCnt         int32
}

func (s *mapNodeStorage) Add(node *adapters.TypedNode) (*pb.NodeIdentifier, error) {
	s.mapLock.Lock()
	id := pb.NodeIdentifier{Id: s.idCnt, NodeType: node.NodeType}
	s.idCnt++
	s.mapLock.Unlock()

	if err := s.objectStorage.Add(id, node); err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *mapNodeStorage) Get(id *pb.NodeIdentifier) (*adapters.TypedNode, error) {
	if val, err := s.objectStorage.Get(*id); err != nil {
		return nil, err
	} else {
		return val.(*adapters.TypedNode), nil
	}
}

func (s *mapNodeStorage) Drop(id *pb.NodeIdentifier) error {
	return s.objectStorage.Drop(id)
}
