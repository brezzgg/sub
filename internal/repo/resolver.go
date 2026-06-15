package repo

import (
	"errors"
	"fmt"

	"github.com/brezzgg/sub/internal/entity"
	R "github.com/brezzgg/sub/internal/repo/repo"
)

type Resolver struct {
	cimpl R.CacheProvider
	simpl R.StorageProvider
}

func NewResolver(opt *Options) (*Resolver, error) {
	opt.Fill()

	rsl := &Resolver{}

	// init cache
	if prFn, ok := R.CacheRegistry[opt.CacheProvider]; !ok {
		return nil, errors.New("cache provider not found")
	} else {
		pr, err := prFn(opt.CacheOpts)
		if err != nil {
			return nil, fmt.Errorf("cache provider error: %s", err)
		}
		rsl.cimpl = pr
	}

	// init storage
	if prFn, ok := R.StorageRegistry[opt.StorageProvider]; !ok {
		return nil, errors.New("storage provider not found")
	} else {
		pr, err := prFn(opt.StorageOpts)
		if err != nil {
			return nil, fmt.Errorf("storage provider error: %s", err)
		}
		rsl.simpl = pr
	}

	return rsl, nil
}

// Get implements [StorageProvider].
// Cache available.
func (r *Resolver) Get(id string) (re *entity.Subscription, err error) {
	re, err = r.cimpl.Get(id)
	if err == nil {
		return
	} else {
		re, err = r.simpl.Get(id)
	}
	return
}

// WhereMeta implements [StorageProvider].
// Cache available.
func (r *Resolver) WhereMeta(equals map[string]any) (re map[string]*entity.Subscription, err error) {
	re, err = r.cimpl.WhereMeta(equals)
	if err == nil {
		return
	} else {
		re, err = r.simpl.WhereMeta(equals)
	}
	return
}

// WhereFn implements [StorageProvider].
func (r *Resolver) WhereFn(fn R.WhereFunc) (re map[string]*entity.Subscription, err error) {
	return r.simpl.WhereFn(fn)
}

// Set implements [StorageProvider].
func (r *Resolver) Set(id string, s *entity.Subscription) error {
	return r.simpl.Set(id, s)
}

// Remove implements [StorageProvider].
func (r *Resolver) Remove(id string) error {
	return r.simpl.Remove(id)
}

var _ R.StorageProvider = (*Resolver)(nil)
