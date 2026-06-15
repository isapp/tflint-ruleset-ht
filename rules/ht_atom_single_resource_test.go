// rules/ht_atom_single_resource_test.go
package rules_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/isapp/tflint-ruleset-ht/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestAtomSingleResourceRule(t *testing.T) {
	rule := rules.NewAtomSingleResourceRule()

	cases := []struct {
		name     string
		files    map[string]string
		expected helper.Issues
	}{
		{
			name: "valid — exactly one resource in atom path",
			files: map[string]string{
				"modules/aws/atoms/s3-bucket-versioning/main.tf": `
resource "aws_s3_bucket_versioning" "this" {
  bucket = var.bucket
}`,
			},
			expected: helper.Issues{},
		},
		{
			name: "invalid — two resource blocks in atom",
			files: map[string]string{
				"modules/aws/atoms/bad-atom/main.tf": `
resource "aws_s3_bucket" "this" {}
resource "aws_s3_bucket_versioning" "this" {}`,
			},
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: "atom must contain exactly 1 resource block, found 2",
					Range: hcl.Range{
						Filename: "modules/aws/atoms/bad-atom/main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 32},
					},
				},
			},
		},
		{
			name: "valid — multiple resources in molecule (not an atom)",
			files: map[string]string{
				"modules/aws/molecules/cloudwatch-email-alarm/main.tf": `
resource "aws_sns_topic" "this" {}
resource "aws_cloudwatch_metric_alarm" "this" {}`,
			},
			expected: helper.Issues{},
		},
		{
			name: "valid — multiple resources in non-module path",
			files: map[string]string{
				"main.tf": `
resource "aws_s3_bucket" "a" {}
resource "aws_s3_bucket" "b" {}`,
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
