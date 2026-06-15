package cli

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/brezzgg/sub/cmd/sub/cli/common/errors"
	"github.com/brezzgg/sub/cmd/sub/cli/common/log"
	"github.com/brezzgg/sub/internal/entity"
	"github.com/brezzgg/sub/internal/transport/grpc"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var getCmd = &cobra.Command{
	Use:   "get <id|.>",
	Short: "Return subscription matching the given id or all subscriptions if set '.'",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, err := cmd.Flags().GetBool(getForceFlag)
		if err != nil {
			log.Fatal(errors.ErrInvalidArg(err).Error())
		}
		full, err := cmd.Flags().GetBool(getFullFlag)
		if err != nil {
			log.Fatal(errors.ErrInvalidArg(err).Error())
		}
		wide, err := cmd.Flags().GetBool(getWideFlag)
		if err != nil {
			log.Fatal(errors.ErrInvalidArg(err).Error())
		}

		id := args[0]

		if id == "." {
			getAll(cmd, args, force, full, wide)
		} else {
			getOne(cmd, args, force, full, wide)
		}
	},
}

func getOne(cmd *cobra.Command, args []string, force, full, wide bool) {
	cl := mgr.GrpcClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	sub, err := cl.Get(ctx, &grpc.IdRequest{Id: args[0]})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.NotFound {
				if force {
					return
				}
				log.Fatal(errors.ErrSubNotFound().Error())
			}
		}
		log.Fatal(errors.ErrRemoteInternal(err).Error())
	}

	_, sb, err := sub.ToSubscription()
	if err != nil {
		log.Fatal("failed to unmarshal subscription: %s", err)
	}

	if !full {
		getWriteSubOneline(args[0], sb, wide)
		return
	}

	getWriteSub(args[0], sb)
}

func getAll(cmd *cobra.Command, args []string, force, full, wide bool) {
	cl := mgr.GrpcClient()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	all, err := cl.GetAll(ctx, &grpc.Empty{})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			if st.Code() == codes.NotFound {
				log.Fatal(errors.ErrSubsNotFound().Error())
			}
		}
		log.Fatal(errors.ErrRemoteInternal(err).Error())
	}

	for _, sub := range all.GetSubs() {
		id, sb, err := sub.ToSubscription()
		if err != nil {
			log.Error("unmarshal subscription: %s", err)
			continue
		}
		if !full {
			getWriteSubOneline(id, sb, wide)
			continue
		}
		getWriteSub(id, sb)
	}
}

func getWriteSubOneline(id string, sub *entity.Subscription, wide bool) {
	n := 15
	shortId := id
	if !wide {
		idRunes := []rune(shortId)
		if len(idRunes) > n+1 {
			shortId = string(idRunes[:n]) + "..."
		} else {
			shortId = shortId + strings.Repeat(" ", n+3-len(idRunes))
		}
	}
	log.Log(
		"id: %s(%s)  expired: %s  payload_body_len: %d  payload_header_count: %d",
		shortId,
		strconv.FormatBool(sub.Ok() == nil),
		expiredString(sub.Expired),
		len(sub.Payload.Body),
		len(sub.Payload.Headers),
	)
}

func getWriteSub(id string, sub *entity.Subscription) {
	log.Log(
		"id: %s\nenabled: %s\nexpired: %s\npayload:\nbody:",
		id,
		strconv.FormatBool(sub.Ok() == nil),
		expiredString(sub.Expired),
	)
	body, err := sub.Payload.UnmarshalBody()
	if err != nil {
		log.Error("failed to unmarshal payload body: %s", err)
		return
	}
	bsplit := strings.Split(body, "\n")
	for i, b := range bsplit {
		log.Log("  [%d] %s", i+1, b)
	}
	log.Log("headers:")
	i := 0
	for k, v := range sub.Payload.Headers {
		i++
		log.Log("  [%d] %s: %s", i, k, v)
	}
}

func expiredString(t time.Time) string {
	if t.Year() < 2000 {
		return "never"
	}
	return time.Now().UTC().Format(time.DateTime)
}

const (
	getForceFlag = "force"
	getFullFlag  = "full"
	getWideFlag  = "wide"
)

func init() {
	getCmd.PersistentFlags().BoolP(getForceFlag, "f", false, "force get")
	getCmd.PersistentFlags().Bool(getFullFlag, false, "show full subscription data")
	getCmd.PersistentFlags().Bool(getWideFlag, false, "show full id")
}
