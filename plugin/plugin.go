package plugin

import (
	"github.com/hashicorp/go-plugin"
	"net/rpc"
)

type Evaluator interface {
	PrepareForEval() error
	//Evaluate(query rego.PreparedEvalQuery) (rego.ResultSet, error)
}

type EvaluatorRPC struct{ client *rpc.Client }

func (g *EvaluatorRPC) PrepareForEval() error {
	var resp string
	err := g.client.Call("Plugin.PrepareForEval", new(interface{}), &resp)
	return err
}

type EvaluatorRPCServer struct {
	// This is the real implementation
	Impl Evaluator
}

func (s *EvaluatorRPCServer) PrepareForEval(args interface{}, resp *string) error {
	err := s.Impl.PrepareForEval()
	return err
}

type EvaluatorPlugin struct {
	// Impl Injection
	Impl Evaluator
}

func (p *EvaluatorPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &EvaluatorRPCServer{Impl: p.Impl}, nil
}

func (EvaluatorPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &EvaluatorRPC{client: c}, nil
}
