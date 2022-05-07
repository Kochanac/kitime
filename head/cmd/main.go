package main

import (
	"flag"
	"fmt"
	pb "github.com/Kochanac/vobla/head/internal/api"
	"github.com/Kochanac/vobla/head/internal/clickhouse"
	"github.com/Kochanac/vobla/head/internal/kafka"
	"github.com/Kochanac/vobla/head/internal/server"
	"github.com/Kochanac/vobla/head/pkg/config"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 9666, "The server port")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	c := config.Config{}
	producer, err := kafka.InitProducer(c)
	if err != nil {
		log.Fatalf("failed to create kafka producer: %s", err)
	}

	pb.RegisterHeadServer(s, &server.HeadServer{
		Config:           c,
		KafkaProducer:    producer,
		ClickhouseClient: clickhouse.Init(c.GetClickhouseHost()),
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
