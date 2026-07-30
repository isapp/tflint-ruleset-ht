// rules/ht_atom_full_inputs.go
package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AtomFullInputsRule enforces that atom-tier modules expose every configurable
// attribute of their resource as an input variable.
type AtomFullInputsRule struct {
	tflint.DefaultRule
}

func NewAtomFullInputsRule() *AtomFullInputsRule {
	return &AtomFullInputsRule{}
}

func (r *AtomFullInputsRule) Name() string              { return "ht_atom_full_inputs" }
func (r *AtomFullInputsRule) Enabled() bool             { return true }
func (r *AtomFullInputsRule) Severity() tflint.Severity { return tflint.ERROR }
func (r *AtomFullInputsRule) Link() string {
	return "https://github.com/isapp/tflint-ruleset-ht/blob/main/docs/rules/ht_atom_full_inputs.md"
}

func (r *AtomFullInputsRule) Check(runner tflint.Runner) error {
	isAtom, err := isAtomModule(runner)
	if err != nil {
		return err
	}
	if !isAtom {
		return nil
	}

	// Determine resource type from the single resource block.
	resourceContent, err := runner.GetModuleContent(
		&hclext.BodySchema{
			Blocks: []hclext.BlockSchema{
				{Type: "resource", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{}},
			},
		},
		&tflint.GetModuleContentOption{},
	)
	if err != nil {
		return err
	}
	if len(resourceContent.Blocks) != 1 {
		return nil // ht_atom_single_resource handles wrong counts
	}
	resourceType := resourceContent.Blocks[0].Labels[0]

	schema, ok := resourceSchemas[resourceType]
	if !ok {
		return nil // no schema entry yet — skip; add via scaffold-atom.py
	}

	varContent, err := runner.GetModuleContent(
		&hclext.BodySchema{
			Blocks: []hclext.BlockSchema{
				{Type: "variable", LabelNames: []string{"name"}, Body: &hclext.BodySchema{}},
			},
		},
		&tflint.GetModuleContentOption{},
	)
	if err != nil {
		return err
	}

	defined := make(map[string]bool, len(varContent.Blocks))
	for _, block := range varContent.Blocks {
		defined[block.Labels[0]] = true
	}

	for _, expected := range schema.inputs {
		if defined[expected] {
			continue
		}
		if err := runner.EmitIssue(
			r,
			fmt.Sprintf("atom is missing variable %q (required by resource schema for %s)", expected, resourceType),
			hcl.Range{},
		); err != nil {
			return err
		}
	}
	return nil
}
