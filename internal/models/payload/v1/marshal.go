package payload

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"google.golang.org/protobuf/proto"
)

func Unmarshal(b []byte, decodeBody bool) (*Payload, error) {
	var p Payload
	if err := proto.Unmarshal(b, &p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pb: %s", err)
	}
	if decodeBody {
		b, err := base64.StdEncoding.DecodeString(p.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to decode body: %s", err)
		}
		p.Body = string(b)
	}
	return &p, nil
}

var (
	newlineRe = regexp.MustCompile(`[\n]+`)
	spaceRe   = regexp.MustCompile(`[\s]+`)
)

func Marshal(p *Payload) ([]byte, error) {
	p.Body = strings.TrimSpace(p.Body)
	p.Body = strings.Join(strings.Fields(p.Body), "\n")
	p.Body = base64.StdEncoding.EncodeToString([]byte(p.Body))
	return proto.Marshal(p)
}
