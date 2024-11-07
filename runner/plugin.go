package runner

import (
	goplugin "github.com/hashicorp/go-plugin"
	"net/rpc"
)

type ConfigItem struct {
	Key   string
	Value interface{}
}

type RunnerConfig []ConfigItem

func (rc RunnerConfig) findKey(key string) interface{} {
	for _, item := range rc {
		if item.Key == key {
			return item.Value
		}
	}
	return nil
}

func (rc RunnerConfig) GetString(key string) string {
	if found := rc.findKey(key); found != nil {
		return found.(string)
	}
	return ""
}

func (rc RunnerConfig) GetInt(key string) int {
	if found := rc.findKey(key); found != nil {
		return found.(int)
	}
	return 0
}

type Runner interface {
	Configure(RunnerConfig) error
	PrepareForEval() error
}

type RunnerRPC struct {
	client *rpc.Client
}

func (g *RunnerRPC) Configure(config RunnerConfig) error {
	var resp any
	err := g.client.Call("Plugin.Configure", config, &resp)
	return err
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

func (s *RunnerRPCServer) Configure(config RunnerConfig, resp *error) error {
	*resp = s.Impl.Configure(config)
	return nil
}

func (s *RunnerRPCServer) PrepareForEval(args interface{}, resp *error) error {
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
