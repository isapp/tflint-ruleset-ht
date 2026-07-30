// rules/atom_helpers_test.go
package rules

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

// The first case is the one that matters: it reproduces a real tflint invocation,
// where GetFiles() yields module-relative filenames with no tier segment. Before
// isAtomModule consulted the working directory, that case returned false and every
// atom rule silently disabled itself in production while the suite stayed green.
func TestIsAtomModule(t *testing.T) {
	cases := []struct {
		name  string
		dir   string // relative dir to chdir into, created under t.TempDir()
		files map[string]string
		want  bool
	}{
		{
			name:  "real invocation — atom cwd, module-relative filenames",
			dir:   filepath.Join("modules", "aws", "atoms", "eks-cluster"),
			files: map[string]string{"main.tf": `resource "aws_eks_cluster" "this" {}`},
			want:  true,
		},
		{
			name:  "real invocation — molecule cwd, module-relative filenames",
			dir:   filepath.Join("modules", "aws", "molecules", "s3-managed-bucket"),
			files: map[string]string{"main.tf": `resource "aws_s3_bucket" "this" {}`},
			want:  false,
		},
		{
			name: "test-harness shape — tier only in the filename",
			dir:  "somewhere-else",
			files: map[string]string{
				"modules/aws/atoms/s3-bucket/main.tf": `resource "aws_s3_bucket" "this" {}`,
			},
			want: true,
		},
		{
			name:  "neither signal present",
			dir:   "plain",
			files: map[string]string{"main.tf": `resource "aws_s3_bucket" "this" {}`},
			want:  false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			base := t.TempDir()
			dir := filepath.Join(base, tc.dir)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				t.Fatalf("mkdir: %s", err)
			}
			t.Chdir(dir)

			runner := helper.TestRunner(t, tc.files)
			got, err := isAtomModule(runner)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if got != tc.want {
				t.Errorf("isAtomModule() = %v, want %v (cwd=%q)", got, tc.want, dir)
			}
		})
	}
}
