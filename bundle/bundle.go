package bundle

import (
	"context"
	"fmt"
	"github.com/open-policy-agent/opa/rego"
	"slices"
	"strings"
)

type Bundle struct {
	path  string
	query rego.PreparedEvalQuery
}

func New(ctx context.Context, policyPath string) *Bundle {
	return &Bundle{path: policyPath}
}

func (b *Bundle) BuildQuery(ctx context.Context, pluginNamespace string) (*Bundle, error) {
	r := rego.New(
		rego.Query("data.compliance_framework"),
		rego.LoadBundle(b.path),
		rego.Package(fmt.Sprintf("compliance_framework.%s", pluginNamespace)),
	)

	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return b, err
	}

	b.query = query

	// Check that it will be able to prepare when we're ready to run
	return b, nil
}

func (b *Bundle) Execute(ctx context.Context, input map[string]interface{}) ([]Result, error) {
	var output []Result

	for _, module := range b.query.Modules() {
		// Exclude any test files for this compilation
		if strings.HasSuffix(module.Package.Location.File, "_test.rego") {
			continue
		}

		result := Result{
			Policy: Policy{
				File:        module.Package.Location.File,
				Package:     Package(module.Package.Path.String()),
				Annotations: module.Annotations,
			},
			AdditionalVariables: map[string]interface{}{},
			Violations:          nil,
		}

		sub := rego.New(
			rego.Query(module.Package.Path.String()),
			rego.LoadBundle(b.path),
			rego.Package(module.Package.Path.String()),
			rego.Input(input),
		)

		evaluation, err := sub.Eval(ctx)
		if err != nil {
			return nil, err
		}

		for _, eval := range evaluation {
			for _, expression := range eval.Expressions {
				moduleOutputs := expression.Value.(map[string]interface{})

				for key, value := range moduleOutputs {
					if !slices.Contains([]string{"violation"}, key) {
						result.AdditionalVariables[key] = value
					}
				}

				for _, tester := range moduleOutputs["violation"].([]interface{}) {
					result.Violations = append(result.Violations, tester.(map[string]interface{}))
				}

			}
		}
		output = append(output, result)
	}

	//compiler
	return output, nil
}
