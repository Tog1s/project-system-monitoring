package metrics

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tog1s/project-system-monitoring/internal/config"
)

func TestMetrics(t *testing.T) {
	t.Run("test empty config", func(t *testing.T) {
		cfg := config.Config{
			Metrics: config.MetricsConfig{
				LoadAverage: false,
			},
		}

		result := Collect(cfg)
		require.Empty(t, result)
	})

	t.Run("test collection", func(t *testing.T) {
		cfg := config.Config{
			Metrics: config.MetricsConfig{
				LoadAverage: true,
			},
		}

		result := Collect(cfg)
		require.Len(t, result, 1)
	})
}
