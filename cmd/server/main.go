package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/tog1s/project-system-monitoring/internal/config"
	"github.com/tog1s/project-system-monitoring/internal/metrics"
	"github.com/tog1s/project-system-monitoring/internal/storage"
	"github.com/tog1s/project-system-monitoring/pb"
	"github.com/tog1s/project-system-monitoring/pkg/pipeline"
	"google.golang.org/grpc"
)

type Service struct {
	pb.UnimplementedMetricsServer
}

var (
	store      = storage.New()
	addr       = flag.String("addr", "localhost:50051", "listen address and port")
	configFile = flag.String("config", "configs/config.yaml", "Path to config file")
	startTime  time.Time
)

func newResponse(m *metrics.SystemMetricsAverage) *pb.Response {
	return &pb.Response{
		Load: &pb.LoadMessage{
			LoadAvg1:  float32(m.LoadAvg1),
			LoadAvg5:  float32(m.LoadAvg5),
			LoadAvg15: float32(m.LoadAvg15),
		},
		Cpu: &pb.CPUMessage{
			User:   float32(m.CPUUser),
			System: float32(m.CPUSystem),
			Idle:   float32(m.CPUIdle),
		},
	}
}

func collectMetrics(cfg config.Config) {
	in := make(pipeline.Bi)
	done := make(chan bool)
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				m := metrics.SystemMetrics{
					ID:          uuid.New(),
					CollectedAt: time.Now(),
				}
				in <- m
			case <-done:
				return
			}
		}
	}()

	stages := metrics.Collect(cfg)
	go func(storage *storage.Store) {
		for metric := range pipeline.ExecutePipeline(in, nil, stages...) {
			err := storage.Write(metric.(metrics.SystemMetrics))
			if err != nil {
				log.Printf("record writing error: %v", err)
			}
		}
	}(store)
}

func startStream(stream pb.Metrics_GetServer, done chan bool, message *pb.Request, ticker time.Ticker) {
	duration := time.Duration(message.AverageWindow) * time.Second

	if time.Since(startTime) < duration {
		time.Sleep(duration - time.Since(startTime))
	}

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
}

func (s *Service) Get(message *pb.Request, stream pb.Metrics_GetServer) error {
	ticker := time.NewTicker(time.Duration(message.ScrapeInterval) * time.Second)

	done := make(chan bool)
	go startStream(stream, done, message, *ticker)
	<-done

	return nil
}

func main() {
	flag.Parse()

	cfg, err := config.ReadFromFile(*configFile)
	if err != nil {
		log.Fatalf("error reading configuration file: %s", err)
	}

	lsn, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMetricsServer(server, &Service{})

	go collectMetrics(*cfg)

	log.Printf("starting server on %s", lsn.Addr().String())
	startTime = time.Now()
	if err := server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}
