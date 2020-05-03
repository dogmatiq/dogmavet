package testhelper

import (
	"path/filepath"
	"sync/atomic"
	"testing"

	_ "github.com/dogmatiq/dogma" // ensure dogma is in go.mod, needed by Makefile
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/analysistest"
)

// Run runs tests for the given analyzer.
func Run(
	t *testing.T,
	a *analysis.Analyzer,
	patterns ...string,
) {
	atomic.AddInt32(&testCount, 1)
	defer atomic.AddInt32(&testCount, -1)

	analysistest.Run(
		t,
		goPath(),
		a,
		patterns...,
	)
}

// goPath returns the root of the GOPATH to use for testing anaylizers.
func goPath() string {
	p, err := filepath.Abs("../../artifacts/analysistest/gopath")
	if err != nil {
		panic(err)
	}

	return p
}
