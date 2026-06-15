package usecase_test

import (
	"encoding/base64"
	"reflect"
	"testing"
	"time"

	"github.com/brezzgg/sub/internal/entity"
	"github.com/brezzgg/sub/internal/usecase"
)

func TestSetRawRequest_ToSubscription(t *testing.T) {
	tests := []struct {
		name            string
		request         *usecase.SubscriptionRawPb
		wantId          string
		wantSub         *entity.Subscription
		wantErr, wantOk bool
	}{
		{
			name:   "headers clean",
			wantOk: true,
			request: &usecase.SubscriptionRawPb{
				Id:                "someid",
				PayloadBodyString: "somebody",
				PayloadHeaders: map[string]string{
					"header1": "value1",
					"header2": "",
				},
				Expired:  0,
				Disabled: false,
				Metadata: nil,
			},
			wantId: "someid",
			wantSub: &entity.Subscription{
				Payload: &entity.Payload{
					Body: []byte("c29tZWJvZHk="),
					Headers: map[string]string{
						"header1": "value1",
					},
				},
				Expired:  time.Time{},
				Disabled: false,
				Metadata: map[string]any{},
			},
			wantErr: false,
		},
		{
			name:   "empty body",
			wantOk: false,
			request: &usecase.SubscriptionRawPb{
				Id:                "someid",
				PayloadBodyString: "",
				PayloadHeaders: map[string]string{
					"header1": "value1",
					"header2": "",
				},
				Expired:  0,
				Disabled: true,
				Metadata: nil,
			},
			wantId: "someid",
			wantSub: &entity.Subscription{
				Payload: &entity.Payload{
					Body: []byte{},
					Headers: map[string]string{
						"header1": "value1",
					},
				},
				Expired:  time.Time{},
				Disabled: true,
				Metadata: map[string]any{},
			},
			wantErr: false,
		},
		{
			name:   "expired",
			wantOk: false,
			request: &usecase.SubscriptionRawPb{
				Id:                "someid",
				PayloadBodyString: "",
				PayloadHeaders: map[string]string{
					"header1": "value1",
					"header2": "",
				},
				Expired:  1009843200,
				Disabled: false,
				Metadata: nil,
			},
			wantId: "someid",
			wantSub: &entity.Subscription{
				Payload: &entity.Payload{
					Body: []byte{},
					Headers: map[string]string{
						"header1": "value1",
					},
				},
				Expired:  (time.Time{}).AddDate(2001, 0, 0),
				Disabled: false,
				Metadata: map[string]any{},
			},
			wantErr: false,
		},
		{
			name:   "all headers empty values",
			wantOk: true,
			request: &usecase.SubscriptionRawPb{
				Id:                "someid",
				PayloadBodyString: "",
				PayloadHeaders: map[string]string{
					"header1": "",
					"header2": "",
				},
				Expired:  0,
				Disabled: false,
				Metadata: nil,
			},
			wantId: "someid",
			wantSub: &entity.Subscription{
				Payload: &entity.Payload{
					Body:    []byte{},
					Headers: map[string]string{},
				},
				Expired:  time.Time{},
				Disabled: false,
				Metadata: map[string]any{},
			},
			wantErr: false,
		},
		{
			name:   "nil headers",
			wantOk: true,
			request: &usecase.SubscriptionRawPb{
				Id:                "someid",
				PayloadBodyString: "",
				PayloadHeaders:    nil,
				Expired:           0,
				Disabled:          false,
				Metadata:          nil,
			},
			wantId: "someid",
			wantSub: &entity.Subscription{
				Payload: &entity.Payload{
					Body:    []byte{},
					Headers: map[string]string{},
				},
				Expired:  time.Time{},
				Disabled: false,
				Metadata: map[string]any{},
			},
			wantErr: false,
		},
		{
			name:   "expired and disabled",
			wantOk: false,
			request: &usecase.SubscriptionRawPb{
				Id:                "someid",
				PayloadBodyString: "somebody",
				PayloadHeaders: map[string]string{
					"header1": "value1",
				},
				Expired:  1009843200,
				Disabled: true,
				Metadata: nil,
			},
			wantId: "someid",
			wantSub: &entity.Subscription{
				Payload: &entity.Payload{
					Body: []byte("c29tZWJvZHk="),
					Headers: map[string]string{
						"header1": "value1",
					},
				},
				Expired:  (time.Time{}).AddDate(2001, 0, 0),
				Disabled: true,
				Metadata: map[string]any{},
			},
			wantErr: false,
		},
		{
			name:   "empty id",
			wantOk: false,
			request: &usecase.SubscriptionRawPb{
				Id:                "",
				PayloadBodyString: "somebody",
				PayloadHeaders: map[string]string{
					"header1": "value1",
				},
				Expired:  0,
				Disabled: false,
				Metadata: nil,
			},
			wantId: "",
			wantSub: &entity.Subscription{
				Payload: &entity.Payload{
					Body: []byte("c29tZWJvZHk="),
					Headers: map[string]string{
						"header1": "value1",
					},
				},
				Expired:  time.Time{},
				Disabled: false,
				Metadata: map[string]any{},
			},
			wantErr: true,
		},
		{
			name:   "unicode body",
			wantOk: true,
			request: &usecase.SubscriptionRawPb{
				Id:                "someid",
				PayloadBodyString: "привет мир",
				PayloadHeaders: map[string]string{
					"content-type": "text/plain; charset=utf-8",
				},
				Expired:  0,
				Disabled: false,
				Metadata: nil,
			},
			wantId: "someid",
			wantSub: &entity.Subscription{
				Payload: &entity.Payload{
					Body: []byte(base64.StdEncoding.EncodeToString([]byte("привет мир"))),
					Headers: map[string]string{
						"content-type": "text/plain; charset=utf-8",
					},
				},
				Expired:  time.Time{},
				Disabled: false,
				Metadata: map[string]any{},
			},
			wantErr: false,
		},
		{
			name:   "subscription test",
			wantOk: true,
			request: &usecase.SubscriptionRawPb{
				Id:                "someid",
				PayloadBodyString: "\nsub1://somebody1#some name 1\nsub2://somebody2#some name 2\n",
				PayloadHeaders: map[string]string{
					"header1": "value1",
				},
				Expired:  0,
				Disabled: false,
				Metadata: nil,
			},
			wantId: "someid",
			wantSub: &entity.Subscription{
				Payload: &entity.Payload{
					Body: []byte(base64.StdEncoding.EncodeToString([]byte("sub1://somebody1#some name 1\nsub2://somebody2#some name 2"))),
					Headers: map[string]string{
						"header1": "value1",
					},
				},
				Expired:  time.Time{},
				Disabled: false,
				Metadata: map[string]any{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotId, gotSub, gotErr := tt.request.ToSubscription()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("ToSubscription() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ToSubscription() succeeded unexpectedly")
			}
			if tt.wantId != gotId {
				t.Errorf("ToSubscription()[id] = %v, want %v", gotId, tt.wantId)
			}
			if !reflect.DeepEqual(tt.wantSub.Payload.Body, gotSub.Payload.Body) {
				t.Errorf("ToSubscription().Payload.Body = %v, want %v", gotSub.Payload.Body, tt.wantSub.Payload.Body)
			}
			if !reflect.DeepEqual(tt.wantSub.Payload.Headers, gotSub.Payload.Headers) {
				t.Errorf("ToSubscription().Payload.Headers = %v, want %v", gotSub.Payload.Headers, tt.wantSub.Payload.Headers)
			}
			tt.wantSub.Payload = nil
			gotSub.Payload = nil
			if !reflect.DeepEqual(tt.wantSub, gotSub) {
				t.Errorf("ToSubscription() = %v, want %v", gotSub, tt.wantSub)
			}
			if err := gotSub.Ok(); (err == nil) != tt.wantOk {
				t.Errorf("ToSubscription().Ok() = %v, want %v", (err == nil), tt.wantOk)
			}
		})
	}
}
