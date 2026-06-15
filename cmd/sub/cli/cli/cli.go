package cli

import (
	"github.com/brezzgg/sub/cmd/sub/cli/common/errors"
	"github.com/brezzgg/sub/cmd/sub/cli/common/log"
	"github.com/brezzgg/sub/internal/manager"
	"github.com/spf13/cobra"
)

var mgr *manager.Manager

var CliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Client commands",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		remote, err := cmd.Flags().GetString(remoteFlag)
		if err != nil {
			log.Fatal(errors.ErrInvalidArg(err).Error())
		}
		m, err := manager.NewClientManager(
			manager.WithGrpcClient(remote),
		)
		if err != nil {
			log.Fatal(errors.ErrFailGrpcConn(err).Error())
		}
		mgr = m
	},
}

const (
	remoteFlag = "remote"
)

func init() {
	addSubcommands()

	CliCmd.PersistentFlags().StringP(remoteFlag, "r", "127.0.0.1:50051", "specify remote host")
}

func addSubcommands() {
	CliCmd.AddCommand(setCmd)
	CliCmd.AddCommand(getCmd)
	CliCmd.AddCommand(rmCmd)
	CliCmd.AddCommand(randidCmd)
	CliCmd.AddCommand(exampleSetCmd)
	CliCmd.AddCommand(enableCmd)
	CliCmd.AddCommand(disableCmd)
}
