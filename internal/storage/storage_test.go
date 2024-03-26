package storage

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tog1s/project-system-monitoring/internal/metrics"
	"github.com/tog1s/project-system-monitoring/pkg/loadavg"
)

func TestWrite(t *testing.T) {
	t.Run("test write func", func(t *testing.T) {
		metric := metrics.SystemMetrics{
			ID:          uuid.New(),
			CollectedAt: time.Now(),
			Load: &loadavg.LoadAverage{
				LoadAvg1:  1.0,
				LoadAvg5:  1.0,
				LoadAvg15: 1.0,
			},
		}
		store := New()
		err := store.Write(metric)
		require.NoError(t, err)
	})
}
