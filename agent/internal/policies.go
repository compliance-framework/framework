package internal

import (
	"context"
	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"io/fs"
	"log"
	"strings"
)

func PolicyCompiler(ctx context.Context, policyPath string) *ast.Compiler {
	r := rego.New(
		rego.Query("data"),
		rego.Load([]string{policyPath}, func(abspath string, info fs.FileInfo, depth int) bool {
			// Exclude files that contain "_test.rego"
			return strings.Contains(abspath, "_test.rego")
		}),
	)

	query, err := r.PrepareForEval(ctx)
	if err != nil {
		log.Fatal(err)
	}

	compiler := ast.NewCompiler()
	compiler.Compile(query.Modules())

	return compiler
}
