package plugin

import (
	"context"
	"github.com/open-policy-agent/opa/rego"
)

type Plugin interface {
	PrepareForEval(ctx context.Context) error
	Evaluate(ctx context.Context, query rego.PreparedEvalQuery) (rego.ResultSet, error)
}
