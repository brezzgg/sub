package cli

import (
	"context"
	"os"
	"time"

	"github.com/brezzgg/sub/cmd/sub/cli/common/errors"
	"github.com/brezzgg/sub/cmd/sub/cli/common/log"
	"github.com/brezzgg/sub/internal/transport/grpc"
	"github.com/brezzgg/sub/internal/usecase"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var setCmd = &cobra.Command{
	Use:     "set <file.yaml>",
	Short:   "Add or update subscription from a file",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"apply"},
	Run: func(cmd *cobra.Command, args []string) {
		id, err := cmd.Flags().GetString(idFlag)
		if err != nil {
			log.Fatal(errors.ErrInvalidArg(err).Error())
		}
		randId, err := cmd.Flags().GetBool(randidFlag)
		if err != nil {
			log.Fatal(errors.ErrInvalidArg(err).Error())
		}
		fileName := args[0]

		in, err := os.ReadFile(fileName)
		if err != nil {
			log.Fatal(errors.ErrReadFile(fileName, err).Error())
		}

		var sr usecase.SetRequest
		if err := yaml.Unmarshal(in, &sr); err != nil {
			log.Fatal(errors.ErrUnmarshal(fileName, err).Error())
		}

		sr.Raw = &usecase.SubscriptionRawPb{}
		if err := yaml.Unmarshal(in, sr.Raw); err != nil {
			log.Fatal(errors.ErrUnmarshal(fileName, err).Error())
		}

		reqId := ""
		if sr.Raw.Id != "" {
			reqId = sr.Raw.Id
		}
		if randId {
			reqId = randIdFunc()
		}
		if id != "" {
			reqId = id
		}
		sr.Raw.Id = reqId

		cl := mgr.GrpcClient()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		_, err = cl.Set(ctx, &grpc.SetRequest{
			Data: &grpc.SetRequest_Req{Req: &sr},
		})
		if err != nil {
			log.Fatal("set error: %s", err)
		}

		if randId {
			log.Print("%s", sr.Raw.Id)
		}
	},
}

const (
	idFlag     = "id"
	randidFlag = "rand-id"
)

func init() {
	setCmd.PersistentFlags().String(idFlag, "", "specify id of subscription")
	setCmd.PersistentFlags().Bool(randidFlag, false, "generate random id of subscription")
}
