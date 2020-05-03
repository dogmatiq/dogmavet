CGO_ENABLED = 1

# We setup a GOPATH style structure in the artifacts for testing statis
# analyzers because the analysistest package does not support modules yet.
GO_TEST_REQ += artifacts/analysistest/gopath/src/github.com/dogmatiq/dogma
GO_TEST_REQ += artifacts/analysistest/gopath/src/github.com/dogmatiq/dogmavet

-include .makefiles/Makefile
-include .makefiles/pkg/go/v1/Makefile

run: artifacts/build/debug/$(GOHOSTOS)/$(GOHOSTARCH)/dogma
	$< $(RUN_ARGS)

.makefiles/%:
	@curl -sfL https://makefiles.dev/v1 | bash /dev/stdin "$@"

artifacts/analysistest/gopath/src/%: $$(shell go list -f {{.Dir}} -m $$*)
	@mkdir -p "$(@D)"
	ln -s "$<" "$@"
