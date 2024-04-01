package metrics

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/tog1s/project-system-monitoring/internal/config"
	"github.com/tog1s/project-system-monitoring/pkg/cpustat"
	"github.com/tog1s/project-system-monitoring/pkg/loadavg"
	"github.com/tog1s/project-system-monitoring/pkg/pipeline"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type SystemMetricsAverage struct {
	LoadAvg1  float64
	LoadAvg5  float64
	LoadAvg15 float64
	CPUUser   float64
	CPUSystem float64
	CPUIdle   float64
}

type SystemMetrics struct {
	ID          uuid.UUID
	CollectedAt time.Time
	Load        *loadavg.LoadAverage
	CPUStat     *cpustat.CPUStat
}

func Collect(cfg config.Config) []pipeline.Stage {
	var stages []pipeline.Stage

	if cfg.Metrics.LoadAverage {
		stages = append(stages, stageGenerator(
			"Load average",
			func(m SystemMetrics) SystemMetrics {
				la, err := loadavg.Get()
				if err != nil {
					log.Printf("failed to get loadavg: %s", err)
				}
				m.Load = la
				return m
			}),
		)
	}

	if cfg.Metrics.CPU {
		stages = append(stages, stageGenerator(
			"CPU",
			func(m SystemMetrics) SystemMetrics {
				cpu, err := cpustat.Get()
				if err != nil {
					log.Printf("failed to get cpustats: %s", err)
				}
				m.CPUStat = cpu
				return m
			}),
		)
	}
	return stages
}

func stageGenerator(_ string, f func(m SystemMetrics) SystemMetrics) pipeline.Stage {
	return func(in In) Out {
		out := make(Bi)
		go func() {
			defer close(out)
			for v := range in {
				out <- f(v.(SystemMetrics))
			}
		}()
		return out
	}
}
