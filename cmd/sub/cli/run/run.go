package run

import (
	"net"
	"strings"

	"github.com/brezzgg/go-packages/lg"
	"github.com/brezzgg/sub/internal/models"
	"github.com/brezzgg/sub/internal/repo/sqlite"
	"github.com/brezzgg/sub/internal/transport/grpc"
	"github.com/brezzgg/sub/internal/transport/grpc/pb"
	"github.com/brezzgg/sub/internal/transport/http"
	"github.com/spf13/cobra"
	G "google.golang.org/grpc"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Start sub server",
	Run: func(cmd *cobra.Command, args []string) {
		repof, err := cmd.Flags().GetString(repoFlag)
		if err != nil {
			lg.Fatal(err)
		}
		repoFile, err := cmd.Flags().GetString(repoFileFlag)
		if err != nil {
			lg.Fatal(err)
		}
		host, err := cmd.Flags().GetString(hostFlag)
		if err != nil {
			lg.Fatal(err)
		}
		grpcHost, err := cmd.Flags().GetString(grpcHostFlag)
		if err != nil {
			lg.Fatal(err)
		}
		httpPattern, err := cmd.Flags().GetString(httpPatternFlag)
		if err != nil {
			lg.Fatal(err)
		}

		if !strings.Contains(httpPattern, "{id}") {
			lg.Warn("http pattern not contain `{id}`, switched to default")
			httpPattern = httpPatternFlagDefault
		}

		lg.Info("options", lg.C{
			"repo":         repof,
			"repo_file":    repoFile,
			"host":         host,
			"grpc_host":    grpcHost,
			"http_pattern": httpPattern,
		})

		var repo models.Repo
		switch repof {
		case "sqlite":
			r, err := sqlite.NewRepo(repoFile)
			if err != nil {
				lg.Fatal("failed to init repo", err)
			}
			repo = r
		default:
			lg.Fatal("unknown repository")
		}
		lg.Info("repo load successful")

		lis, err := net.Listen("tcp", grpcHost)
		if err != nil {
			lg.Fatal("failed to listen", err, lg.C{"host": grpcHost})
		}
		grpcserv := G.NewServer(G.UnaryInterceptor(grpc.LoggingInterceptor))
		pb.RegisterSubServiceServer(grpcserv, grpc.NewServer(repo))

		httpserv := http.New(host, repo, httpPattern)

		sch := &scheduler{}

		lg.Info("grpc server configured")
		sch.Add(func() error {
			return grpcserv.Serve(lis)
		}, func() error {
			grpcserv.GracefulStop()
			return nil
		})

		lg.Info("http server configured")
		sch.Add(httpserv.Run, httpserv.Stop)

		if err := sch.Run(); err != nil {
			lg.Error(err, lg.Sync{})
		}
	},
}

const (
	hostFlag               = "host"
	repoFlag               = "repo"
	repoFileFlag           = "repo-file"
	grpcHostFlag           = "grpc-host"
	httpPatternFlag        = "http-pattern"
	httpPatternFlagDefault = "GET /i/{id}"
)

func init() {
	RunCmd.PersistentFlags().StringP(hostFlag, "H", "0.0.0.0:8080", "select host")
	RunCmd.PersistentFlags().String(grpcHostFlag, "0.0.0.0:50051", "select grpc host")
	RunCmd.PersistentFlags().String(repoFlag, "sqlite", "select repository")
	RunCmd.PersistentFlags().String(repoFileFlag, "./subs.db", "choose repository file if needed")
	RunCmd.PersistentFlags().String(
		httpPatternFlag,
		httpPatternFlagDefault,
		`select an HTTP pattern for the subscription endpoint. the pattern must contain {id}.`,
	)
}
