package testhelper

import (
	"strings"
	"sync/atomic"

	"golang.org/x/tools/go/analysis"
)

// testCount is a count of the currently running tests.
var testCount int32 // testCount

// Ignore returns true if the given pass should be ignored.
func Ignore(pass *analysis.Pass) bool {
	if atomic.LoadInt32(&testCount) > 0 {
		return false
	}

	// Don't perform any checks on the testdata directory when we're not running
	// a test - it's full of things that fail the checks!
	return strings.Contains(pass.Pkg.Path(), "dogmatiq/dogmavet") &&
		strings.Contains(pass.Pkg.Path(), "testdata")
}
