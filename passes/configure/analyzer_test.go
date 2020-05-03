package configure_test

import (
	"testing"

	"github.com/dogmatiq/dogmavet/internal/testhelper"
	. "github.com/dogmatiq/dogmavet/passes/configure"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	analysistest.Run(
		t,
		testhelper.GoPath(),
		Analyzer,
		"github.com/dogmatiq/dogmavet/passes/configure/testdata/common",
		"github.com/dogmatiq/dogmavet/passes/configure/testdata/ignore",
	)
}
