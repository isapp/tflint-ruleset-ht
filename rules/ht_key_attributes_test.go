// rules/ht_key_attributes_test.go
package rules_test

import (
	"testing"

	"github.com/isapp/tflint-ruleset-ht/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// These cases pin the ordering contract for the EKS family: every key attribute
// first (contiguously), then the remainder A-Z.
//
// Note for atom authors: scaffold-atom.py only special-cases `name`, so for
// resources whose key attributes are named otherwise (aws_eks_addon,
// aws_eks_access_entry, aws_eks_node_group) the generated main.tf needs manual
// reordering before it satisfies this rule.
func TestKeyAttributesRule_EksFamily(t *testing.T) {
	rule := rules.NewKeyAttributesRule()

	cases := []struct {
		name      string
		files     map[string]string
		wantCount int
	}{
		{
			name: "eks-cluster — name first, remainder A-Z",
			files: map[string]string{
				"modules/aws/atoms/eks-cluster/main.tf": `
resource "aws_eks_cluster" "this" {
  name                          = var.name
  bootstrap_self_managed_addons = var.bootstrap_self_managed_addons
  deletion_protection           = var.deletion_protection
  region                        = var.region
  role_arn                      = var.role_arn
  tags                          = local.tags
  version                       = var.kubernetes_version
}`,
			},
			wantCount: 0,
		},
		{
			name: "eks-addon — both key attrs first, then A-Z",
			files: map[string]string{
				"modules/aws/atoms/eks-addon/main.tf": `
resource "aws_eks_addon" "this" {
  cluster_name  = var.cluster_name
  addon_name    = var.addon_name
  addon_version = var.addon_version
  preserve      = var.preserve
  tags          = local.tags
}`,
			},
			wantCount: 0,
		},
		{
			name: "eks-addon — non-key attr wedged between the key attrs",
			files: map[string]string{
				"modules/aws/atoms/eks-addon/main.tf": `
resource "aws_eks_addon" "this" {
  addon_name    = var.addon_name
  addon_version = var.addon_version
  cluster_name  = var.cluster_name
  tags          = local.tags
}`,
			},
			wantCount: 1,
		},
		{
			name: "eks-access-entry — both key attrs first, then A-Z",
			files: map[string]string{
				"modules/aws/atoms/eks-access-entry/main.tf": `
resource "aws_eks_access_entry" "this" {
  cluster_name      = var.cluster_name
  principal_arn     = var.principal_arn
  kubernetes_groups = var.kubernetes_groups
  region            = var.region
  tags              = local.tags
}`,
			},
			wantCount: 0,
		},
		{
			name: "eks-node-group — both key attrs first, then A-Z",
			files: map[string]string{
				"modules/aws/atoms/eks-node-group/main.tf": `
resource "aws_eks_node_group" "this" {
  cluster_name    = var.cluster_name
  node_group_name = var.node_group_name
  ami_type        = var.ami_type
  capacity_type   = var.capacity_type
  disk_size       = var.disk_size
  instance_types  = var.instance_types
  node_role_arn   = var.node_role_arn
  tags            = local.tags
}`,
			},
			wantCount: 0,
		},
		{
			name: "violation — non-key attribute before the key attribute",
			files: map[string]string{
				"modules/aws/atoms/eks-cluster/main.tf": `
resource "aws_eks_cluster" "this" {
  role_arn = var.role_arn
  name     = var.name
}`,
			},
			wantCount: 1,
		},
		{
			name: "violation — remainder not sorted A-Z",
			files: map[string]string{
				"modules/aws/atoms/eks-cluster/main.tf": `
resource "aws_eks_cluster" "this" {
  name     = var.name
  tags     = local.tags
  role_arn = var.role_arn
}`,
			},
			wantCount: 1,
		},
		{
			name: "unknown resource type is skipped",
			files: map[string]string{
				"modules/aws/atoms/whatever/main.tf": `
resource "aws_totally_unknown_thing" "this" {
  zebra = 1
  apple = 2
}`,
			},
			wantCount: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runner := helper.TestRunner(t, tc.files)
			if err := rule.Check(runner); err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if len(runner.Issues) != tc.wantCount {
				t.Errorf("got %d issue(s), want %d", len(runner.Issues), tc.wantCount)
				for _, iss := range runner.Issues {
					t.Logf("  - %s", iss.Message)
				}
			}
		})
	}
}
