package repo

import (
	"github.com/brezzgg/sub/internal/repo/cache/placeholder"
	"github.com/brezzgg/sub/internal/repo/repo"
	"github.com/brezzgg/sub/internal/repo/storage/badger"
)

func init() {
	// cache
	repo.CacheRegistry = map[string]repo.CacheProviderConstructor{
		"placeholder": placeholder.Constructor,
	}
	// storage
	repo.StorageRegistry = map[string]repo.StorageProviderConstructor{
		"badger": badger.Constructor,
	}
}
