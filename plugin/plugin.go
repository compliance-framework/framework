package plugin

import (
	goplugin "github.com/hashicorp/go-plugin"
	"github.com/open-policy-agent/opa/rego"
	"net/rpc"
)

type Evaluator interface {
	PrepareForEval() error
	Evaluate(query rego.PreparedEvalQuery) (rego.ResultSet, error)
}

type EvaluatorPlugin struct {
	// Impl Injection
	Impl Evaluator
}

func (p *EvaluatorPlugin) Server(*goplugin.MuxBroker) (interface{}, error) {
	return &EvaluatorRPCServer{Impl: p.Impl}, nil
}

func (EvaluatorPlugin) Client(b *goplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &EvaluatorRPCClient{client: c}, nil
}
