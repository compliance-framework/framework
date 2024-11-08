package bundle

import (
	"fmt"
	"github.com/open-policy-agent/opa/ast"
	"strings"
)

type Violation map[string]interface{}

type Package string

func (p Package) PurePackage() string {
	return strings.TrimPrefix(string(p), "data.")
}

type Policy struct {
	File        string
	Package     Package
	Annotations []*ast.Annotations
}

type Result struct {
	Policy              Policy
	AdditionalVariables map[string]interface{}
	Violations          []Violation
}

func (res Result) String() string {
	return fmt.Sprintf(`
Policy:
	file: %s
	package: %s
	annotations: %s
AdditionalVariables: %v
Violations: %s
`, res.Policy.File, res.Policy.Package.PurePackage(), res.Policy.Annotations, res.AdditionalVariables, res.Violations)
}
