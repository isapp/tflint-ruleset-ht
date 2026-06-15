// rules/ht_atom_full_outputs_test.go
package rules_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/isapp/tflint-ruleset-ht/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestAtomFullOutputsRule(t *testing.T) {
	rule := rules.NewAtomFullOutputsRule()

	cases := []struct {
		name     string
		files    map[string]string
		expected helper.Issues
	}{
		{
			name: "valid — all expected outputs present",
			files: map[string]string{
				"modules/aws/atoms/s3-bucket-versioning/main.tf": `
resource "aws_s3_bucket_versioning" "this" {
  bucket = var.bucket
}`,
				"modules/aws/atoms/s3-bucket-versioning/outputs.tf": `
output "bucket" { value = aws_s3_bucket_versioning.this.bucket }
output "id" { value = aws_s3_bucket_versioning.this.id }
output "mfa" { value = aws_s3_bucket_versioning.this.mfa }
output "region" { value = aws_s3_bucket_versioning.this.region }
output "versioning_configuration" { value = aws_s3_bucket_versioning.this.versioning_configuration }`,
			},
			expected: helper.Issues{},
		},
		{
			name: "invalid — missing id output",
			files: map[string]string{
				"modules/aws/atoms/s3-bucket-versioning/main.tf": `
resource "aws_s3_bucket_versioning" "this" {
  bucket = var.bucket
}`,
				"modules/aws/atoms/s3-bucket-versioning/outputs.tf": `
output "bucket" { value = aws_s3_bucket_versioning.this.bucket }
output "mfa" { value = aws_s3_bucket_versioning.this.mfa }
output "region" { value = aws_s3_bucket_versioning.this.region }
output "versioning_configuration" { value = aws_s3_bucket_versioning.this.versioning_configuration }`,
			},
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: `atom is missing output "id" (required by resource schema for aws_s3_bucket_versioning)`,
					Range:   hcl.Range{},
				},
			},
		},
		{
			name: "skipped — resource type not in schema table",
			files: map[string]string{
				"modules/aws/atoms/unknown-resource/main.tf": `
resource "aws_unknown_thing" "this" {}`,
			},
			expected: helper.Issues{},
		},
		{
			name: "skipped — molecule path",
			files: map[string]string{
				"modules/aws/molecules/cloudwatch-email-alarm/main.tf": `
resource "aws_cloudwatch_metric_alarm" "this" {}`,
			},
			expected: helper.Issues{},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			helper.AssertIssues(t, tc.expected, runner.Issues)
		})
	}
}
