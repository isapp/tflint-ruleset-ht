package rules_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/isapp/tflint-ruleset-ht/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestModuleSourceRule(t *testing.T) {
	rule := rules.NewModuleSourceRule()

	cases := []struct {
		name     string
		content  string
		expected helper.Issues
	}{
		{
			name: "valid HTTPS with ref — atom",
			content: `
module "topic" {
  source = "git::https://github.com/isapp/ht-terraform-modules?ref=modules/atoms/sns-topic/v1.0.0"
}`,
			expected: helper.Issues{},
		},
		{
			name: "valid HTTPS with ref — molecule",
			content: `
module "alert" {
  source = "git::https://github.com/isapp/ht-terraform-modules?ref=modules/molecules/cloudwatch-email-alarm/v1.1.0"
}`,
			expected: helper.Issues{},
		},
		{
			name: "ignored — local relative path",
			content: `
module "s3" {
  source = "../../../modules/aws/s3_bucket"
}`,
			expected: helper.Issues{},
		},
		{
			name: "ignored — terraform registry",
			content: `
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "5.0.0"
}`,
			expected: helper.Issues{},
		},
		{
			name: "invalid — SSH scheme",
			content: `
module "topic" {
  source = "git::ssh://git@github.com/isapp/ht-terraform-modules?ref=modules/atoms/sns-topic/v1.0.0"
}`,
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: "module source must use git::https://, not git::ssh://",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 9},
					},
				},
			},
		},
		{
			name: "invalid — //subdir notation",
			content: `
module "topic" {
  source = "git::https://github.com/isapp/ht-terraform-modules//modules/atoms/sns-topic?ref=modules/atoms/sns-topic/v1.0.0"
}`,
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: "module source must not use //subdir notation — ht-terraform-modules tags are flat commits",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 9},
					},
				},
			},
		},
		{
			name: "invalid — missing ?ref=",
			content: `
module "topic" {
  source = "git::https://github.com/isapp/ht-terraform-modules"
}`,
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: "module source must be pinned with ?ref=<tag>",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 9},
					},
				},
			},
		},
		{
			name: "invalid — all three violations",
			content: `
module "topic" {
  source = "git::ssh://git@github.com/isapp/ht-terraform-modules//modules/atoms/sns-topic"
}`,
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: "module source must use git::https://, not git::ssh://",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 9},
					},
				},
				{
					Rule:    rule,
					Message: "module source must not use //subdir notation — ht-terraform-modules tags are flat commits",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 9},
					},
				},
				{
					Rule:    rule,
					Message: "module source must be pinned with ?ref=<tag>",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 9},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": tc.content})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssues(t, tc.expected, runner.Issues)
		})
	}
}
