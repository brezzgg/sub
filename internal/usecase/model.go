package usecase

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/brezzgg/sub/internal/entity"
	"github.com/brezzgg/sub/internal/pkg/yamlreflect"
	"gopkg.in/yaml.v3"
)

// ToSubscription converts pb format to [entity.Subscription].
func (r *SubscriptionRawPb) ToSubscription() (string, *entity.Subscription, error) {
	sub := &entity.Subscription{}
	sub.Normalize()

	id := strings.TrimSpace(r.GetId())
	if len([]rune(id)) < 3 {
		return "", nil, fmt.Errorf("id must be >= 3")
	}

	// body
	sub.Payload = &entity.Payload{}
	if err := sub.Payload.MarshalBody(r.GetPayloadBodyString()); err != nil {
		return "", nil, fmt.Errorf("marshal raw body: %s", err)
	}

	// headers
	sub.Payload.Headers = r.GetPayloadHeaders()
	if sub.Payload.Headers == nil {
		sub.Payload.Headers = map[string]string{}
	}

	// expired
	sub.Expired = time.Unix(r.GetExpired(), 0).UTC()
	if sub.Expired.Year() < 2000 {
		sub.Expired = time.Time{}
	}

	// disabled
	sub.Disabled = r.GetDisabled()

	// metadata
	meta, err := r.GetMetadata().Unmarshal()
	if err != nil {
		return "", nil, err
	}
	sub.Metadata = meta

	// normalize
	sub.Normalize()

	return id, sub, nil
}

// FromSubscription converts [entity.Subscription] to pb format.
func (r *SubscriptionRawPb) FromSubscription(id string, s *entity.Subscription) error {
	s.Normalize()

	// id
	r.Id = strings.TrimSpace(id)
	if len([]rune(r.Id)) < 3 {
		return fmt.Errorf("if length must be >= 3")
	}

	// payload
	if len(s.Payload.Body) > 0 {
		str, err := s.Payload.UnmarshalBody()
		if err != nil {
			return fmt.Errorf("failed to unmarshal body: %s", err)
		}
		r.PayloadBodyString = str
	}
	r.PayloadHeaders = s.Payload.Headers

	// expired
	if s.Expired.UTC().Year() < 2000 {
		r.Expired = 0
	} else {
		r.Expired = s.Expired.UTC().Unix()
	}

	// disabled
	r.Disabled = s.Disabled

	// meta
	r.Metadata = &MetadataRaw{}
	if err := r.Metadata.Marshal(s.Metadata); err != nil {
		return err
	}

	return nil
}

// ToSubscription converts request struct to [entity.Subscription].
func (r *SetRequest) ToSubscription() (string, *entity.Subscription, error) {
	id, sub, err := r.Raw.ToSubscription()
	if err != nil {
		return "", nil, err
	}

	// custom headers
	for k, v := range r.GetCustomHeaders() {
		sub.Payload.Headers[k] = v
	}

	// title
	if strings.TrimSpace(r.Title) == "" {
		r.Title = "Subscription"
	}
	sub.Payload.Headers["Profile-Title"] = r.Title

	// update interval
	interval, err := strconv.Atoi(r.UpdateInterval)
	if err != nil {
		return "", nil, fmt.Errorf("update_interval must be integer")
	}
	if interval < 0 {
		interval = 6
	}
	sub.Payload.Headers["profile-update-interval"] = strconv.Itoa(interval)

	// user web page
	userWeb := ""
	if strings.TrimSpace(r.GetUserWebPage()) != "" {
		url, err := url.Parse(r.GetUserWebPage())
		if err != nil {
			return "", nil, fmt.Errorf("user_web_page bad format: %s", err)
		}
		userWeb = url.String()
	}
	sub.Payload.Headers["profile-web-page-url"] = userWeb

	// moved permanently
	moved := ""
	if strings.TrimSpace(r.GetMovedPermanently()) != "" {
		url, err := url.Parse(r.GetMovedPermanently())
		if err != nil {
			return "", nil, fmt.Errorf("moved_permanently bad format: %s", err)
		}
		moved = url.String()
	}
	sub.Payload.Headers["moved-permanently-to"] = moved

	// support uri
	support := ""
	if strings.TrimSpace(r.GetSupportUri()) != "" {
		url, err := url.Parse(r.GetSupportUri())
		if err != nil {
			return "", nil, fmt.Errorf("support_uri bad format: %s", err)
		}
		support = url.String()
	}
	sub.Payload.Headers["support-url"] = support

	// dns
	sub.Payload.Headers["DNS"] = r.GetDns()

	// user info
	userInfoRes := ""
	userInfoAddRes := func(k string, v int64) {
		userInfoRes += fmt.Sprintf(" %s=%d;", k, v)
	}
	userInfoIntExp := false
	if sub.Expired.Year() > 2000 {
		userInfoAddRes("expire", sub.Expired.Unix())
		userInfoIntExp = true
	}
	if ui := r.GetUserInfo(); ui != nil {
		if i := ui.GetDownload(); i != 0 {
			userInfoAddRes("download", i)
		}
		if i := ui.GetUpload(); i != 0 {
			userInfoAddRes("upload", i)
		}
		if i := ui.GetTotal(); i != 0 {
			userInfoAddRes("total", i)
		}
		if i := ui.GetExpired(); i != 0 && !userInfoIntExp {
			userInfoAddRes("expire", i)
		}
	}
	if userInfoRes != "" {
		sub.Payload.Headers["Subscription-Userinfo"] = strings.Trim(userInfoRes, " ;")
	}

	// bodies
	if len(sub.Payload.Body) > 0 && len(r.GetBodies()) > 0 {
		return "", nil, fmt.Errorf("body is contained in either payload_body_string or bodies, use one of the two")
	}
	if len(r.Bodies) > 0 {
		if err := sub.Payload.MarshalBody(strings.Join(r.GetBodies(), "\n")); err != nil {
			return "", nil, fmt.Errorf("bodies failed to marshal: %s", err)
		}
	}

	// normalize
	sub.Normalize()

	return id, sub, nil
}

// Unmarshal umarshals metadata json format to map[string]any.
func (m *MetadataRaw) Unmarshal() (map[string]any, error) {
	if len(m.GetData()) == 0 {
		return map[string]any{}, nil
	}
	var res map[string]any
	if err := json.Unmarshal(m.GetData(), &res); err != nil {
		return nil, fmt.Errorf("unmarshal metadata: %s", err)
	}
	return res, nil
}

// Marshal marshals map[string]any to metadata json format.
func (m *MetadataRaw) Marshal(meta map[string]any) error {
	if b, err := json.Marshal(meta); err != nil {
		return fmt.Errorf("marhsal metadata: %s", err)
	} else {
		m.Data = b
		return nil
	}
}

// MarshalYAML implements [yaml.Marshaler]
func (r *SetRequest) MarshalYAML() (any, error) {
	return yamlreflect.CopyToYAMLStruct(r), nil
}

// UnmarshalYAML implements [yaml.Unmarshaler]
func (r *SetRequest) UnmarshalYAML(value *yaml.Node) error {
	yamlType := yamlreflect.ConvertToYAMLStruct(r)
	yamlDst := reflect.New(yamlType)

	if err := value.Decode(yamlDst.Interface()); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	yamlreflect.CopyFromYAMLStruct(r, yamlDst.Interface())
	return nil
}

// UnmarshalYAML implements [yaml.Unmarshaler]
func (r *SubscriptionRawPb) UnmarshalYAML(value *yaml.Node) error {
	yamlType := yamlreflect.ConvertToYAMLStruct(r)
	yamlDst := reflect.New(yamlType)

	if err := value.Decode(yamlDst.Interface()); err != nil {
		return fmt.Errorf("unmarshal yaml: %w", err)
	}

	yamlreflect.CopyFromYAMLStruct(r, yamlDst.Interface())
	return nil
}

// MarshalYAML implements [yaml.Marshaler]
func (r *SubscriptionRawPb) MarshalYAML() (any, error) {
	return yamlreflect.CopyToYAMLStruct(r), nil
}

// interface compile-time checker
var (
	_ yaml.Marshaler   = (*SetRequest)(nil)
	_ yaml.Unmarshaler = (*SetRequest)(nil)
	_ yaml.Marshaler   = (*SubscriptionRawPb)(nil)
	_ yaml.Unmarshaler = (*SubscriptionRawPb)(nil)
)
