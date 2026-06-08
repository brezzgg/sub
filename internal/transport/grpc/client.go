package grpc

import (
	"fmt"

	"github.com/brezzgg/sub/internal/transport/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	clientService pb.SubServiceClient
	clientConn    *grpc.ClientConn
)

func GetClient(remote string) (pb.SubServiceClient, error) {
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

	clientService = pb.NewSubServiceClient(clientConn)

	return clientService, nil
}
