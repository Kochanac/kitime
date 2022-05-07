package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	pb "github.com/Kochanac/vobla/head/internal/api"
	"github.com/Kochanac/vobla/head/internal/clickhouse"
	"github.com/Kochanac/vobla/head/internal/kafka"
	"github.com/Kochanac/vobla/head/pkg/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
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
	log.Printf("Received: %v", request)

	data, err := marshalData(request)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to marshal: %w", err)
	}

	err = s.KafkaProducer.Produce(context.WithValue(ctx, "config", s.Config), data)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "Kafka error: %w", err)
	}

	return &pb.SetReply{}, nil
}

func (s *HeadServer) Get(ctx context.Context, request *pb.GetRequest) (*pb.GetReply, error) {
	row, err := s.ClickhouseClient.GetRow(request.GetUserId(), request.GetVideoId())

	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Errorf(codes.NotFound, "No data: %s", err)
	}
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "Failed to receive data from Clickhouse: %s", err)
	}

	return &pb.GetReply{VideoTime: row.VideoTimestamp}, nil
}
