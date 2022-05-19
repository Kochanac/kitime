package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	pb "github.com/Kochanac/kitime/service/internal/api"
	"github.com/Kochanac/kitime/service/internal/clickhouse"
	"github.com/Kochanac/kitime/service/internal/kafka"
	"github.com/Kochanac/kitime/service/internal/metrics"
	"github.com/Kochanac/kitime/service/pkg/config"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
	"time"
)

type HeadServer struct {
	pb.HeadServer
	Config           config.Config
	KafkaProducer    kafka.Producer
	ClickhouseClient clickhouse.Client
}

func abs(i int64) int64 {
	if i < 0 {
		return -i
	} else {
		return i
	}
}

func marshalData(dataRaw *pb.SetRequest) (string, error) {
	type marshalFormat struct {
		UserId         uint32 `json:"user_id"`
		EventTime      uint32 `json:"event_time"`
		EventType      uint8  `json:"event_type"`
		VideoId        uint32 `json:"video_id"`
		VideoTimestamp uint32 `json:"video_timestamp"`
	}

	var data marshalFormat
	data.UserId = dataRaw.GetUserId()
	data.EventTime = uint32(abs(dataRaw.GetEventTime().Seconds))
	data.EventType = uint8(dataRaw.GetEventType())
	data.VideoId = dataRaw.GetVideoId()
	data.VideoTimestamp = dataRaw.GetVideoTime()

	marshalled, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(marshalled), nil
}

func (s *HeadServer) Set(ctx context.Context, request *pb.SetRequest) (*pb.SetReply, error) {
	time0 := time.Now()

	log.Printf("Received: %v", request)

	data, err := marshalData(request)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to marshal: %w", err)
	}

	err = s.KafkaProducer.Produce(context.WithValue(ctx, "config", s.Config), data)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "Kafka error: %w", err)
	}

	metrics.RequestsMetric.With(prometheus.Labels{
		"type":      "set",
		"resp_type": "ok",
		"time":      strconv.Itoa(int(time.Since(time0).Milliseconds())),
	}).Inc()

	metrics.RequestsTimeSum.With(prometheus.Labels{
		"type":      "set",
		"resp_type": "ok",
	}).Add(float64(time.Since(time0).Milliseconds()))
	return &pb.SetReply{}, nil
}

func (s *HeadServer) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetReply, error) {
	time0 := time.Now()

	row, err := s.ClickhouseClient.GetRow(request.GetUserId(), request.GetVideoId())

	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Errorf(codes.NotFound, "No data: %s", err)
	}
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "Failed to receive data from Clickhouse: %s", err)
	}

	metrics.RequestsMetric.With(prometheus.Labels{
		"type":      "get",
		"resp_type": "ok",
	}).Inc()
	metrics.RequestsTimeSum.With(prometheus.Labels{
		"type":      "get",
		"resp_type": "ok",
	}).Add(float64(time.Since(time0).Milliseconds()))
	return &pb.GetReply{VideoTime: row.VideoTimestamp}, nil
}
