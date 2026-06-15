package cli

import (
	"crypto/rand"
	"math/big"

	"github.com/brezzgg/sub/cmd/sub/cli/common/log"
	"github.com/spf13/cobra"
)

var randidCmd = &cobra.Command{
	Use:   "rand_id",
	Short: "Generate random subscription id",
	Run: func(cmd *cobra.Command, args []string) {
		log.Print(randIdFunc())
	},
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_"

func randIdFunc() string {
	const n = 128
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		x, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[x.Int64()]
	}
	return string(b)
}
