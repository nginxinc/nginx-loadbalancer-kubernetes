package synchronization

import (
	"sync"

	v1 "k8s.io/api/core/v1"
)

// cache contains the most recent definitions for services monitored by NLK.
// We need these so that if a service is deleted from the shared informer cache, the
// caller can access the spec of the deleted service for cleanup.
type cache struct {
	mu    sync.RWMutex
	store map[ServiceKey]*v1.Service
}

func newCache() *cache {
	return &cache{
		store: make(map[ServiceKey]*v1.Service),
	}
}

func (s *cache) get(key ServiceKey) (*v1.Service, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	service, ok := s.store[key]
	return service, ok
}

func (s *cache) add(key ServiceKey, service *v1.Service) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = service
}

func (s *cache) delete(key ServiceKey) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, key)
}
