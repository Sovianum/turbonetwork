package common

import (
	"fmt"
	"sync"
)

type ObjectStorage interface {
	Add(key, value interface{}) error
	Get(key interface{}) (interface{}, error)
	Drop(key interface{}) error
}

func NewMapObjectStorage() ObjectStorage {
	return &mapObjectStorage{
		mapLock:   sync.Mutex{},
		objectMap: make(map[interface{}]interface{}),
	}
}

type mapObjectStorage struct {
	mapLock   sync.Mutex
	objectMap map[interface{}]interface{}
}

func (s *mapObjectStorage) Add(key, value interface{}) error {
	s.mapLock.Lock()
	defer s.mapLock.Unlock()

	if _, ok := s.objectMap[key]; ok {
		return fmt.Errorf("duplicate key")
	}

	s.objectMap[key] = value
	return nil
}

func (s *mapObjectStorage) Get(key interface{}) (interface{}, error) {
	s.mapLock.Lock()
	defer s.mapLock.Unlock()

	if _, ok := s.objectMap[key]; !ok {
		return nil, fmt.Errorf("not found object with key %v", key)
	}
	return s.objectMap[key], nil
}

func (s *mapObjectStorage) Drop(key interface{}) error {
	s.mapLock.Lock()
	defer s.mapLock.Unlock()

	delete(s.objectMap, key)
	return nil
}
