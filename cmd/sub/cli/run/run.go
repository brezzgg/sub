package run

import (
	"strings"

	"github.com/brezzgg/go-packages/lg"
	"github.com/brezzgg/sub/internal/manager"
	"github.com/brezzgg/sub/internal/repo"
	"github.com/spf13/cobra"
)

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Start sub server",
	Run: func(cmd *cobra.Command, args []string) {
		storageType, err := cmd.Flags().GetString(storageFlag)
		if err != nil {
			lg.Fatal(err)
		}
		storageFile, err := cmd.Flags().GetString(storageFileFlag)
		if err != nil {
			lg.Fatal(err)
		}
		storageRemote, err := cmd.Flags().GetString(storageRemoteFlag)
		if err != nil {
			lg.Fatal(err)
		}
		cacheType, err := cmd.Flags().GetString(cacheFlag)
		if err != nil {
			lg.Fatal(err)
		}
		cacheFile, err := cmd.Flags().GetString(storageFileFlag)
		if err != nil {
			lg.Fatal(err)
		}
		cacheRemote, err := cmd.Flags().GetString(storageRemoteFlag)
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
			"storage_type":   storageType,
			"storage_file":   storageFile,
			"storage_remote": storageRemote,
			"host":           host,
			"grpc_host":      grpcHost,
			"http_pattern":   httpPattern,
		})

		repoOpts := &repo.Options{
			StorageProvider: storageType,
			StorageOpts: map[string]any{
				"file":   storageFile,
				"remote": storageRemote,
			},
			CacheProvider: cacheType,
			CacheOpts: map[string]any{
				"file":   cacheFile,
				"remote": cacheRemote,
			},
		}

		mgr, err := manager.NewManager(
			repoOpts,
			nil,
			manager.WithGrpcService(grpcHost),
			manager.WithHttpHandler(host, httpPattern),
		)
		if err != nil {
			lg.Fatal("manager init", err)
		}

		if err := mgr.Run(); err != nil {
			lg.Fatal(err, lg.Sync{})
		}
	},
}

const (
	hostFlag               = "host"
	cacheFlag              = "cache"
	cacheFileFlag          = "cache-file"
	cacheRemoteFlag        = "cache-remote"
	storageFlag            = "storage"
	storageFileFlag        = "storage-file"
	storageRemoteFlag      = "storage-remote"
	grpcHostFlag           = "grpc-host"
	httpPatternFlag        = "http-pattern"
	httpPatternFlagDefault = "GET /i/{id}"
)

func init() {
	RunCmd.PersistentFlags().StringP(hostFlag, "H", "0.0.0.0:8080", "select host")
	RunCmd.PersistentFlags().String(grpcHostFlag, "0.0.0.0:50051", "select grpc host")
	RunCmd.PersistentFlags().String(cacheFlag, "", "select cache")
	RunCmd.PersistentFlags().String(cacheFileFlag, "", "select cache file/directory if needed")
	RunCmd.PersistentFlags().String(cacheRemoteFlag, "", "select cache remote host if needed")
	RunCmd.PersistentFlags().String(storageFlag, "", "select storage")
	RunCmd.PersistentFlags().String(storageFileFlag, "", "select storage file/directory if needed")
	RunCmd.PersistentFlags().String(storageRemoteFlag, "", "select storage remote host if needed")
	RunCmd.PersistentFlags().String(
		httpPatternFlag,
		httpPatternFlagDefault,
		`select an HTTP pattern for the subscription endpoint. the pattern must contain {id}.`,
	)
}
