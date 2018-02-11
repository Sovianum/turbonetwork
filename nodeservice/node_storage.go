package nodeservice

import (
	"github.com/Sovianum/turbonetwork/common"
	"github.com/Sovianum/turbonetwork/nodeservice/adapters"
	"github.com/Sovianum/turbonetwork/pb"
	"sync"
)

// NodeStorage is a wrapper around ObjectStorage which also
// automatically generates unique ids and casts TypedNode objects to and from interface{}
type NodeStorage interface {
	Add(node *adapters.TypedNode) (*pb.NodeIdentifier, error)
	Get(id *pb.NodeIdentifier) (*adapters.TypedNode, error)
	Drop(id *pb.NodeIdentifier) error
}

// NewMapNodeStorage creates NodeStorage based on map based ObjectStorage
func NewMapNodeStorage() NodeStorage {
	return &mapNodeStorage{
		idCnt:         1,
		idLock:        sync.Mutex{},
		objectStorage: common.NewMapObjectStorage(),
	}
}

type mapNodeStorage struct {
	objectStorage common.ObjectStorage
	idLock        sync.Mutex
	idCnt         int32
}

func (s *mapNodeStorage) Add(node *adapters.TypedNode) (*pb.NodeIdentifier, error) {
	s.idLock.Lock()
	id := pb.NodeIdentifier{Id: s.idCnt, NodeType: node.NodeType}
	s.idCnt++
	s.idLock.Unlock()

	if err := s.objectStorage.Add(id, node); err != nil {
		return nil, err
	}

	return &id, nil
}

func (s *mapNodeStorage) Get(id *pb.NodeIdentifier) (*adapters.TypedNode, error) {
	var (
		val interface{}
		err error
	)
	if val, err = s.objectStorage.Get(*id); err != nil {
		return nil, err
	}
	return val.(*adapters.TypedNode), nil
}

func (s *mapNodeStorage) Drop(id *pb.NodeIdentifier) error {
	return s.objectStorage.Drop(id)
}
