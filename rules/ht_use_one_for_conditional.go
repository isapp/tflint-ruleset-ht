package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

// skipRoots lists traversal roots that are NOT resource references and must not
// be flagged by this rule (they are never count-conditional resources).
var skipRoots = map[string]bool{
	"var":       true,
	"local":     true,
	"path":      true,
	"terraform": true,
	"each":      true,
	"count":     true,
	"self":      true,
	"module":    true,
}

// UseOneForConditionalRule flags resource_type.resource_name[0].attr traversals
// and suggests using one(resource_type.resource_name[*].attr) instead, which
// returns null safely when count = 0 rather than panicking.
type UseOneForConditionalRule struct {
	tflint.DefaultRule
}

// NewUseOneForConditionalRule creates a new UseOneForConditionalRule.
func NewUseOneForConditionalRule() *UseOneForConditionalRule {
	return &UseOneForConditionalRule{}
}

func (r *UseOneForConditionalRule) Name() string              { return "ht_use_one_for_conditional" }
func (r *UseOneForConditionalRule) Enabled() bool             { return true }
func (r *UseOneForConditionalRule) Severity() tflint.Severity { return tflint.ERROR }
func (r *UseOneForConditionalRule) Link() string {
	return "https://github.com/isapp/tflint-ruleset-ht/blob/main/docs/rules/ht_use_one_for_conditional.md"
}

// Check walks all HCL expression trees in every file and emits an issue for
// each ScopeTraversalExpr that contains a literal [0] index on a resource
// reference.
func (r *UseOneForConditionalRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		var visitErr error
		diags := hclsyntax.VisitAll(body, func(node hclsyntax.Node) hcl.Diagnostics {
			if visitErr != nil {
				return nil
			}

			expr, ok := node.(*hclsyntax.ScopeTraversalExpr)
			if !ok {
				return nil
			}

			traversal := expr.Traversal
			if len(traversal) < 2 {
				return nil
			}

			// The first step is always a TraverseRoot — get the root name.
			root, ok := traversal[0].(hcl.TraverseRoot)
			if !ok {
				return nil
			}
			if skipRoots[root.Name] {
				return nil
			}

			// Find the position of the [0] TraverseIndex step.
			zeroIdx := -1
			for i, step := range traversal {
				ti, ok := step.(hcl.TraverseIndex)
				if !ok {
					continue
				}
				if ti.Key.Type() == cty.Number {
					bf := ti.Key.AsBigFloat()
					if bf.IsInt() {
						iv, _ := bf.Int64()
						if iv == 0 {
							zeroIdx = i
							break
						}
					}
				}
			}
			if zeroIdx < 0 {
				return nil
			}

			// Build the resource reference prefix (everything before the [0]).
			prefixParts := make([]string, 0, zeroIdx)
			for _, step := range traversal[:zeroIdx] {
				switch s := step.(type) {
				case hcl.TraverseRoot:
					prefixParts = append(prefixParts, s.Name)
				case hcl.TraverseAttr:
					prefixParts = append(prefixParts, s.Name)
				}
			}
			prefix := strings.Join(prefixParts, ".")

			// Build the attribute suffix (everything after the [0]).
			var suffixParts []string
			for _, step := range traversal[zeroIdx+1:] {
				if ta, ok := step.(hcl.TraverseAttr); ok {
					suffixParts = append(suffixParts, ta.Name)
				}
			}
			suffix := strings.Join(suffixParts, ".")

			var msg string
			if suffix != "" {
				msg = fmt.Sprintf(
					`use one() instead of [0] indexing on %q — one(%s[*].%s) returns null safely when count = 0`,
					prefix, prefix, suffix,
				)
			} else {
				msg = fmt.Sprintf(
					`use one() instead of [0] indexing on %q — one(%s[*]) returns null safely when count = 0`,
					prefix, prefix,
				)
			}

			if err := runner.EmitIssue(r, msg, expr.SrcRange); err != nil {
				visitErr = err
			}
			return nil
		})
		if visitErr != nil {
			return visitErr
		}
		if diags.HasErrors() {
			return diags
		}
	}
	return nil
}
