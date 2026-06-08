package cli

import (
	"github.com/brezzgg/sub/cmd/sub/cli/cli"
	"github.com/brezzgg/sub/cmd/sub/cli/run"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "sub",
}

func Execute() {
	_ = rootCmd.Execute()
}

func addSubcommand() {
	rootCmd.AddCommand(run.RunCmd)
	rootCmd.AddCommand(cli.CliCmd)
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	addSubcommand()
}
