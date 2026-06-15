package entity_test

import (
	"testing"
	"time"

	"github.com/brezzgg/sub/internal/entity"
)

func TestSubscription_Ok(t *testing.T) {
	tests := []struct {
		name    string
		sub     *entity.Subscription
		wantErr error
	}{
		{
			name: "not expired",
			sub: &entity.Subscription{
				Payload:  &entity.Payload{},
				Expired:  time.Now().Add(time.Second),
				Disabled: false,
			},
			wantErr: nil,
		},
		{
			name: "expired",
			sub: &entity.Subscription{
				Payload:  &entity.Payload{},
				Expired:  time.Now().Add(-time.Second),
				Disabled: false,
			},
			wantErr: entity.ErrSubExpired,
		},
		{
			name: "second at second",
			sub: &entity.Subscription{
				Payload:  &entity.Payload{},
				Expired:  time.Now(),
				Disabled: false,
			},
			wantErr: nil,
		},
		{
			name: "not expired but disabled",
			sub: &entity.Subscription{
				Payload:  &entity.Payload{},
				Expired:  time.Now().Add(time.Second),
				Disabled: true,
			},
			wantErr: entity.ErrSubDisabled,
		},
		{
			name: "restore payload",
			sub: &entity.Subscription{
				Payload:  nil,
				Expired:  time.Now().Add(time.Second),
				Disabled: false,
			},
			wantErr: nil,
		},
		{
			name: "unlimited",
			sub: &entity.Subscription{
				Payload:  &entity.Payload{},
				Expired:  time.Time{},
				Disabled: false,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.sub.Ok()
			if gotErr != tt.wantErr {
				t.Errorf("Ok() = %s, want = %s", gotErr, tt.wantErr)
			}
		})
	}
}
