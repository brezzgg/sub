package repo

import (
	"github.com/brezzgg/sub/internal/entity"
)

type StorageProvider interface {
	CacheProvider
	Set(id string, s *entity.Subscription) error
	Remove(id string) error
	WhereFn(fn WhereFunc) (map[string]*entity.Subscription, error)
}

type CacheProvider interface {
	Get(id string) (*entity.Subscription, error)
	WhereMeta(equals map[string]any) (map[string]*entity.Subscription, error)
}

// WhereFunc is a function for [StorageLevel.Where]/[CacheLevel.Where] repositories.
type WhereFunc func(s *entity.Subscription) bool
