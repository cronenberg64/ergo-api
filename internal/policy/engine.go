package policy

import (
	"context"
	"embed"
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/rego"
)

//go:embed rbac.rego
var embeddedPolicies embed.FS

type Engine struct {
	query rego.PreparedEvalQuery
}

func NewEngine(ctx context.Context, policyPath string) (*Engine, error) {
	var policyContent []byte
	var err error

	// If policyPath is provided and exists, use it. Otherwise use embedded.
	if policyPath != "" {
		if _, err := os.Stat(policyPath); err == nil {
			policyContent, err = os.ReadFile(policyPath)
			if err != nil {
				return nil, fmt.Errorf("failed to read policy file: %w", err)
			}
		}
	}

	if len(policyContent) == 0 {
		// Fallback to embedded
		policyContent, err = embeddedPolicies.ReadFile("rbac.rego")
		if err != nil {
			return nil, fmt.Errorf("failed to read embedded policy: %w", err)
		}
		fmt.Println("📜 Using embedded RBAC policy")
	} else {
		fmt.Printf("📜 Using external policy: %s\n", policyPath)
	}

	r := rego.New(
		rego.Query("data.authz.allow"),
		rego.Module("rbac.rego", string(policyContent)),
	)

	query, err := r.PrepareForEval(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare rego query: %w", err)
	}

	return &Engine{query: query}, nil
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
