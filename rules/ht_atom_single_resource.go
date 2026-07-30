// rules/ht_atom_single_resource.go
package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AtomSingleResourceRule enforces that atom-tier modules contain exactly one resource block.
type AtomSingleResourceRule struct {
	tflint.DefaultRule
}

func NewAtomSingleResourceRule() *AtomSingleResourceRule {
	return &AtomSingleResourceRule{}
}

func (r *AtomSingleResourceRule) Name() string              { return "ht_atom_single_resource" }
func (r *AtomSingleResourceRule) Enabled() bool             { return true }
func (r *AtomSingleResourceRule) Severity() tflint.Severity { return tflint.ERROR }
func (r *AtomSingleResourceRule) Link() string {
	return "https://github.com/isapp/tflint-ruleset-ht/blob/main/docs/rules/ht_atom_single_resource.md"
}

func (r *AtomSingleResourceRule) Check(runner tflint.Runner) error {
	isAtom, err := isAtomModule(runner)
	if err != nil {
		return err
	}
	if !isAtom {
		return nil
	}

	content, err := runner.GetModuleContent(
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

	if len(content.Blocks) == 1 {
		return nil
	}

	issueRange := hcl.Range{}
	if len(content.Blocks) > 0 {
		issueRange = content.Blocks[0].DefRange
	}

	return runner.EmitIssue(
		r,
		fmt.Sprintf("atom must contain exactly 1 resource block, found %d", len(content.Blocks)),
		issueRange,
	)
}
