package subset

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type Model struct {
	Id               string            `yaml:"id"`
	Title            string            `yaml:"title"`
	UpdateInterval   string            `yaml:"update_interval"`
	UserWebPage      string            `yaml:"user_web_page"`
	MovedPermanently string            `yaml:"moved_permanently"`
	SupportUri       string            `yaml:"support_uri"`
	PreferredDNS     string            `yaml:"preferred_dns"`
	CustomHeaders    map[string]string `yaml:"custom_headers"`
	TTL              int               `yaml:"ttl"`
	Bodies           []string          `yaml:"bodies"`
}

// UnmarshalYAML implements [yaml.Unmarshaler].
func (m *Model) UnmarshalYAML(value *yaml.Node) error {
	type plain Model
	var tmp plain
	if err := value.Decode(&tmp); err != nil {
		return err
	}

	if len(tmp.Id) < 2 {
		return fmt.Errorf("id: length must be >= 2")
	}
	tmp.Id = strings.TrimSpace(tmp.Id)
	for i, c := range tmp.Id {
		if !isURISafe(c) {
			return fmt.Errorf("id: invalid character %q at position %d", c, i)
		}
	}

	if tmp.Title == "" {
		tmp.Title = "Subscription"
	}

	if tmp.UpdateInterval != "" {
		n, err := strconv.Atoi(tmp.UpdateInterval)
		if err != nil || n <= 0 {
			return fmt.Errorf("update_interval: must be a positive integer, got %q", tmp.UpdateInterval)
		}
	}

	if tmp.TTL < 0 {
		return fmt.Errorf("ttl: must be a positive integer")
	}

	for i := range tmp.Bodies {
		tmp.Bodies[i] = strings.TrimSpace(tmp.Bodies[i])
		if tmp.Bodies[i] == "" {
			return fmt.Errorf("bodies[%d]: must not be empty", i)
		}
	}

	*m = Model(tmp)
	return nil
}

func isURISafe(c rune) bool {
	return (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') ||
		(c >= '0' && c <= '9') ||
		c == '-' || c == '_' || c == '.' || c == '~'
}

var _ yaml.Unmarshaler = (*Model)(nil)
