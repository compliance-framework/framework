package cmd

import (
	"container-solutions.com/continuous-compliance/internal"
	"container-solutions.com/continuous-compliance/plugin"
	"container-solutions.com/continuous-compliance/plugins"
	"context"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"github.com/spf13/cobra"
	"log"
)

func AgentCmd() *cobra.Command {
	var agentCmd = &cobra.Command{
		Use:   "agent",
		Short: "long running agent for continuously checking policies against plugin data",
		Long: `The Continuous Compliance Agent is a long running process that continuously checks policy controls
with plugins to ensure continuous compliance.`,
		Run: RunAgent,
	}

	agentCmd.Flags().StringArray("policy-path", []string{}, "Directory where policies are stored")
	agentCmd.Flags().StringArray("policy-bundle", []string{}, "Directory where policies are stored")
	agentCmd.MarkFlagsOneRequired("policy-path", "policy-bundle")

	//agentCmd.Flags().StringArray("plugin-path", []string{}, "Directory where policies are stored")
	//agentCmd.MarkFlagsOneRequired("policy-path")

	return agentCmd
}

func RunAgent(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	// First we want to fetch all the policy bundles
	// Hardcoded until plugin system is incorporated
	pluginList := []plugin.Plugin{
		plugins.NewLocalSSH(),
	}

	bundles, err := cmd.Flags().GetStringArray("policy-bundle")
	if err != nil {
		internal.OnError(err, func(err error) {
			log.Fatal("Unable to retrieve policy bundles", err)
		})
	}

	// We have to load and evaluate the bundles with the plugins one by one, due to
	// https://www.openpolicyagent.org/docs/latest/management-bundles/#multiple-sources-of-policy-and-data.
	var queryBundles []*rego.Rego
	for _, inputBundle := range bundles {
		r := rego.New(
			rego.Query("data"),
			rego.LoadBundle(inputBundle),
		)
		// Check that it will be able to prepare when we're ready to run
		_, err = r.PrepareForEval(ctx)
		if err != nil {
			log.Fatal(err)
		}
		queryBundles = append(queryBundles, r)
	}

	for _, runnablePlugin := range pluginList {
		for _, queryBundle := range queryBundles {
			err = runnablePlugin.PrepareForEval(ctx)
			if err != nil {
				log.Fatal(err)
			}

			query, err := queryBundle.PrepareForEval(ctx)
			if err != nil {
				log.Fatal(err)
			}

			result, err := runnablePlugin.Evaluate(ctx, query)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(result)
		}
	}
}
