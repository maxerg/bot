package auth

import (
	"context"
	"sync"
	"time"
)

type MemStore struct {
	mu sync.RWMutex
	m  map[string]Session
}

func NewMemStore() *MemStore {
	return &MemStore{m: make(map[string]Session)}
}

func (s *MemStore) Create(ctx context.Context, sess Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[sess.ID] = sess
	return nil
}

func (s *MemStore) Get(ctx context.Context, id string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[id]
	if !ok {
		return Session{}, false
	}
	if time.Now().After(v.ExpiresAt) {
		return Session{}, false
	}
	return v, true
}

func (s *MemStore) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, id)
	return nil
}
