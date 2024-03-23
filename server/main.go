package main

import (
	"context"
	"log"
	"net"

	"github.com/tog1s/project-system-monitoring/pb"
	"github.com/tog1s/project-system-monitoring/pkg/loadavg"
	"google.golang.org/grpc"
)

type Service struct {
	pb.UnimplementedMetricsServer
}

func (s *Service) Get(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	la := loadavg.Read()
	return &pb.Response{
		OneMinutes: la.OneMinutes,
	}, nil
}

func main() {
	lsn, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	pb.RegisterMetricsServer(server, new(Service))

	log.Printf("starting server on %s", lsn.Addr().String())
	if err := server.Serve(lsn); err != nil {
		log.Fatal(err)
	}
}
