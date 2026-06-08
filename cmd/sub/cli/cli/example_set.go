package cli

import (
	"fmt"
	"os"

	"github.com/brezzgg/sub/internal/option/subset"
	"github.com/spf13/cobra"
)

var exampleSetCmd = &cobra.Command{
	Use:   "example_set out.yaml",
	Short: "Generate example file for set command, user stdout to write in stdout",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "stdout" {
			if err := subset.WriteExampleStdout(); err != nil {
				fmt.Printf("failed to write example: %s\n", err)
				os.Exit(1)
			}
		} else {
			if err := subset.WriteExample(args[0]); err != nil {
				fmt.Printf("failed to write example file: %s\n", err)
			}
		}
	},
}
