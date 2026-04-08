package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// VariableFieldOrderRule enforces that fields within each variable block appear in
// this order: type, default (if present), description.
type VariableFieldOrderRule struct {
	tflint.DefaultRule
}

func NewVariableFieldOrderRule() *VariableFieldOrderRule {
	return &VariableFieldOrderRule{}
}

func (r *VariableFieldOrderRule) Name() string              { return "isapp_variable_field_order" }
func (r *VariableFieldOrderRule) Enabled() bool             { return true }
func (r *VariableFieldOrderRule) Severity() tflint.Severity { return tflint.ERROR }
func (r *VariableFieldOrderRule) Link() string {
	return "https://github.com/isapp/tflint-ruleset-ht/blob/main/docs/rules/isapp_variable_field_order.md"
}

func (r *VariableFieldOrderRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(
		&hclext.BodySchema{
			Blocks: []hclext.BlockSchema{
				{
					Type:       "variable",
					LabelNames: []string{"name"},
					Body: &hclext.BodySchema{
						Attributes: []hclext.AttributeSchema{
							{Name: "type", Required: false},
							{Name: "default", Required: false},
							{Name: "description", Required: false},
						},
					},
				},
			},
		},
		&tflint.GetModuleContentOption{},
	)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		attrs := block.Body.Attributes
		typeAttr := attrs["type"]
		defaultAttr := attrs["default"]
		descAttr := attrs["description"]

		// type must come before default
		if typeAttr != nil && defaultAttr != nil {
			if typeAttr.Range.Start.Line > defaultAttr.Range.Start.Line {
				if err := runner.EmitIssue(r,
					fmt.Sprintf(`in variable "%s": "type" must appear before "default"`, block.Labels[0]),
					typeAttr.Range,
				); err != nil {
					return err
				}
			}
		}

		// type must come before description
		if typeAttr != nil && descAttr != nil {
			if typeAttr.Range.Start.Line > descAttr.Range.Start.Line {
				if err := runner.EmitIssue(r,
					fmt.Sprintf(`in variable "%s": "type" must appear before "description"`, block.Labels[0]),
					typeAttr.Range,
				); err != nil {
					return err
				}
			}
		}

		// default must come before description
		if defaultAttr != nil && descAttr != nil {
			if defaultAttr.Range.Start.Line > descAttr.Range.Start.Line {
				if err := runner.EmitIssue(r,
					fmt.Sprintf(`in variable "%s": "default" must appear before "description"`, block.Labels[0]),
					defaultAttr.Range,
				); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
