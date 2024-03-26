package metrics

import (
	"time"

	"github.com/google/uuid"
	"github.com/tog1s/project-system-monitoring/pkg/loadavg"
)

type SystemMetricsAverage struct {
	LoadAvg1  float64
	LoadAvg5  float64
	LoadAvg15 float64
}

type SystemMetrics struct {
	ID          uuid.UUID
	CollectedAt time.Time
	Load        *loadavg.LoadAverage
}
