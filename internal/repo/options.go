package repo

type Options struct {
	CacheProvider, StorageProvider string
	CacheOpts, StorageOpts         map[string]any
}

func (o *Options) Fill() {
	if o.CacheProvider == "" {
		o.CacheProvider = "placeholder"
	}
	if o.StorageProvider == "" {
		o.StorageProvider = "badger"
	}
	if o.CacheOpts == nil {
		o.CacheOpts = map[string]any{}
	}
	if o.StorageOpts == nil {
		o.StorageOpts = map[string]any{}
	}
}
