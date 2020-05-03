// +build release

package report

import (
	"strings"

	"golang.org/x/tools/go/analysis"
)

// Ignore returns true if the given pass should be ignored.
func Ignore(pass *analysis.Pass) bool {
	// When built in release mode, don't perform any checks on the test
	// directory, which is full of things that fail the checks!
	return strings.Contains(pass.Pkg.Path(), "dogmatiq/dogmavet") &&
		strings.Contains(pass.Pkg.Path(), "testdata")
}
