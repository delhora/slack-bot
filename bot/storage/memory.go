package storage

import (
	"encoding/json"
	"fmt"
	"sync"
)

type memoryCollection map[string][]byte

func newMemoryStorage() storage {
	return &memoryStorage{
		storage: make(map[string]memoryCollection),
		mutex:   sync.Mutex{},
	}
}

type memoryStorage struct {
	storage map[string]memoryCollection
	mutex   sync.Mutex
}

func (s memoryStorage) Write(collection, key string, v interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.storage[collection]; !ok {
		s.storage[collection] = make(memoryCollection)
	}

	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	s.storage[collection][key] = data

	return nil
}

func (s memoryStorage) Read(collection, key string, v interface{}) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.storage[collection]; !ok {
		return fmt.Errorf("collection is empty")
	}

	if _, ok := s.storage[collection][key]; !ok {
		return fmt.Errorf("value is empty")
	}

	return json.Unmarshal(s.storage[collection][key], v)
}

func (s memoryStorage) GetKeys(collection string) ([]string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if len(s.storage[collection]) == 0 {
		return nil, fmt.Errorf("collection is empty")
	}

	var keys = make([]string, 0, len(s.storage[collection]))

	for key := range s.storage[collection] {
		keys = append(keys, key)
	}

	return keys, nil
}

func (s memoryStorage) Delete(collection, key string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.storage[collection], key)

	return nil
}
