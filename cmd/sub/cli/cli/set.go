package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/brezzgg/go-packages/lg"
	"github.com/brezzgg/sub/internal/models"
	"github.com/brezzgg/sub/internal/models/payload/v1"
	"github.com/brezzgg/sub/internal/option/subset"
	"github.com/brezzgg/sub/internal/transport/grpc"
	"github.com/brezzgg/sub/internal/transport/grpc/pb"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var setCmd = &cobra.Command{
	Use:   "set <file.yaml>",
	Short: "add or update subscription from a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		remote, err := cmd.Flags().GetString(remoteFlag)
		if err != nil {
			panic(err)
		}
		id, err := cmd.Flags().GetString(idFlag)
		if err != nil {
			panic(err)
		}
		randId, err := cmd.Flags().GetBool(randidFlag)
		if err != nil {
			panic(err)
		}

		fileb, err := os.ReadFile(args[0])
		if err != nil {
			fmt.Printf("read file error: %s\n", err)
			os.Exit(1)
		}

		var subcfg subset.Model
		if err := yaml.Unmarshal(fileb, &subcfg); err != nil {
			fmt.Printf("invalid yaml format: %s\n", err)
			os.Exit(1)
		}

		var p payload.Payload
		p.Body = strings.Join(subcfg.Bodies, "\n")

		p.Headers = make(map[string]string, 10+len(subcfg.CustomHeaders))
		for k, v := range subcfg.CustomHeaders {
			p.Headers[k] = v
		}

		p.Headers["Profile-Title"] = subcfg.Title
		p.Headers["profile-update-interval"] = subcfg.UpdateInterval
		p.Headers["profile-web-page-url"] = subcfg.UserWebPage
		p.Headers["support-url"] = subcfg.SupportUri
		p.Headers["moved-permanently-to"] = subcfg.MovedPermanently
		p.Headers["DNS"] = subcfg.PreferredDNS

		pmsg, err := payload.Marshal(&p)
		if err != nil {
			fmt.Errorf("failed to marshal payload: %s", err)
		}

		setreq := &pb.SetRequest{
			Payload: pmsg,
			Ttl:     int64(subcfg.TTL),
		}

		if subcfg.Id != "" {
			setreq.Id = subcfg.Id
		} else {
			if randId {
				setreq.Id = randIdFunc()
			} else {
				if id != "" {
					if models.ValidateSubId(id) {
						setreq.Id = id
					} else {
						fmt.Printf("invalid id format: %s\n", id)
						os.Exit(1)
					}
				} else {
					fmt.Printf("id for subscription not set\n")
					os.Exit(1)
				}
			}
		}

		cl, err := grpc.GetClient(remote)
		if err != nil {
			fmt.Errorf("grpc conn error: %s\n", err)
			os.Exit(1)
		}

		_, err = cl.Set(context.TODO(), setreq)
		if err != nil {
			lg.Fatal("set error: %s", err)
		}

		if randId {
			fmt.Printf("%s\n", setreq.Id)
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
