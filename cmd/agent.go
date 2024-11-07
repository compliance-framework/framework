package cmd

import (
	"github.com/chris-cmsoft/concom/runner"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
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
				Level:  hclog.Trace,
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

	queryBundles []*rego.Rego
}

func (ar AgentRunner) Run(cmd *cobra.Command, args []string) error {
	//ctx := context.TODO()

	//bundles, err := cmd.Flags().GetStringArray("policy-bundle")
	//if err != nil {
	//	internal.OnError(err, func(err error) {
	//		log.Fatal("Unable to retrieve policy bundles", err)
	//	})
	//}
	//
	//// First we'll load the file based bundles as Rego queries.
	//// These will be evaluated one at a time, to avoid any root conflicts in packages as they
	//// all will fall under `package compliance_framework.XXX`
	////
	//// Why this is necessary:
	//// https://www.openpolicyagent.org/docs/latest/management-bundles/#multiple-sources-of-policy-and-data.
	//for _, inputBundle := range bundles {
	//	r := rego.New(
	//		rego.Query("data"),
	//		rego.LoadBundle(inputBundle),
	//	)
	//
	//	// Check that it will be able to prepare when we're ready to run
	//	_, err = r.PrepareForEval(ctx)
	//	if err != nil {
	//		return err
	//	}
	//	ar.queryBundles = append(ar.queryBundles, r)
	//}

	plugins, err := cmd.Flags().GetStringArray("plugin-path")
	if err != nil {
		return err
	}

	defer ar.closePluginClients()

	for _, path := range plugins {
		logger := hclog.New(&hclog.LoggerOptions{
			Name:   "runner",
			Output: os.Stdout,
			Level:  hclog.Debug,
		})

		runnerInstance, err := ar.GetRunnerInstance(logger, path)
		if err != nil {
			return err
		}

		//err = runnerInstance.Configure(runner.RunnerConfig{
		//	"host": "192.168.1.1",
		//})
		//if err != nil {
		//	return err
		//}

		err = runnerInstance.Configure(map[string]string{
			"host": "127.0.0.1",
			"port": "80",
		})
		if err != nil {
			return err
		}

		err = runnerInstance.PrepareForEval()
		if err != nil {
			return err
		}

		//
		//for _, queryBundle := range runner.queryBundles {
		//	fmt.Println("-------------")
		//	query, err := queryBundle.PrepareForEval(ctx)
		//	if err != nil {
		//		return err
		//	}
		//	fmt.Println(query)
		//
		//	//result, err := evaluator.Evaluate(query)
		//	//if err != nil {
		//	//	log.Fatal(err)
		//	//}
		//	//fmt.Println(result)
		//}

	}

	return nil
}

func (ar AgentRunner) GetRunnerInstance(logger hclog.Logger, path string) (runner.Runner, error) {
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
