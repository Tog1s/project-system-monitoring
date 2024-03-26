package storage

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/tog1s/project-system-monitoring/internal/metrics"
)

var ErrObjectExist = errors.New("object already created")

type Store struct {
	mu      sync.RWMutex
	metrics map[uuid.UUID]metrics.SystemMetrics
}

func New() *Store {
	return &Store{
		metrics: make(map[uuid.UUID]metrics.SystemMetrics),
	}
}

func (s *Store) Write(m metrics.SystemMetrics) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.metrics[m.ID]; ok {
		return ErrObjectExist
	}

	s.metrics[m.ID] = m
	return nil
}

func (s *Store) Remove(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.metrics[id]; !ok {
		return nil
	}
	delete(s.metrics, id)
	return nil
}

func (s *Store) All() ([]metrics.SystemMetrics, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metrics := make([]metrics.SystemMetrics, 0, len(s.metrics))
	for _, systemMetrics := range s.metrics {
		metrics = append(metrics, systemMetrics)
	}
	return metrics, nil
}

func (s *Store) Averaged(duration time.Duration) (*metrics.SystemMetricsAverage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	total := 0.0
	var loadAvg1, loadAvg5, loadAvg15 float64
	for _, systemMetrics := range s.metrics {
		if now.Sub(systemMetrics.CollectedAt) <= duration {
			loadAvg1 += systemMetrics.Load.LoadAvg1
			loadAvg5 += systemMetrics.Load.LoadAvg5
			loadAvg15 += systemMetrics.Load.LoadAvg15
			total++
		} else {
			delete(s.metrics, systemMetrics.ID)
		}
	}
	return &metrics.SystemMetricsAverage{
		LoadAvg1:  loadAvg1 / total,
		LoadAvg5:  loadAvg5 / total,
		LoadAvg15: loadAvg15 / total,
	}, nil
}
