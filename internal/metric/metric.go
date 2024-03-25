package metric

import (
	"time"

	"github.com/tog1s/project-system-monitoring/pkg/loadavg"
)

type SystemMetricsAverage struct {
}

type SystemMetrics struct {
	CollectedAt time.Time
	Load        *loadavg.LoadAverage
}
