package grpc

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	clientService SubServiceClient
	clientConn    *grpc.ClientConn
)

func GetClient(remote string) (SubServiceClient, error) {
	if clientService != nil && clientConn != nil {
		if clientConn.GetState() == connectivity.Ready {
			return clientService, nil
		}
	}

	var err error
	clientConn, err = grpc.NewClient(
		remote,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("connection: %s", err)
	}

	clientService = NewSubServiceClient(clientConn)

	return clientService, nil
}
