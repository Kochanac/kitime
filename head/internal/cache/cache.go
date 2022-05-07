package cache

import (
	"context"
	"fmt"
	pb "github.com/Kochanac/vobla/head/internal/api"
)

func CheckCache(ctx context.Context, req pb.GetRequest) (*pb.GetReply, error) {
	return nil, fmt.Errorf("Not implemented")
}

func SaveToCache(ctx context.Context, req pb.GetRequest, reply pb.GetReply) error {
	return fmt.Errorf("Not implemented")
}
