package cli

import (
	"context"
	"time"

	"github.com/brezzgg/sub/cmd/sub/cli/common/log"
	"github.com/brezzgg/sub/internal/transport/grpc"
	"github.com/spf13/cobra"
)

var enableCmd = &cobra.Command{
	Use:   "enable <subscription_id>",
	Short: "Enable subscription",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		setEnable(id, true)
	},
}

var disableCmd = &cobra.Command{
	Use:   "disable <subscription_id>",
	Short: "Disable subscription",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		setEnable(id, false)
	},
}

func setEnable(id string, enabled bool) {
	cl := mgr.GrpcClient()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	_, err := cl.SetEnabled(ctx, &grpc.SetEnabledRequest{
		Id: &grpc.IdRequest{
			Id: id,
		},
		Enabled: enabled,
	})
	if err != nil {
		log.Fatal("%s", err)
	}
}
