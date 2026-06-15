package manager

import (
	"errors"
	"fmt"
	"net"

	"github.com/brezzgg/go-packages/lg"
	"github.com/brezzgg/sub/internal/manager/scheduler"
	"github.com/brezzgg/sub/internal/repo"
	tgrpc "github.com/brezzgg/sub/internal/transport/grpc"
	"github.com/brezzgg/sub/internal/transport/http"
	"github.com/brezzgg/sub/internal/usecase"
	"google.golang.org/grpc"
)

type Manager struct {
	rsl   *repo.Resolver
	usec  *usecase.Usecase
	httph *http.Handler
	grpcs *tgrpc.Service
	grpcc tgrpc.SubServiceClient
	sch   *scheduler.Scheduler
}

type Option func(m *Manager) error

func WithGrpcClient(remote string) Option {
	return func(m *Manager) error {
		var err error
		m.grpcc, err = tgrpc.GetClient(remote)
		if err != nil {
			return fmt.Errorf("failed to init grpc client: %s", err)
		}
		return nil
	}
}

func WithHttpHandler(host, pattern string) Option {
	return func(m *Manager) error {
		if m.usec == nil {
			return errors.New("usecase cannot be nil")
		}
		m.httph = http.New(m.usec, host, pattern)
		m.sch.Add(m.httph.Run, m.httph.Stop)
		lg.Info("http server configured")
		return nil
	}
}

func WithGrpcService(host string) Option {
	return func(m *Manager) error {
		if m.usec == nil {
			return errors.New("usecase cannot be nil")
		}
		lis, err := net.Listen("tcp", host)
		if err != nil {
			return fmt.Errorf("failed to listen '%s': %s", host, err)
		}
		srv := grpc.NewServer(grpc.UnaryInterceptor(tgrpc.LoggingInterceptor))
		tgrpc.RegisterSubServiceServer(srv, tgrpc.NewService(m.usec))
		m.sch.Add(func() error {
			return srv.Serve(lis)
		}, func() error {
			srv.GracefulStop()
			return nil
		})
		lg.Info("grpc server configured")
		return nil
	}
}

func NewClientManager(opts ...Option) (*Manager, error) {
	m := Manager{}

	// apply options
	for i, optFn := range opts {
		if err := optFn(&m); err != nil {
			return nil, fmt.Errorf("opts[%d] failed to apply: %s", i, err)
		}
	}

	return &m, nil
}

func NewManager(repoOpt *repo.Options, usecOpt []usecase.Option, opts ...Option) (*Manager, error) {
	m := Manager{}

	// repo resolver
	var err error
	m.rsl, err = repo.NewResolver(repoOpt)
	if err != nil {
		return nil, fmt.Errorf("failed to init repository: %s", err)
	}

	// usecase
	m.usec, err = usecase.NewUsecase(m.rsl, usecOpt...)
	if err != nil {
		return nil, fmt.Errorf("failed to init usecase: %s", err)
	}

	// scheduler
	m.sch = &scheduler.Scheduler{}

	// apply options
	for i, optFn := range opts {
		if err := optFn(&m); err != nil {
			return nil, fmt.Errorf("opts[%d] failed to apply: %s", i, err)
		}
	}

	return &m, nil
}

func (m *Manager) Usecase() *usecase.Usecase {
	return m.usec
}

func (m *Manager) GrpcClient() tgrpc.SubServiceClient {
	return m.grpcc
}

func (m *Manager) Run() error {
	return m.sch.Run()
}

func (m *Manager) Stop() {
	m.sch.Stop()
}
