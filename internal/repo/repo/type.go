package repo

import (
	"fmt"
)

func GetOption[T any](m map[string]any, key string) (*T, error) {
	vala, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("options[%s] cannot be zero", key)
	}
	r, ok := vala.(T)
	if !ok {
		return nil, fmt.Errorf("options[%s] bad value type", key)
	}
	return &r, nil
}
