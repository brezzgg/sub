package cli

import (
	"os"

	"github.com/brezzgg/sub/cmd/sub/cli/common/errors"
	"github.com/brezzgg/sub/cmd/sub/cli/common/log"
	"github.com/brezzgg/sub/example"
	"github.com/spf13/cobra"
)

var exampleSetCmd = &cobra.Command{
	Use:   "example_set out.yaml",
	Short: "Generate example file for set command, user stdout to write in stdout",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		buf := example.SubscriptionSet
		var file *os.File
		if args[0] == "stdout" {
			file = os.Stdout
		} else {
			f, err := os.Create(args[0])
			defer f.Close()
			if err != nil {
				log.Fatal(errors.ErrWriteFile(args[0], err).Error())
			}
			file = f
		}
		file.Write(buf)
	},
}
