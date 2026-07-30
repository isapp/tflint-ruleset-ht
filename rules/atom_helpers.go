// rules/atom_helpers.go
package rules

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// atomPathSegment is the directory that marks a module as atom-tier under
// ADR-0004's modules/{provider}/{tier}/{name} layout.
const atomPathSegment = "/atoms/"

// isAtomModule reports whether the module under inspection is atom-tier.
//
// Two signals, ORed, because neither alone is sufficient:
//
//   - The process working directory. tflint changes into each module directory
//     before running rules, for both `--recursive` and `--chdir`, so this is the
//     reliable signal in real runs.
//   - The filenames from runner.GetFiles(). In real runs these are
//     module-relative ("main.tf") and never contain the tier segment. The
//     plugin-sdk test helper, however, uses its map keys verbatim as filenames,
//     so unit tests supply full paths here.
//
// The atom rules originally used only the filename signal. That meant they
// returned early in every real invocation — both `tflint --recursive` from a repo
// root and the per-directory pre-commit hook — while their unit tests passed,
// because the test helper's map keys happen to carry the full path. Verified by
// planting a two-resource module under modules/aws/atoms/ and observing zero
// findings from either invocation.
func isAtomModule(runner tflint.Runner) (bool, error) {
	if wd, err := os.Getwd(); err == nil {
		if strings.Contains(filepath.ToSlash(wd), atomPathSegment) {
			return true, nil
		}
	}

	files, err := runner.GetFiles()
	if err != nil {
		return false, err
	}
	for filename := range files {
		if strings.Contains(filepath.ToSlash(filename), atomPathSegment) {
			return true, nil
		}
	}
	return false, nil
}
