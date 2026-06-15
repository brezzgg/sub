package errors

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/status"
)

var (
	ErrFailedClientConn = errors.New("failed to connect to remote ")
)

func ErrSubNotFound() error {
	return errors.New("subscription not found")
}

func ErrSubsNotFound() error {
	return errors.New("subscriptions not found")
}

func ErrFailGrpcConn(grpcErr error) error {
	const plain = "connection to grpc server failed: %s"
	gerr, ok := status.FromError(grpcErr)
	if ok {
		return fmt.Errorf(plain, gerr.Message())
	} else {
		return fmt.Errorf(plain, gerr)
	}
}

func ErrReadFile(path string, err error) error {
	return fmt.Errorf("failed to read file %s: %s", path, err)
}

func ErrWriteFile(path string, err error) error {
	return fmt.Errorf("failed to write file %s: %s", path, err)
}

func ErrUnmarshal(target string, err error) error {
	return fmt.Errorf("failed to unmarshal %s: %s", target, err)
}

func ErrMarshal(target string, err error) error {
	return fmt.Errorf("failed to marshal %s: %s", target, err)
}

func ErrInvalidArg(err error) error {
	return fmt.Errorf("invalid argument: %s", err)
}

func ErrRemoteInternal(err error) error {
	return fmt.Errorf("remote internal error: %s", err)
}

func ErrUnknown(err error) error {
	return fmt.Errorf("unknown: %s", err)
}
