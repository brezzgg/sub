package models

// Subscription is a subscription data struct.
type Subscription struct {
	// Payload is a subscription body and headers in protocol buffers format.
	Payload []byte
	// TTL is a time to live of subscription in unix timestamp format.
	TTL int64
	// CreatedAt is a creation date in unix timestamp format.
	CreatedAt int64
}

// DefaultTTL is a default value of ttl if user ttl not valid.
var DefaultTTL int64 = 0

// Repo is a repository of subscriptions.
type Repo interface {
	// Get returns subscription by id.
	Get(id string) (*Subscription, error)
	// GetAll returns all subscriptions.
	GetAll() (map[string]*Subscription, error)
	// GetFunc returns all subscriptions where fn returns true.
	GetFunc(fn func(s *Subscription) bool) ([]*Subscription, error)
	// Set saves or update subscription.
	Set(id string, s *Subscription) error
	// Remove deletes a subscription.
	Remove(id string) error
}
