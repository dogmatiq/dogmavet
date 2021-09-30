# Dogma Vet Tool

[![Build Status](https://github.com/dogmatiq/dogmavet/workflows/CI/badge.svg)](https://github.com/dogmatiq/dogmavet/actions?workflow=CI)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/dogmavet/main.svg)](https://codecov.io/github/dogmatiq/dogmavet)
[![Latest Version](https://img.shields.io/github/tag/dogmatiq/dogmavet.svg?label=semver)](https://semver.org)
[![Documentation](https://img.shields.io/badge/go.dev-reference-007d9c)](https://pkg.go.dev/github.com/dogmatiq/dogmavet)
[![Go Report Card](https://goreportcard.com/badge/github.com/dogmatiq/dogmavet)](https://goreportcard.com/report/github.com/dogmatiq/dogmavet)

This repositoriy contains a custom Go [vet](https://golang.org/cmd/vet/) tool
that implements checkers for common mistakes in
[Dogma](https://github.com/dogmatiq/dogma) application and message handler
implementations.

## Installation

    go install github.com/dogmatiq/dogmavet/...

### Visual Studio Code

Assuming your `GOROOT` is in the default location, simply update the
`go.vetFlags` setting as follows:

```json
"go.vetFlags": [
  "-all",
  "-vettool=~/go/bin/dogmavet"
]
```

## Checks

The following checks are currently supported:

- Ensures `Configure()` methods call `Identity()` exactly once
- Ensures `Identity()` is called with valid names and keys
- Ensures identity keys are UUIDs, formatted as per RFC-4122

## Fixes

The checkers are able to provide fixes to common problems, however it seems that
`go vet` does not yet expose this information for consumptions by IDEs.

The fixers can be run manually using:

```
dogmavet -fix ./...
```

Please note that **all** of the following fixes will be applied:

- Rename the configurer parameter name to `c`
- Replace non-UUID identity keys with UUIDs
- Reformat non-standard UUID representations by the RFC-4122 grammar

## Caveats

The `go vet` command does not allow more than one `-vettool` argument to be
provided. Furthermore, setting `-vettool` disables all of the standard checks.

As a workaround, `dogmavet` includes all of the standard checks along with the
Dogma-specific ones.

See also [vscode-go/issues#3219](https://github.com/microsoft/vscode-go/issues/3219).
