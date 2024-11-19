package cmd

import (
	"github.com/chris-cmsoft/concom/internal/downloader"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
)

func DownloadPluginCmd() *cobra.Command {
	var agentCmd = &cobra.Command{
		Use:   "download-plugin",
		Short: "downloads plugins from OCI or URLs",
		Run: func(cmd *cobra.Command, args []string) {
			logger := hclog.New(&hclog.LoggerOptions{
				Output: os.Stdout,
				Level:  hclog.Debug,
			})
			downloadCmd := DownloadRunner{
				logger: logger,
			}
			err := downloadCmd.Run(cmd, args)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	var source []string
	agentCmd.Flags().StringArrayVarP(&source, "source", "s", source, "OCI or URL sources of the plugins")
	agentCmd.MarkFlagsOneRequired("source")

	return agentCmd
}

type DownloadRunner struct {
	logger hclog.Logger
}

func (d *DownloadRunner) Run(cmd *cobra.Command, args []string) error {
	sources, err := cmd.Flags().GetStringArray("source")
	if err != nil {
		return err
	}

	basePath, loopErr := os.Getwd()
	if loopErr != nil {
		return loopErr
	}

	// At some point, we will wrap this in go routine to download concurrently.
	// For the moment, we've left it without for the sake of simplicity and easy amendments.
	// We don't want to be hassled with channels and scoped variables if we need to refactor this during implementation.
	for _, source := range sources {
		d.logger.Debug("Received source", "source", source)

		downloaderImpl, err := downloader.Download(source, path.Join(basePath, AgentPluginDir))
		if err != nil {
			return err
		}

		finalDestination, err := downloaderImpl.GetFinalDestination()
		if err != nil {
			return err
		}

		d.logger.Debug("Downloaded plugin", "path", finalDestination)
	}

	return nil
}
