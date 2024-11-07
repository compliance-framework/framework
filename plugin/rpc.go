package plugin

import (
	"net/rpc"
)

type EvaluatorRPCClient struct {
	client *rpc.Client
}

func (g *EvaluatorRPCClient) PolicyNamespace() error {
	var resp string
	err := g.client.Call("Plugin.PolicyNamespace", new(interface{}), &resp)
	return err
}

func (g *EvaluatorRPCClient) PrepareForEval() error {
	var resp string
	err := g.client.Call("Plugin.PrepareForEval", new(interface{}), &resp)
	return err
}

//func (g *EvaluatorRPCClient) Evaluate(query rego.PreparedEvalQuery) (rego.ResultSet, error) {
//	var resp rego.ResultSet
//	err := g.client.Call("Plugin.Evaluate", query, &resp)
//	return resp, err
//}

type EvaluatorRPCServer struct {
	// This is the real implementation
	Impl Evaluator
}

func (s *EvaluatorRPCServer) PolicyNamespace() string {
	return s.Impl.PolicyNamespace()
}

func (s *EvaluatorRPCServer) PrepareForEval() error {
	err := s.Impl.PrepareForEval()
	return err
}

//func (s *EvaluatorRPCServer) Evaluate(query rego.PreparedEvalQuery, resp *rego.ResultSet) error {
//	v, err := s.Impl.Evaluate(query)
//	*resp = v
//
//	return err
//}
