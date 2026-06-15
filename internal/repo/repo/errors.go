package repo

import "errors"

var (
	ErrNotFound      = errors.New("subscription not found")
	ErrInternal      = errors.New("internal error")
	ErrCacheExpired  = errors.New("cache expired")
	ErrSerialization = errors.New("serealization error")
	ErrAlreadyExists = errors.New("already exists")
)
