package storage

import (
	"sync"
	"time"

	"github.com/tog1s/project-system-monitoring/internal/metric"
)

type Store struct {
	mu      sync.RWMutex
	metrics map[time.Time]metric.SystemMetrics
}

func New() *Store {
	return &Store{
		metrics: make(map[time.Time]metric.SystemMetrics),
	}
}

func (s *Store) Write(m metric.SystemMetrics) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.metrics[m.CollectedAt]; ok {
		return nil
	}

	s.metrics[m.CollectedAt] = m
	return nil
}

func (s *Store) Remove() error {
	return nil
}

func (s *Store) All() error {
	return nil
}

func (s *Store) Averaged() error {
	return nil
}
