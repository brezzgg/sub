package cli

import (
	"github.com/brezzgg/go-packages/lg"
	"github.com/brezzgg/sub/internal/transport/grpc"
	"github.com/spf13/cobra"
)

var CliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Client commands",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		remote, err := cmd.Flags().GetString(remoteFlag)
		if err != nil {
			lg.Fatal(err)
		}
		if _, err := grpc.GetClient(remote); err != nil {
			lg.Fatal("failed to connect to remote", err)
		}
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
	CliCmd.AddCommand(randidCmd)
	CliCmd.AddCommand(exampleSetCmd)
}
