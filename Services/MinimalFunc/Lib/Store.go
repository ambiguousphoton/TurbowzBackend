package store

import "sync"

type Function struct{
	ID  string
	Code  string
}

type MemoryStore struct {
	mu               sync.RWMutex
	functions        map[string]Functions
}

func NewStore () *MemoryStore{
	return &MemoryStore{
		functions: make(map[string]Function),
	}
}

func (s *MemoryStore) Save(fn Function){
	s.mu.Lock()
	defer s.mu.Unlock()
	s.function[Fn.ID] = fn
}

func (s *MemoryStore) Get(id string)(Function, bool){
	s.mu.Lock()
	defer s.mu.Unlock()
	fn, ok := s.functions[id]
	return fn, ok
}