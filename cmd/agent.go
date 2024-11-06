package cmd

import (
	"context"
	"fmt"
	"github.com/chris-cmsoft/concom/internal"
	cfplugin "github.com/chris-cmsoft/concom/plugin"
	"github.com/hashicorp/go-hclog"
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/open-policy-agent/opa/rego"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
)

func AgentCmd() *cobra.Command {
	var agentCmd = &cobra.Command{
		Use:   "agent",
		Short: "long running agent for continuously checking policies against plugin data",
		Long: `The Continuous Compliance Agent is a long running process that continuously checks policy controls
with plugins to ensure continuous compliance.`,
		Run: func(cmd *cobra.Command, args []string) {
			logger := hclog.New(&hclog.LoggerOptions{
				Name:   "agent",
				Output: os.Stdout,
				Level:  hclog.Debug,
			})
			runner := AgentRunner{
				logger: logger,
			}
			err := runner.Run(cmd, args)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	agentCmd.Flags().StringArray("policy-path", []string{}, "Directory where policies are stored")
	agentCmd.Flags().StringArray("policy-bundle", []string{}, "Directory where policies are stored")
	agentCmd.MarkFlagsOneRequired("policy-path", "policy-bundle")

	agentCmd.Flags().StringArray("plugin-path", []string{}, "Plugin executable")
	agentCmd.MarkFlagsOneRequired("plugin-path")

	return agentCmd
}

type AgentRunner struct {
	logger hclog.Logger
}

func (runner AgentRunner) Run(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()

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
			return err
		}
		queryBundles = append(queryBundles, r)
	}

	plugins, err := cmd.Flags().GetStringArray("plugin-path")
	if err != nil {
		return err
	}

	for _, path := range plugins {
		evaluator, err := runner.getExecPluginClient(path)
		if err != nil {
			return err
		}

		err = evaluator.PrepareForEval()
		if err != nil {
			return err
		}

		for _, queryBundle := range queryBundles {

			query, err := queryBundle.PrepareForEval(ctx)
			if err != nil {
				return err
			}

			fmt.Println(query)
			//result, err := evaluator.Evaluate(query)
			//if err != nil {
			//	log.Fatal(err)
			//}
		}

	}

	runner.closePluginClients()
	return nil
}

func (runner AgentRunner) getExecPluginClient(command string) (*cfplugin.EvaluatorRPC, error) {
	// We're a host! Start by launching the plugin process.
	client := goplugin.NewClient(&goplugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Managed:         true,
		Cmd:             exec.Command(command),
		Logger:          runner.logger,
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("evaluator")
	if err != nil {
		return nil, err
	}

	pluginRpc := raw.(*cfplugin.EvaluatorRPC)
	return pluginRpc, err
}

func (runner AgentRunner) closePluginClients() {
	goplugin.CleanupClients()
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = goplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]goplugin.Plugin{
	"evaluator": &cfplugin.EvaluatorPlugin{},
}
