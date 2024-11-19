package cmd

import (
	"fmt"
	"github.com/chris-cmsoft/concom/internal/downloader"
	"github.com/chris-cmsoft/concom/runner"
	"github.com/chris-cmsoft/concom/runner/proto"
	"github.com/coreos/go-systemd/v22/daemon"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/open-policy-agent/opa/rego"
	"github.com/spf13/cobra"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"syscall"
	"time"
)

const AgentPluginDir = ".compliance-framework/plugins"

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

func (ar *AgentRunner) Run(cmd *cobra.Command, args []string) error {
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
		ar.runDaemon(plugins, policyBundles)
	} else {
		err := ar.runInstance(plugins, policyBundles)

		if err != nil {
			ar.logger.Error("error running instance", "error", err)
			return err
		}
	}

	return nil
}

// Should never return, either handles any error or panics.
func (ar *AgentRunner) runDaemon(plugins []string, policyBundles []string) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		ar.logger.Info("received signal to terminate plugins and exit", "signal", sig)
		ar.closePluginClients()
		os.Exit(0)
	}()

	go daemon.SdNotify(false, "READY=1")

	for {
		err := ar.runInstance(plugins, policyBundles)

		if err != nil {
			ar.logger.Error("error running instance", "error", err)
			// No return for now, we keep retrying.
			// TODO: Should we have a retry limit maybe?
		}

		time.Sleep(time.Second * 60)
	}
}

// Run the agent as an instance, this is a single run of the agent that will check the
// policies against the plugins.
//
// Arguments:
// - plugins: list of plugin paths
// - policyBundles: list of policy bundle paths
// Returns:
// - error: any error that occurred during the run
func (ar *AgentRunner) runInstance(plugins []string, policyBundles []string) error {
	defer ar.closePluginClients()

	for _, source := range plugins {
		logger := hclog.New(&hclog.LoggerOptions{
			Name:   "runner",
			Output: os.Stdout,
			Level:  hclog.Debug,
		})

		pluginPath, err := ar.DownloadPlugin(source)
		if err != nil {
			return err
		}

		runnerInstance, err := ar.getRunnerInstance(logger, pluginPath)

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

func (ar *AgentRunner) getRunnerInstance(logger hclog.Logger, path string) (runner.Runner, error) {
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

// DownloadPlugin checks whether we need to download the source plugin, or whether it is already on the file system.
// If it isn't on the filesystem, we'll download it and return the final destination for use.
func (ar *AgentRunner) DownloadPlugin(source string) (usablePlugin string, err error) {
	// First we check if the source is a path that exists on the fs.
	// If it does exist, it means we've been passed a binary, and we can just use it as is.
	_, err = os.ReadFile(source)

	if err == nil {
		// The file exists. Just return it.
		return source, err
	}

	if !os.IsNotExist(err) {
		// The error we've received is something other than not exists.
		// Exit early with the error
		return "", err
	}

	loader, err := downloader.Download(source, AgentPluginDir)
	if err != nil {
		return "", err
	}

	dest, err := loader.GetFinalDestination()
	if err != nil {
		return "", err
	}

	// If we've downloaded a plugin, we should assume it's called "plugin" for the moment.
	// We may need to find a better method of deciding what is and is not the right thing to return later.
	// As an example, if we download a Python GRPC plugin, we will probably end up returning "main.py"
	return path.Join(dest, "plugin"), nil
}

func (ar *AgentRunner) closePluginClients() {
	plugin.CleanupClients()
}
