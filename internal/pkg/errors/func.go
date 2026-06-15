package errors

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
)

func Directf(f string, opts ...any) error {
	return Direct(fmt.Sprintf(f, opts...))
}

func Errorf(code ErrorCode, f string, opts ...any) error {
	return Error(code, fmt.Sprintf(f, opts...))
}

func Direct(msg string) error {
	return NewDomain(CodeDirect, msg)
}

func Error(code ErrorCode, msg string) error {
	return NewDomain(code, msg)
}

func Code(code ErrorCode) error {
	return NewDomain(code, "")
}

func AsDomain(err error) (*Domain, bool) {
	if err == nil {
		return nil, false
	}
	d, ok := err.(*Domain)
	return d, ok
}

func IsUnfixable(err error) bool {
	d, ok := AsDomain(err)
	if !ok {
		return false
	}
	if d.Code() != CodeUnfixable {
		return false
	}
	return true
}

func CodeIs(err error, code ErrorCode) bool {
	if d, ok := AsDomain(err); ok {
		return d.Code() == code
	}
	return false
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err, target error) bool {
	return errors.As(err, target)
}

func AddContext(err error, context string) error {
	d, ok := err.(*Domain)
	if ok {
		d.msg = fmt.Sprintf("%s: %s", context, d.msg)
		return d
	}
	return err
}

func internal(msg string) {
	os.Stderr.WriteString(fmt.Sprintf("errors internal: %s\n", msg))
	debug.PrintStack()
}
