package grpc

import (
	"context"

	"github.com/brezzgg/sub/internal/entity"
	"github.com/brezzgg/sub/internal/pkg/errors"
	"github.com/brezzgg/sub/internal/usecase"
)

func NewService(usec *usecase.Usecase) *Service {
	return &Service{usec: usec}
}

type Service struct {
	UnimplementedSubServiceServer
	usec *usecase.Usecase
}

// Get implements [SubServiceServer].
func (s *Service) Get(ctx context.Context, req *IdRequest) (*usecase.SubscriptionRawPb, error) {
	sub, err := s.usec.Get(req.Id)
	if err != nil {
		return nil, errors.AddContext(err, "repository error")
	}
	raw := &usecase.SubscriptionRawPb{}
	if err := raw.FromSubscription(req.Id, sub); err != nil {
		return nil, errors.Directf("decode subscription: %s", err)
	}
	return raw, nil
}

// GetAll implements [SubServiceServer].
func (s *Service) GetAll(context.Context, *Empty) (*GetAllResponse, error) {
	all, err := s.usec.GetAll()
	if err != nil {
		return nil, errors.AddContext(err, "get all")
	}
	subs := &GetAllResponse{
		Subs: make([]*usecase.SubscriptionRawPb, 0, len(all)),
	}
	for id, sub := range all {
		raw := &usecase.SubscriptionRawPb{}
		if err := raw.FromSubscription(id, sub); err != nil {
			return nil, errors.Directf("decode subscription: %s", err)
		}
		subs.Subs = append(subs.Subs, raw)
	}
	return subs, nil
}

// Remove implements [SubServiceServer].
func (s *Service) Remove(ctx context.Context, req *IdRequest) (*Empty, error) {
	if err := s.usec.Remove(req.Id); err != nil {
		return nil, err
	}
	return &Empty{}, nil
}

// Set implements [SubServiceServer].
func (s *Service) Set(ctx context.Context, req *SetRequest) (*Empty, error) {
	var (
		sb  *entity.Subscription
		id  string
		err error
	)
	raw := req.GetRaw()
	if raw != nil {
		id, sb, err = raw.ToSubscription()
		if err != nil {
			return nil, errors.Error(errors.CodeInvalidArgument, err.Error())
		}
	} else {
		re := req.GetReq()
		id, sb, err = re.ToSubscription()
		if err != nil {
			return nil, errors.Error(errors.CodeInvalidArgument, err.Error())
		}
	}
	if err := s.usec.Set(id, sb); err != nil {
		return nil, err
	}
	return nil, nil
}

// SetEnabled implements [SubServiceServer].
func (s *Service) SetEnabled(ctx context.Context, req *SetEnabledRequest) (*Empty, error) {
	err := s.usec.SetEnabled(req.GetId().GetId(), req.GetEnabled())
	if err != nil {
		return nil, err
	}
	return nil, nil
}

var _ SubServiceServer = (*Service)(nil)
