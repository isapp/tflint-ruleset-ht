package rules_test

import (
	"testing"

	"github.com/isapp/tflint-ruleset-ht/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestUseOneForConditionalRule(t *testing.T) {
	rule := rules.NewUseOneForConditionalRule()

	cases := []struct {
		name     string
		content  string
		expected helper.Issues
	}{
		{
			name: "aws_instance[0].id in locals — error",
			content: `
locals {
  instance_id = aws_instance.web[0].id
}`,
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: `use one() instead of [0] indexing on "aws_instance.web" — one(aws_instance.web[*].id) returns null safely when count = 0`,
				},
			},
		},
		{
			name: "data.aws_route53_zone[0].zone_id in resource attribute — error",
			content: `
resource "aws_route53_record" "example" {
  zone_id = data.aws_route53_zone.this[0].zone_id
}`,
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: `use one() instead of [0] indexing on "data.aws_route53_zone.this" — one(data.aws_route53_zone.this[*].zone_id) returns null safely when count = 0`,
				},
			},
		},
		{
			name: "one(aws_instance.web[*].id) — pass",
			content: `
locals {
  instance_id = one(aws_instance.web[*].id)
}`,
			expected: helper.Issues{},
		},
		{
			name: "var.subnets[0] — pass (var root skipped)",
			content: `
locals {
  subnet = var.subnets[0]
}`,
			expected: helper.Issues{},
		},
		{
			name: "local.items[0] — pass (local root skipped)",
			content: `
locals {
  item = local.items[0]
}`,
			expected: helper.Issues{},
		},
		{
			name: "aws_instance.web[1].id — pass (index 1, not 0)",
			content: `
locals {
  instance_id = aws_instance.web[1].id
}`,
			expected: helper.Issues{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": tc.content})
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssuesWithoutRange(t, tc.expected, runner.Issues)
		})
	}
}
