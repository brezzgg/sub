package errors

import (
	"fmt"
	"strings"

	"github.com/brezzgg/go-packages/lg"
)

type Domain struct {
	msg  string
	code ErrorCode
}

func NewDomain(code ErrorCode, msg string) *Domain {
	if code != CodeDirect && msg != "" {
		level.Log(lg.GlobalLogger, "error", lg.C{"error": msg})
	}
	return &Domain{
		code: code,
		msg:  msg,
	}
}

var level = lg.NewLogLevel(lg.ClrFgPink, "Domain", lg.LevelOptionDisableCaller)

// Error implements [error].
func (e *Domain) Error() string {
	if e == nil {
		internal("error is nil")
		return ""
	}
	if e.code == CodeDirect {
		return e.msg
	}
	if answ, ok := CodeMap[e.code]; ok {
		return fmt.Sprintf("%d %s", e.code, answ)
	}
	internal(fmt.Sprintf("unknown code %d", int(e.code)))
	return ""
}

func (e *Domain) Code() ErrorCode {
	return e.code
}

type Domains struct {
	errs []*Domain
}

// Error implements [error].
func (e *Domains) Error() string {
	if e == nil {
		internal("error slice is nil")
		return ""
	}
	if len(e.errs) == 0 {
		internal("error slice length is zero")
		return ""
	}
	errs := make([]string, 0, len(e.errs))
	for _, err := range e.errs {
		if err != nil {
			if e := err.Error(); e != "" {
				errs = append(errs, e)
			}
		}
	}
	if len(errs) == 0 {
		internal("error slice length is zero")
		return ""
	}
	return fmt.Sprintf("%d erros was occured: %s", len(errs), strings.Join(errs, "; "))
}

type trivialError struct {
	msg string
}

// Error implements [error].
func (t *trivialError) Error() string {
	return t.msg
}

var _ error = (*Domain)(nil)
var _ error = (*Domains)(nil)
var _ error = (*trivialError)(nil)
