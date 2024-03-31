package storage

import (
	"math"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tog1s/project-system-monitoring/internal/metrics"
	"github.com/tog1s/project-system-monitoring/pkg/loadavg"
)

func TestStorage(t *testing.T) {
	t.Run("test write, all, remove", func(t *testing.T) {
		some_metrics := metrics.SystemMetrics{
			ID:          uuid.New(),
			CollectedAt: time.Now(),
			Load: &loadavg.LoadAverage{
				LoadAvg1:  1.0,
				LoadAvg5:  1.0,
				LoadAvg15: 1.0,
			},
		}
		store := New()
		err := store.Write(some_metrics)
		require.NoError(t, err)

		m, err := store.All()
		require.NoError(t, err)
		require.Len(t, m, 1)
		require.Equal(t, some_metrics, m[0])

		err = store.Remove(some_metrics.ID)
		require.NoError(t, err)

		m, err = store.All()
		require.NoError(t, err)
		require.Len(t, m, 0)
	})

	t.Run("test averaged", func(t *testing.T) {
		store := New()
		some_metrics := []metrics.SystemMetrics{
			{
				ID:          uuid.New(),
				CollectedAt: time.Now().Add(-10 * time.Second),
				Load: &loadavg.LoadAverage{
					LoadAvg1:  2.0,
					LoadAvg5:  2.0,
					LoadAvg15: 2.0,
				},
			},
			{
				ID:          uuid.New(),
				CollectedAt: time.Now(),
				Load: &loadavg.LoadAverage{
					LoadAvg1:  1.0,
					LoadAvg5:  1.0,
					LoadAvg15: 1.0,
				},
			},
			{
				ID:          uuid.New(),
				CollectedAt: time.Now(),
				Load: &loadavg.LoadAverage{
					LoadAvg1:  1.5,
					LoadAvg5:  1.5,
					LoadAvg15: 1.5,
				},
			},
		}

		for _, v := range some_metrics {
			err := store.Write(v)
			if err != nil {
				t.FailNow()
				return
			}
		}

		require.Len(t, store.metrics, 3)
		avg, err := store.Averaged(5 * time.Second)
		require.NoError(t, err)
		require.Len(t, store.metrics, 2)

		require.Equal(t, math.Round(avg.LoadAvg1*100)/100, 1.25)
		require.Equal(t, math.Round(avg.LoadAvg5*100)/100, 1.25)
		require.Equal(t, math.Round(avg.LoadAvg15*100)/100, 1.25)

	})
}
