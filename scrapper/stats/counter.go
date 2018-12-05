package stats

import "sync"

type CounterService struct {
	mu sync.Mutex
	counters map[string]int64
}

func NewCounterService() *CounterService {
	return &CounterService{
		sync.Mutex{},
		make(map[string]int64),
	}
}

func (s *CounterService) IncreaseURIStat(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if c, exists := s.counters[key]; !exists {
		s.counters[key] = 1
	} else {
		s.counters[key] = c + 1
	}
}

func (s *CounterService) GetCounters() map[string]int64 {
	return s.counters
}