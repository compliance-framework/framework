package runner

import (
	goplugin "github.com/hashicorp/go-plugin"
	"net/rpc"
)

type Runner interface {
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

type RunnerRPCServer struct {
	// This is the real implementation
	Impl Runner
}

func (s *RunnerRPCServer) Greet(args interface{}, resp *error) error {
	*resp = s.Impl.PrepareForEval()
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
