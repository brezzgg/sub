package repo

type StorageProviderConstructor func(options map[string]any) (StorageProvider, error)
type CacheProviderConstructor func(options map[string]any) (CacheProvider, error)

var (
	StorageRegistry map[string]StorageProviderConstructor
	CacheRegistry   map[string]CacheProviderConstructor
)
