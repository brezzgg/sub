package placeholder

import (
	"github.com/brezzgg/sub/internal/entity"
	"github.com/brezzgg/sub/internal/pkg/errors"
	"github.com/brezzgg/sub/internal/repo/repo"
)

// Placeholder is a implemetaion of [repo.CacheProvider]
// that used when cache provider dont needed.
type Placeholder struct{}

func Constructor(options map[string]any) (repo.CacheProvider, error) {
	return &Placeholder{}, nil
}

// Get implements [repo.CacheProvider].
func (p *Placeholder) Get(id string) (*entity.Subscription, error) {
	return nil, errors.Code(errors.CodeCacheExpired)
}

// WhereMeta implements [repo.CacheProvider].
func (p *Placeholder) WhereMeta(equals map[string]any) (map[string]*entity.Subscription, error) {
	return nil, errors.Code(errors.CodeCacheExpired)
}

var _ repo.CacheProvider = (*Placeholder)(nil)
