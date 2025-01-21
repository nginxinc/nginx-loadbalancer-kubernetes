package synchronization

import (
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
)

// cache contains the most recent definitions for services monitored by NLK.
// We need these so that if a service is deleted from the shared informer cache, the
// caller can access the spec of the deleted service for cleanup.
type cache struct {
	mu    sync.RWMutex
	store map[ServiceKey]service
}

type service struct {
	service *v1.Service
	// removedAt indicates when the service was removed from NGINXaaS
	// monitoring. A zero time indicates that the service is still actively
	// being monitored by NGINXaaS.
	removedAt time.Time
}

func newCache() *cache {
	return &cache{
		store: make(map[ServiceKey]service),
	}
}

func (s *cache) get(key ServiceKey) (service, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	svc, ok := s.store[key]
	return svc, ok
}

func (s *cache) add(key ServiceKey, service service) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = service
}

func (s *cache) delete(key ServiceKey) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, key)
}
