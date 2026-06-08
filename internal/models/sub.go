package models

import "time"

// NewSubscription creates a new subscription.
// TTL will be set to default value if it is not valid.
func NewSubscription(ttl int64, payload []byte) *Subscription {
	if ttl < 0 {
		ttl = DefaultTTL
	}
	return &Subscription{
		TTL:       ttl,
		CreatedAt: time.Now().Unix(),
		Payload:   payload,
	}
}

// Validate validates ttl and payload.
func (s *Subscription) Validate() bool {
	if len(s.Payload) < 1 {
		return false
	}

	// created at not set
	if s.CreatedAt == 0 {
		return false
	}

	// if ttl < 1 then ttl is unlimited
	if s.TTL < 1 {
		return true
	}

	// ttl check
	if (s.CreatedAt + s.TTL) > time.Now().Unix() {
		return true
	} else {
		return false
	}
}
