// rules/ht_atom_full_inputs_test.go
package rules_test

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/isapp/tflint-ruleset-ht/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func TestAtomFullInputsRule(t *testing.T) {
	rule := rules.NewAtomFullInputsRule()

	cases := []struct {
		name     string
		files    map[string]string
		expected helper.Issues
	}{
		{
			name: "valid — all expected inputs present",
			files: map[string]string{
				"modules/aws/atoms/s3-bucket-versioning/main.tf": `
resource "aws_s3_bucket_versioning" "this" {
  bucket = var.bucket
}`,
				"modules/aws/atoms/s3-bucket-versioning/inputs.tf": `
variable "bucket" {}
variable "expected_bucket_owner" {}
variable "id" {}
variable "mfa" {}
variable "region" {}
variable "versioning_configuration" {}`,
			},
			expected: helper.Issues{},
		},
		{
			name: "invalid — missing versioning_configuration variable",
			files: map[string]string{
				"modules/aws/atoms/s3-bucket-versioning/main.tf": `
resource "aws_s3_bucket_versioning" "this" {
  bucket = var.bucket
}`,
				"modules/aws/atoms/s3-bucket-versioning/inputs.tf": `
variable "bucket" {}
variable "expected_bucket_owner" {}
variable "id" {}
variable "mfa" {}
variable "region" {}`,
			},
			expected: helper.Issues{
				{
					Rule:    rule,
					Message: `atom is missing variable "versioning_configuration" (required by resource schema for aws_s3_bucket_versioning)`,
					Range:   hcl.Range{},
				},
			},
		},
		{
			name: "skipped — resource type not in schema table",
			files: map[string]string{
				"modules/aws/atoms/unknown-resource/main.tf": `
resource "aws_unknown_thing" "this" {}`,
				"modules/aws/atoms/unknown-resource/inputs.tf": `
variable "id" {}`,
			},
			expected: helper.Issues{},
		},
		{
			name: "skipped — molecule path",
			files: map[string]string{
				"modules/aws/molecules/cloudwatch-email-alarm/main.tf": `
resource "aws_sns_topic" "this" {}`,
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
