package plugin

import (
	"net/rpc"
)

type EvaluatorRPCClient struct {
	client *rpc.Client
}

func (g *EvaluatorRPCClient) PrepareForEval() error {
	var resp error
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
