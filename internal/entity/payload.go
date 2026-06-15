package entity

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// Payload is a struct that contains a
// subscription response for http protocol.
type Payload struct {
	Body    []byte
	Headers map[string]string
}

// Ok returns error if payload cannot be send to user.
func (p *Payload) Ok() error {
	var del []string
	for key, val := range p.Headers {
		if strings.TrimSpace(val) == "" {
			del = append(del, key)
		}
	}
	for _, d := range del {
		delete(p.Headers, d)
	}

	return nil
}

// Normalize normalize payload fields.
func (p *Payload) Normalize() {
	if p.Headers == nil {
		p.Headers = map[string]string{}
	}
	var del []string
	for k, v := range p.Headers {
		if strings.TrimSpace(v) == "" {
			del = append(del, k)
		}
	}
	for _, d := range del {
		delete(p.Headers, d)
	}

	if p.Body == nil {
		p.Body = []byte{}
	}
}

// MarshalBody encodes string into base64 bytes [Payload.Body].
func (p *Payload) MarshalBody(b string) error {
	if p == nil {
		return fmt.Errorf("payload is nil")
	}
	b = p.FormatBody(b)
	if b == "" {
		p.Body = nil
		return nil
	}
	encoded := base64.StdEncoding.EncodeToString([]byte(b))
	p.Body = []byte(encoded)
	return nil
}

// UnmarshalBody decodes [Payload.Body] to string.
func (p *Payload) UnmarshalBody() (string, error) {
	if p == nil {
		return "", fmt.Errorf("payload is nil")
	}
	res, err := base64.StdEncoding.DecodeString(string(p.Body))
	if err != nil {
		return "", fmt.Errorf("base64 decode: %s", err)
	}
	return string(res), nil
}

// FormatBody formats string to subscription format.
// Any white space replaces to \n symbol.
func (p *Payload) FormatBody(s string) string {
	s = strings.TrimSpace(s)
	split := strings.Split(s, "\n")
	for i, sp := range split {
		split[i] = strings.TrimSpace(sp)
	}
	s = strings.Join(split, "\n")
	return s
}

// Default implements [Defaultable].
func (p *Payload) Default() *Payload {
	return &Payload{
		Body:    []byte{},
		Headers: make(map[string]string, 5),
	}
}

var _ Defaultable[Payload] = (*Payload)(nil)
