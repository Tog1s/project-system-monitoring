package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/tog1s/project-system-monitoring/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewMetricsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &pb.Request{Query: ""})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("LoadAverage one minutes: %f", r.GetOneMinutes())
}
