package runner

import (
	goplugin "github.com/hashicorp/go-plugin"
	"net/rpc"
)

type SmallConfig struct {
	smallKey string
}

type Config struct {
	Somekey  string
	Otherkey uint

	SmallPart SmallConfig
}

type Runner interface {
	Namespace() string
	Configure(Config) error
	PrepareForEval() error
}

type RunnerRPC struct {
	client *rpc.Client
}

func (g *RunnerRPC) PrepareForEval() error {
	var resp any
	err := g.client.Call("Plugin.PrepareForEval", new(interface{}), &resp)
	return err
}

func (g *RunnerRPC) Namespace() string {
	var resp string
	err := g.client.Call("Plugin.Namespace", new(interface{}), &resp)
	if err != nil {
		// Not much to do here right now
		panic(err)
	}
	return resp
}

func (g *RunnerRPC) Configure(config Config) error {
	var resp any
	err := g.client.Call("Plugin.Configure", config, &resp)
	return err
}

type RunnerRPCServer struct {
	// This is the real implementation
	Impl Runner
}

func (s *RunnerRPCServer) PrepareForEval(args interface{}, resp *error) error {
	*resp = s.Impl.PrepareForEval()
	return nil
}

func (s *RunnerRPCServer) Namespace(args interface{}, resp *string) error {
	*resp = s.Impl.Namespace()
	return nil
}

func (s *RunnerRPCServer) Configure(config Config, resp *error) error {
	*resp = s.Impl.Configure(config)
	return nil
}

type RunnerPlugin struct {
	// Impl Injection
	Impl Runner
}

func (p *RunnerPlugin) Server(*goplugin.MuxBroker) (interface{}, error) {
	return &RunnerRPCServer{Impl: p.Impl}, nil
}

func (RunnerPlugin) Client(b *goplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RunnerRPC{client: c}, nil
}
