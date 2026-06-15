package entity

import (
	"errors"
	"time"
)

// Subscription is a subscription data struct.
type Subscription struct {
	// Payload is a subscription body and headers.
	Payload *Payload
	// Expired indicate when subsciption be expired.
	// Suscription unlimited if default type value.
	Expired time.Time
	// Disabled indicate status of subscription.
	// If disable http server returns empty payload body.
	Disabled bool
	// Metadata is a key-value storage that
	// can contain any subscription metadata.
	Metadata map[string]any
}

// Ok returns error if subscription cannot be send to user.
func (s *Subscription) Ok() error {
	s.Normalize()

	if s.Disabled {
		return ErrSubDisabled
	}

	if s.Expired.Year() >= 2000 {
		// check expired
		if time.Now().Unix() > s.Expired.Unix() {
			return ErrSubExpired
		}
	} else {
		// set zero
		s.Expired = time.Time{}
	}

	if s.Payload == nil {
		return ErrSubBroken
	}

	return s.Payload.Ok()
}

// Normalize normalize subscription fields.
func (s *Subscription) Normalize() {
	if s == nil {
		return
	}

	// expired
	if s.Expired.Year() < 2000 {
		s.Expired = time.Time{}
	}

	// payload
	if s.Payload == nil {
		s.Payload = &Payload{}
	}
	s.Payload.Normalize()

	// metadata
	if s.Metadata == nil {
		s.Metadata = map[string]any{}
	}
}

func (s *Subscription) GetMeta(key string) (any, bool) {
	if s.Metadata == nil {
		return nil, false
	}
	val, ok := s.Metadata[key]
	return val, ok
}

func (s *Subscription) SetMeta(key string, val any) {
	if s.Metadata == nil {
		s.Metadata = make(map[string]any, 1)
	}
	s.Metadata[key] = val
}

// Default implements [Defaultable].
func (s *Subscription) Default() *Subscription {
	return &Subscription{
		Payload:  (&Payload{}).Default(),
		Expired:  time.Time{},
		Disabled: false,
		Metadata: map[string]any{},
	}
}

var (
	ErrSubDisabled = errors.New("subscription manualy disabled")
	ErrSubExpired  = errors.New("subscription expired")
	ErrSubBroken   = errors.New("subscription integrity is broken")
)

var _ Defaultable[Subscription] = (*Subscription)(nil)
