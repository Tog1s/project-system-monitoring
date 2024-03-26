package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/tog1s/project-system-monitoring/internal/metrics"
	"github.com/tog1s/project-system-monitoring/internal/storage"
	"github.com/tog1s/project-system-monitoring/pb"
	"github.com/tog1s/project-system-monitoring/pkg/loadavg"

	"google.golang.org/grpc"
)

type Service struct {
	pb.UnimplementedMetricsServer
}

var store = storage.New()

var loadAverage bool = true

func newResponse(m *metrics.SystemMetricsAverage) *pb.Response {
	return &pb.Response{
		LoadAvg1:  float32(m.LoadAvg1),
		LoadAvg5:  float32(m.LoadAvg5),
		LoadAvg15: float32(m.LoadAvg15),
	}
}

func (s *Service) Get(message *pb.Request, stream pb.Metrics_GetServer) error {
	duration, err := time.ParseDuration("15s")
	if err != nil {
		fmt.Println(err)
	}

	done := make(chan bool)
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		time.Sleep(15 * time.Second)
		for {
			select {
			case <-ticker.C:
				m, err := store.Averaged(duration)
				if err != nil {
					log.Printf("error %s", err)
					done <- true
				}

				err = stream.Send(newResponse(m))
				if err != nil {
					log.Printf("send error %s", err)
					done <- true
					return
				}
				log.Printf("metrics sended")

			case <-done:
				return
			}
		}
	}()
	<-done
	return nil
}

func collectMetrics() {
	for {
		record := new(metrics.SystemMetrics)
		record.ID = uuid.New()
		record.CollectedAt = time.Now()
		if loadAverage {
			loadAvg, err := loadavg.Get()
			if err != nil {
				fmt.Println(err)
			}
			record.Load = loadAvg
		}
		store.Write(*record)
	}
}

func main() {
	go collectMetrics()

	lsn, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMetricsServer(server, &Service{})

	log.Printf("starting server on %s", lsn.Addr().String())
	if err := server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}
