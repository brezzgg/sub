package models

import "errors"

var (
	ErrNotFound = errors.New("subscription not found")
	ErrExpired  = errors.New("subscription expired")
	ErrInternal = errors.New("internal error")
)
