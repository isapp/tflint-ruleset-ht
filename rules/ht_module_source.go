package rules

import (
	"strings"

	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

const htModulesRepo = "github.com/isapp/ht-terraform-modules"

// ModuleSourceRule enforces source conventions for ht-terraform-modules references.
type ModuleSourceRule struct {
	tflint.DefaultRule
}

// NewModuleSourceRule creates a new ModuleSourceRule.
func NewModuleSourceRule() *ModuleSourceRule {
	return &ModuleSourceRule{}
}

func (r *ModuleSourceRule) Name() string              { return "ht_module_source" }
func (r *ModuleSourceRule) Enabled() bool             { return true }
func (r *ModuleSourceRule) Severity() tflint.Severity { return tflint.ERROR }
func (r *ModuleSourceRule) Link() string {
	return "https://github.com/isapp/tflint-ruleset-ht/blob/main/docs/rules/ht_module_source.md"
}

// Check walks all module blocks and validates source values that reference ht-terraform-modules.
func (r *ModuleSourceRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		for _, block := range body.Blocks {
			if block.Type != "module" || len(block.Labels) < 1 {
				continue
			}

			sourceAttr, exists := block.Body.Attributes["source"]
			if !exists {
				continue
			}

			val, diags := sourceAttr.Expr.Value(nil)
			if diags.HasErrors() || val.Type() != cty.String {
				continue
			}
			source := val.AsString()

			if !strings.Contains(source, htModulesRepo) {
				continue
			}

			if strings.Contains(source, "git::ssh://") {
				if err := runner.EmitIssue(
					r,
					"module source must use git::https://, not git::ssh://",
					sourceAttr.NameRange,
				); err != nil {
					return err
				}
			}

			// Strip the scheme (e.g. "https://") before checking for //subdir,
			// otherwise the :// in the URL scheme itself triggers a false positive.
			withoutScheme := source
			if idx := strings.Index(source, "://"); idx >= 0 {
				withoutScheme = source[idx+3:]
			}
			if strings.Contains(withoutScheme, "//") {
				if err := runner.EmitIssue(
					r,
					"module source must not use //subdir notation — ht-terraform-modules tags are flat commits",
					sourceAttr.NameRange,
				); err != nil {
					return err
				}
			}

			if !strings.Contains(source, "?ref=") {
				if err := runner.EmitIssue(
					r,
					"module source must be pinned with ?ref=<tag>",
					sourceAttr.NameRange,
				); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
