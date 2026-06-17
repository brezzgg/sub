package usecase

import (
	"github.com/brezzgg/go-packages/lg"
	"github.com/brezzgg/sub/internal/entity"
	"github.com/brezzgg/sub/internal/pkg/errors"
	"github.com/brezzgg/sub/internal/repo"
)

type Usecase struct {
	rsl *repo.Resolver

	payloadFuncs    []PayloadFunc
	nilPayloadFuncs []PayloadFunc
}

type Option func(u *Usecase) error

type PayloadFunc func(dst *entity.Payload, id string, s *entity.Subscription) error

func WithPayloadFuncs(fn ...PayloadFunc) Option {
	return func(u *Usecase) error {
		u.nilPayloadFuncs = append(u.nilPayloadFuncs, fn...)
		return nil
	}
}

func WithNilPayloadFuncs(fn ...PayloadFunc) Option {
	return func(u *Usecase) error {
		u.nilPayloadFuncs = append(u.nilPayloadFuncs, fn...)
		return nil
	}
}

func NewUsecase(rsl *repo.Resolver, opts ...Option) (*Usecase, error) {
	u := &Usecase{
		rsl: rsl,
	}
	for _, fn := range opts {
		if err := fn(u); err != nil {
			return nil, errors.Directf("option error: %s", err)
		}
	}
	return u, nil
}

func (u *Usecase) Set(id string, s *entity.Subscription) error {
	s.Normalize()
	err := u.rsl.Set(id, s)
	if err != nil {
		if errors.IsUnfixable(err) {
			u.unfixable(id)
		}
		return err
	}
	return nil
}

func (u *Usecase) SetEnabled(id string, enabled bool) error {
	re, err := u.rsl.Get(id)
	if err != nil {
		domain, ok := errors.AsDomain(err)
		if !ok {
			return errors.Error(errors.CodeInternal, err.Error())
		}
		if domain.Code() == errors.CodeUnfixable {
			u.unfixable(id)
		}
		return err
	}
	re.Disabled = !enabled
	re.Normalize()
	err = u.rsl.Set(id, re)
	if err != nil {
		if errors.IsUnfixable(err) {
			u.unfixable(id)
		}
		return err
	}
	return nil
}

func (u *Usecase) Remove(id string) error {
	err := u.rsl.Remove(id)
	if err != nil {
		domain, ok := errors.AsDomain(err)
		if !ok {
			return errors.Error(errors.CodeInternal, err.Error())
		}
		if domain.Code() == errors.CodeNotFound {
			return nil
		}
		return err
	}
	return nil
}

func (u *Usecase) Get(id string) (*entity.Subscription, error) {
	re, err := u.rsl.Get(id)
	if err != nil {
		if errors.IsUnfixable(err) {
			u.unfixable(id)
		}
		return nil, err
	}
	re.Normalize()
	return re, nil
}

func (u *Usecase) GetPayload(id string) (*entity.Payload, error) {
	re, err := u.Get(id)
	if err != nil {
		return nil, err
	}
	if err := re.Ok(); err != nil {
		lg.Warn("not ok", err)
		return nil, err
	}
	if len(re.Payload.Body) == 0 && len(re.Payload.Headers) == 0 {
		for _, fn := range u.nilPayloadFuncs {
			if err := fn(re.Payload, id, re); err != nil {
				return nil, errors.Errorf(errors.CodeInternal, "nil payload func error: %s", err)
			}
		}
	}
	for _, fn := range u.payloadFuncs {
		if err := fn(re.Payload, id, re); err != nil {
			return nil, errors.Errorf(errors.CodeInternal, "payload func error: %s", err)
		}
	}
	re.Payload.Normalize()
	return re.Payload, nil
}

func (u *Usecase) GetAll() (map[string]*entity.Subscription, error) {
	re, resp := u.rsl.WhereFn(func(s *entity.Subscription) bool {
		return true
	})
	if resp != nil {
		return nil, resp
	}
	for _, s := range re {
		s.Normalize()
	}
	return re, nil
}

func (u *Usecase) unfixable(id string) {
	lg.Warn("unfixable", "remove_error", u.rsl.Remove(id))
}
