package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/tog1s/project-system-monitoring/internal/config"
	"github.com/tog1s/project-system-monitoring/internal/metrics"
	"github.com/tog1s/project-system-monitoring/internal/storage"
	"github.com/tog1s/project-system-monitoring/pb"
	"github.com/tog1s/project-system-monitoring/pkg/loadavg"
	"google.golang.org/grpc"
)

type Service struct {
	pb.UnimplementedMetricsServer
}

var (
	store            = storage.New()
	loadAverage bool = true
	addr             = flag.String("addr", "localhost:50051", "listen address and port")
	configFile       = flag.String("config", "configs/config.yaml", "Path to config file")
	startTime   time.Time
)

func newResponse(m *metrics.SystemMetricsAverage) *pb.Response {
	return &pb.Response{
		LoadAvg1:  float32(m.LoadAvg1),
		LoadAvg5:  float32(m.LoadAvg5),
		LoadAvg15: float32(m.LoadAvg15),
	}
}

func startStream(stream pb.Metrics_GetServer, done chan bool,
	message *pb.Request, ticker time.Ticker) {

	duration := time.Duration(time.Duration(message.AverageWindow) * time.Second)

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

func collectMetrics(cfg config.Config) {
	for {
		record := new(metrics.SystemMetrics)
		record.ID = uuid.New()
		record.CollectedAt = time.Now()
		if cfg.Metrics.LoadAverage {
			loadAvg, err := loadavg.Get()
			if err != nil {
				fmt.Println(err)
			}
			record.Load = loadAvg
		}
		err := store.Write(*record)
		if err != nil {
			log.Printf("record writing error: %v", err)
		}
	}
}

func main() {
	flag.Parse()

	cfg, err := config.ReadFromFile(*configFile)
	if err != nil {
		log.Fatalf("error reading configuration file: %s", err)
	}

	go collectMetrics(*cfg)

	lsn, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterMetricsServer(server, &Service{})

	log.Printf("starting server on %s", lsn.Addr().String())
	startTime = time.Now()
	if err := server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}
