package testhelper

import "path/filepath"

// GoPath returns the root of the GOPATH to use for testing anaylizers.
func GoPath() string {
	p, err := filepath.Abs("../../artifacts/analysistest/gopath")
	if err != nil {
		panic(err)
	}

	return p
}
