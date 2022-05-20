package main

import (
	"flag"
	"fmt"
	pb "github.com/Kochanac/kitime/service/internal/api"
	"github.com/Kochanac/kitime/service/internal/cache"
	"github.com/Kochanac/kitime/service/internal/clickhouse"
	"github.com/Kochanac/kitime/service/internal/kafka"
	"github.com/Kochanac/kitime/service/internal/metrics"
	"github.com/Kochanac/kitime/service/internal/server"
	"github.com/Kochanac/kitime/service/pkg/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net"
	"net/http"
	"time"

	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 9666, "The server port")
)

func main() {
	metrics.Init()

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
		ClickhouseClient: clickhouse.Init(c.GetClickhouseHost(), c.GetClickhouseUser(), c.GetClickhousePassword()),
		CacheClient:      cache.InitCache(c.GetRedisHost(), time.Minute*10),
	})

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		err := http.ListenAndServe(":9100", nil)
		if err != nil {
			log.Printf("Error at http server: %s", err)
		}
	}()

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
