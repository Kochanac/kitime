package cache

import (
	"context"
	"fmt"
	pb "github.com/Kochanac/kitime/service/internal/api"
	"github.com/Kochanac/kitime/service/internal/metrics"
	"github.com/gomodule/redigo/redis"
	"strconv"
	"time"
)

const (
	keyFormat = "uid:%d/vid:%d/latest"
)

type Cache interface {
	CheckCache(ctx context.Context, req *pb.GetRequest) (*pb.GetReply, error)
	SaveToCache(ctx context.Context, req *pb.GetRequest, reply *pb.GetReply) error
}

type RedisCache struct {
	pool            redis.Pool
	cacheExpireTime time.Duration
}

func InitCache(redisAddress string, cacheExpireTime time.Duration) *RedisCache {
	return &RedisCache{
		pool: redis.Pool{
			MaxIdle:   4,
			MaxActive: 8,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", redisAddress)
				if err != nil {
					panic(err.Error())
				}
				return c, err
			},
		},
		cacheExpireTime: cacheExpireTime,
	}
}

func (c *RedisCache) CheckCache(ctx context.Context, req *pb.GetRequest) (reply *pb.GetReply, err error) {
	defer metrics.ObserveCacheHits(reply != nil)

	client := c.pool.Get()
	defer client.Close()

	key, err := client.Do("get", fmt.Sprintf(keyFormat, req.GetUserId(), req.GetVideoId()))
	if err != nil {
		return nil, err
	}

	if key == nil {
		return nil, nil
	}

	videoTime, err := strconv.ParseUint(string(key.([]byte)), 10, 64)
	if err != nil {
		return nil, err
	}

	return &pb.GetReply{VideoTime: uint32(videoTime)}, nil
}

func (c *RedisCache) SaveToCache(ctx context.Context, req *pb.GetRequest, reply *pb.GetReply) error {
	client := c.pool.Get()
	defer client.Close()

	_, err := client.Do("set",
		fmt.Sprintf(keyFormat, req.GetUserId(), req.GetVideoId()),
		reply.GetVideoTime())
	if err != nil {
		return err
	}

	_, err = client.Do("expire",
		fmt.Sprintf(keyFormat, req.GetUserId(), req.GetVideoId()),
		c.cacheExpireTime.Seconds(),
	)
	if err != nil {
		return err
	}
	return nil
}
