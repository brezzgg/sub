package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/brezzgg/sub/internal/models/payload/v1"
	"github.com/brezzgg/sub/internal/transport/grpc"
	"github.com/brezzgg/sub/internal/transport/grpc/pb"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Return subscription matching the given id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remote, err := cmd.Flags().GetString(remoteFlag)
		if err != nil {
			panic(err)
		}

		cl, err := grpc.GetClient(remote)
		if err != nil {
			fmt.Printf("grpc client error: %w\n", err)
			os.Exit(1)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		sub, err := cl.Get(ctx, &pb.GetRequest{Id: args[0]})
		if err != nil {
			fmt.Printf("get subscription: %s\n", err)
			os.Exit(1)
		}
		pay, err := payload.Unmarshal(sub.Payload, true)
		if err != nil {
			fmt.Printf("unmarshal payload: %s\n")
			os.Exit(1)
		}
		fmt.Printf(
			"id: %s\nttl: %d\ncreated_at: %d\npayload:\nbody:\n",
			args[0],
			sub.GetTtl(),
			sub.GetCreatedAt(),
		)
		bsplit := strings.Split(pay.Body, "\n")
		for i, b := range bsplit {
			fmt.Printf("  [%d] %s\n", i, b)
		}
		fmt.Println("headers:")
		i := 0
		for k, v := range pay.Headers {
			i++
			fmt.Printf("  [%d] %s: %s\n", i, k, v)
		}
	},
}
