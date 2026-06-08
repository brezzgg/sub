package grpc

import (
	"context"
	"time"

	"github.com/brezzgg/go-packages/lg"
	"google.golang.org/grpc"
)

func LoggingInterceptor(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (any, error) {
	start := time.Now()
	resp, err := handler(ctx, req)
	lg.Info("grpc request", lg.C{
		"method":   info.FullMethod,
		"duration": time.Since(start).Abs().String(),
		"err":      err,
	})
	return resp, err
}
