package grpc

import (
	"context"
	"errors"

	"github.com/brezzgg/go-packages/lg"
	"github.com/brezzgg/sub/internal/models"
	"github.com/brezzgg/sub/internal/transport/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewServer(repo models.Repo) *grpcServer {
	return &grpcServer{repo: repo}
}

type grpcServer struct {
	pb.UnimplementedSubServiceServer
	repo models.Repo
}

// Get implements [pb.SubServiceServer].
func (g *grpcServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	sub, err := g.repo.Get(req.Id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "subscription with id %s not found", req.Id)
		}
		lg.Error("repository error: %s", err)
		return nil, status.Error(codes.Internal, "internal repository error")
	}
	return &pb.GetResponse{
		Payload:   sub.Payload,
		Ttl:       sub.TTL,
		CreatedAt: sub.CreatedAt,
	}, nil
}

// GetAll implements [pb.SubServiceServer].
func (g *grpcServer) GetAll(ctx context.Context, _ *pb.Empty) (*pb.GetAllResponse, error) {
	all, err := g.repo.GetAll()
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "not found")
		}
		return nil, status.Error(codes.Internal, "internal")
	}
	subs := make(map[string]*pb.GetResponse)
	for id, sub := range all {
		subs[id] = &pb.GetResponse{
			Payload:   sub.Payload,
			Ttl:       sub.TTL,
			CreatedAt: sub.CreatedAt,
		}
	}
	return &pb.GetAllResponse{
		Subscriptions: subs,
	}, nil
}

// Set implements [pb.SubServiceServer].
func (g *grpcServer) Set(ctx context.Context, req *pb.SetRequest) (*pb.Empty, error) {
	sub := models.NewSubscription(req.Ttl, req.Payload)
	if !sub.Validate() {
		return nil, status.Error(codes.InvalidArgument, "invalid subscription format")
	}
	if err := g.repo.Set(req.Id, sub); err != nil {
		lg.Error("repository error", err)
		return nil, status.Errorf(codes.Internal, "internal repository error")
	}
	return nil, nil
}

var _ pb.SubServiceServer = (*grpcServer)(nil)
