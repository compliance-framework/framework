package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/chris-cmsoft/concom/runner"
	"github.com/chris-cmsoft/concom/runner/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/open-policy-agent/opa/rego"
	"github.com/spf13/cobra"
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
				Level:  hclog.Trace,
			})
			pluginRunner := AgentRunner{
				logger: logger,
			}
			err := pluginRunner.Run(cmd, args)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	agentCmd.Flags().StringArray("policy", []string{}, "Directory or Bundle archive where policies are stored")
	err := agentCmd.MarkFlagRequired("policy")
	if err != nil {
		log.Fatal(err)
	}

	agentCmd.Flags().StringArray("plugin", []string{}, "Plugin executable or directory")
	agentCmd.MarkFlagsOneRequired("plugin")

	agentCmd.Flags().BoolP("daemon", "d", false, "Specify to run as a long running daemon")

	return agentCmd
}

type AgentRunner struct {
	logger hclog.Logger

	queryBundles []*rego.Rego
}

func (ar AgentRunner) Run(cmd *cobra.Command, args []string) error {
	//ctx := context.TODO()

	policyBundles, err := cmd.Flags().GetStringArray("policy")
	if err != nil {
		return err
	}

	plugins, err := cmd.Flags().GetStringArray("plugin")
	if err != nil {
		return err
	}

	daemon, err := cmd.Flags().GetBool("daemon")
	if err != nil {
		return err
	}

	if daemon == true {
		for {
			err := ar.runInstance(plugins, policyBundles)

			if err != nil {
				ar.logger.Error("error running instance", "error", err)
				// No return for now, we'll do a retry afterwards.
				// TODO: Should we have a retry limit maybe?
			}

			time.Sleep(time.Second * 60)
		}
	} else {
		err := ar.runInstance(plugins, policyBundles)

		if err != nil {
			ar.logger.Error("error running instance", "error", err)
			return err
		}
	}

	return nil
}

func (ar AgentRunner) runInstance(
	plugins []string,
	policyBundles []string,
) error {
	defer ar.closePluginClients()

	for _, path := range plugins {
		logger := hclog.New(&hclog.LoggerOptions{
			Name:   "runner",
			Output: os.Stdout,
			Level:  hclog.Debug,
		})

		runnerInstance, err := ar.getRunnerInstance(logger, path)
		if err != nil {
			return err
		}

		_, err = runnerInstance.Configure(&proto.ConfigureRequest{
			Config: map[string]string{
				"host": "127.0.0.1",
				"port": "22",
			},
		})
		if err != nil {
			return err
		}

		_, err = runnerInstance.PrepareForEval(&proto.PrepareForEvalRequest{})
		if err != nil {
			return err
		}

		for _, inputBundle := range policyBundles {
			res, err := runnerInstance.Eval(&proto.EvalRequest{
				BundlePath: inputBundle,
			})
			if err != nil {
				return err
			}

			fmt.Println("Output from runner:")
			fmt.Println("Findings:", res.Findings)
			fmt.Println("Observations:", res.Observations)
			fmt.Println("Log Entries:", res.Logs)

			// Here we'll send the data back to NATS
		}
	}

	return nil
}

func (ar AgentRunner) getRunnerInstance(logger hclog.Logger, path string) (runner.Runner, error) {
	// We're a host! Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig:  runner.HandshakeConfig,
		Plugins:          runner.PluginMap,
		Managed:          true,
		Cmd:              exec.Command(path),
		Logger:           logger,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("runner")
	if err != nil {
		return nil, err
	}

	// We should have a Greeter now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	runnerInstance := raw.(runner.Runner)
	return runnerInstance, nil
}

func (ar AgentRunner) closePluginClients() {
	plugin.CleanupClients()
}
