// +build !release

package report

import (
	"golang.org/x/tools/go/analysis"
)

// Ignore returns true if the given pass should be ignored.
func Ignore(pass *analysis.Pass) bool {
	return false
}
