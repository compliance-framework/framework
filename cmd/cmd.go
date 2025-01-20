package cmd

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cf",
		Short: "cf manages policies for the compliance framework",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
	cmd.AddCommand(
		ApiCmd(),
		AgentCmd(),
		DownloadPluginCmd(),
		DownloadPolicyCmd(),
	)
	return cmd
}
