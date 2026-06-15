package cli

import (
	"context"
	"time"

	"github.com/brezzgg/sub/cmd/sub/cli/common/errors"
	"github.com/brezzgg/sub/cmd/sub/cli/common/log"
	"github.com/brezzgg/sub/internal/transport/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var rmCmd = &cobra.Command{
	Use:     "remove <id>",
	Aliases: []string{"rm"},
	Short:   "Remove subscription by id",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, err := cmd.Flags().GetBool(rmForceFlag)
		if err != nil {
			log.Fatal(errors.ErrInvalidArg(err).Error())
		}

		id := args[0]

		cl := mgr.GrpcClient()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		_, err = cl.Remove(ctx, &grpc.IdRequest{Id: id})
		if err != nil {
			st, ok := status.FromError(err)
			if !ok {
				log.Fatal(errors.ErrUnknown(err).Error())
			}
			if st.Code() == codes.NotFound {
				if force {
					return
				}
				log.Fatal(errors.ErrSubNotFound().Error())
			}
			log.Fatal(errors.ErrRemoteInternal(err).Error())
		}
	},
}

const (
	rmForceFlag = "force"
)

func init() {
	rmCmd.PersistentFlags().BoolP(rmForceFlag, "f", false, "force remove")
}
