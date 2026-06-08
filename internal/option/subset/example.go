package subset

import (
	"os"

	"github.com/brezzgg/sub/example"
)

func WriteExample(path string) error {
	return os.WriteFile(path, example.SubscriptionSet, 0755)
}

func WriteExampleStdout() error {
	_, err := os.Stdout.Write(example.SubscriptionSet)
	return err
}
