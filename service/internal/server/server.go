package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	pb "github.com/Kochanac/kitime/service/internal/api"
	"github.com/Kochanac/kitime/service/internal/cache"
	"github.com/Kochanac/kitime/service/internal/clickhouse"
	"github.com/Kochanac/kitime/service/internal/kafka"
	"github.com/Kochanac/kitime/service/internal/metrics"
	"github.com/Kochanac/kitime/service/pkg/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type HeadServer struct {
	pb.HeadServer
	Config           config.Config
	KafkaProducer    kafka.Producer
	ClickhouseClient clickhouse.ClickhouseClient
	CacheClient      cache.Cache
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

func (s *HeadServer) Set(ctx context.Context, request *pb.SetRequest) (reply *pb.SetReply, err error) {
	time0 := time.Now()
	defer func() {
		if err == nil {
			metrics.ObserveRequests("set", "ok")
			metrics.ObserveRequestsTimeSum("set", "ok", time.Since(time0).Seconds())
		} else {
			errStatus, _ := status.FromError(err)
			metrics.ObserveRequests("set", errStatus.Code().String())
			metrics.ObserveRequestsTimeSum("set", errStatus.Code().String(), time.Since(time0).Seconds())
		}
	}()

	log.Printf("Received: %v", request)

	data, err := marshalData(request)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to marshal: %w", err)
	}

	err = s.KafkaProducer.Produce(context.WithValue(ctx, "config", s.Config), data)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "Kafka error: %w", err)
	}

	err = s.CacheClient.SaveToCache(ctx,
		&pb.GetRequest{UserId: request.GetUserId(), VideoId: request.GetVideoId()},
		&pb.GetReply{VideoTime: request.VideoTime})
	if err != nil {
		log.Printf("failed to save to cache: %s", err)
	}

	return &pb.SetReply{}, nil
}

func (s *HeadServer) Get(ctx context.Context, request *pb.GetRequest) (resp *pb.GetReply, err error) {
	time0 := time.Now()
	defer func() {
		if err == nil {
			metrics.ObserveRequests("get", "ok")
			metrics.ObserveRequestsTimeSum("set", "ok", time.Since(time0).Seconds())
		} else {
			errStatus, _ := status.FromError(err)
			metrics.ObserveRequests("get", errStatus.Code().String())
			metrics.ObserveRequestsTimeSum("set", errStatus.Code().String(), time.Since(time0).Seconds())
		}
	}()

	fromCache, err := s.CacheClient.CheckCache(ctx, request)
	if err != nil {
		log.Printf("failed to get data from cache: %s", err)
	}

	var res *pb.GetReply
	if fromCache != nil {
		res = fromCache
	} else {
		row, err := s.ClickhouseClient.GetRow(request.GetUserId(), request.GetVideoId())

		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "No data: %s", err)
		}
		if err != nil {
			return nil, status.Errorf(codes.Unavailable, "Failed to receive data from Clickhouse: %s", err)
		}

		res = &pb.GetReply{VideoTime: row.VideoTimestamp}
	}

	return res, nil
}
