package configure_test

import (
	"testing"

	"github.com/dogmatiq/dogmavet/internal/testhelper"
	. "github.com/dogmatiq/dogmavet/passes/configure"
)

func TestAnalyzer(t *testing.T) {
	testhelper.Run(
		t,
		Analyzer,
		"github.com/dogmatiq/dogmavet/passes/configure/testdata/common",
		"github.com/dogmatiq/dogmavet/passes/configure/testdata/ignore",
	)
}
