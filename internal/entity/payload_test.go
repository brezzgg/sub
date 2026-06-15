package entity_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/brezzgg/sub/internal/entity"
)

func TestPayload_FormatBody(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			name: "oneline",
			s:    "vless://some-uuid@192.168.0.1:443/?key1=value&key2=value2\n",
			want: "vless://some-uuid@192.168.0.1:443/?key1=value&key2=value2",
		},
		{
			name: "space in the middle",
			s:    "vless://some-uuid@192.168.0.1:443/ key1=value&key2=value2\n",
			want: "vless://some-uuid@192.168.0.1:443/ key1=value&key2=value2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entity.Payload{}
			got := p.FormatBody(tt.s)
			if got != tt.want {
				t.Errorf("FormatBody() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPayload_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		b       string
		wantErr bool
		want    string
	}{
		{
			name:    "ok",
			b:       "somebody",
			wantErr: false,
		},
		{
			name:    "empty string",
			b:       "",
			wantErr: false,
		},
		{
			name:    "with space symbols",
			b:       "some\n pay \n\t load",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &entity.Payload{}
			gotErr := p.MarshalBody(tt.b)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MarshalBody() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MarshalBody() succeeded unexpectedly")
			}
			unmarshRes, gotErr := p.UnmarshalBody()
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("UnmarshalBody() failed: %v", gotErr)
				}
			}
			if p.FormatBody(tt.b) != unmarshRes {
				t.Errorf("Marshaled != Unmarshal: '%s' != '%s'", tt.b, unmarshRes)
			}
		})
	}
}

func TestPayload_Ok(t *testing.T) {
	tests := []struct {
		name     string
		pl, want *entity.Payload
		err      error
	}{
		{
			name: "header deletion",
			pl: &entity.Payload{
				Body: []byte{},
				Headers: map[string]string{
					"header1": "value1",
					"header2": "",
				},
			},
			want: &entity.Payload{
				Body: []byte{},
				Headers: map[string]string{
					"header1": "value1",
				},
			},
			err: nil,
		},
		{
			name: "header deletion",
			pl: &entity.Payload{
				Body: []byte{},
				Headers: map[string]string{
					"header1": "value1",
					"header2": "",
				},
			},
			want: &entity.Payload{
				Body: []byte{},
				Headers: map[string]string{
					"header1": "value1",
				},
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr := tt.pl.Ok()
			if gotErr != nil {
				if errors.Is(tt.err, gotErr) {
					return
				}
				t.Errorf("Ok() = %s, want %s", gotErr, tt.err)
			}
			if !reflect.DeepEqual(tt.want.Headers, tt.pl.Headers) {
				t.Errorf("Payload.Headers = %v, want %v", tt.pl.Headers, tt.want.Headers)
			}
		})
	}
}
