package plugin

type Evaluator interface {
	PrepareForEval() error
	//Evaluate(query rego.PreparedEvalQuery) (rego.ResultSet, error)
}

type EvaluatorPlugin struct {
	// Impl Injection
	Impl Evaluator
}
