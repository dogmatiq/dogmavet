package main

import (
	"math/rand"
	"time"

	"github.com/dogmatiq/dogmavet/passes/configure"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	multichecker.Main(
		configure.Analyzer,
	)
}
