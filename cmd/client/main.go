package main

import (
	"context"
	"errors"
	"flag"
	"io"
	"log"

	"github.com/tog1s/project-system-monitoring/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr           = flag.String("addr", "localhost:50051", "the address to connect to")
	scrapeInterval = flag.Int("scrape", 5, "how frequently to scrape metrics from server (seconds)")
	averageWindow  = flag.Int("average", 15, "average window (seconds)")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewMetricsClient(conn)
	stream, err := client.Get(context.Background(), &pb.Request{
		Query:          "",
		ScrapeInterval: int32(*scrapeInterval),
		AverageWindow:  int32(*averageWindow),
	})
	//nolint:gocritic
	if err != nil {
		log.Fatalln("client error", err)
	}

	done := make(chan bool)
	go func() {
		for {
			metrics, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				log.Fatalln("stream recive error", err)
			}
			log.Println(metrics)
		}
	}()

	<-done
	log.Printf("finish recive")
}
