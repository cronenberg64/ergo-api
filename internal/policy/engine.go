package policy

import (
	"context"
	"fmt"

	"github.com/open-policy-agent/opa/rego"
)

type Engine struct {
	query rego.PreparedEvalQuery
}

func NewEngine(ctx context.Context, policyPath string) (*Engine, error) {
	r := rego.New(
		rego.Query("data.rbac.allow"),
		rego.Load([]string{policyPath}, nil),
	)

	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare rego query: %w", err)
	}

	return &Engine{
		query: query,
	}, nil
}

func (e *Engine) Evaluate(ctx context.Context, input map[string]interface{}) (bool, error) {
	results, err := e.query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return false, err
	}

	if len(results) == 0 {
		return false, nil
	}

	if allowed, ok := results[0].Expressions[0].Value.(bool); ok {
		return allowed, nil
	}

	return false, nil
}
