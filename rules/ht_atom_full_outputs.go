// rules/ht_atom_full_outputs.go
package rules

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AtomFullOutputsRule enforces that atom-tier modules expose every provider-exported
// attribute of their resource as an output.
type AtomFullOutputsRule struct {
	tflint.DefaultRule
}

func NewAtomFullOutputsRule() *AtomFullOutputsRule {
	return &AtomFullOutputsRule{}
}

func (r *AtomFullOutputsRule) Name() string              { return "ht_atom_full_outputs" }
func (r *AtomFullOutputsRule) Enabled() bool             { return true }
func (r *AtomFullOutputsRule) Severity() tflint.Severity { return tflint.ERROR }
func (r *AtomFullOutputsRule) Link() string {
	return "https://github.com/isapp/tflint-ruleset-ht/blob/main/docs/rules/ht_atom_full_outputs.md"
}

func (r *AtomFullOutputsRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	isAtom := false
	for filename := range files {
		if strings.Contains(filepath.ToSlash(filename), "/atoms/") {
			isAtom = true
			break
		}
	}
	if !isAtom {
		return nil
	}

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
		return nil
	}
	resourceType := resourceContent.Blocks[0].Labels[0]

	schema, ok := resourceSchemas[resourceType]
	if !ok {
		return nil
	}

	outputContent, err := runner.GetModuleContent(
		&hclext.BodySchema{
			Blocks: []hclext.BlockSchema{
				{Type: "output", LabelNames: []string{"name"}, Body: &hclext.BodySchema{}},
			},
		},
		&tflint.GetModuleContentOption{},
	)
	if err != nil {
		return err
	}

	defined := make(map[string]bool, len(outputContent.Blocks))
	for _, block := range outputContent.Blocks {
		defined[block.Labels[0]] = true
	}

	for _, expected := range schema.outputs {
		if defined[expected] {
			continue
		}
		if err := runner.EmitIssue(
			r,
			fmt.Sprintf("atom is missing output %q (required by resource schema for %s)", expected, resourceType),
			hcl.Range{},
		); err != nil {
			return err
		}
	}
	return nil
}
