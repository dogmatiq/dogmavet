<div align="center">

# Dogma Vet Tool

A custom Go [vet](https://golang.org/cmd/vet/) tool that checks for common
mistakes in [Dogma](https://github.com/dogmatiq/dogma) applications.

[![Latest Version](https://img.shields.io/github/tag/dogmatiq/dogmavet.svg?&style=for-the-badge&label=semver)](https://github.com/dogmatiq/dogmavet/releases)
[![Build Status](https://img.shields.io/github/actions/workflow/status/dogmatiq/dogmavet/ci.yml?style=for-the-badge&branch=main)](https://github.com/dogmatiq/dogmavet/actions/workflows/ci.yml)
[![Code Coverage](https://img.shields.io/codecov/c/github/dogmatiq/dogmavet/main.svg?style=for-the-badge)](https://codecov.io/github/dogmatiq/dogmavet)

</div>

## Installation

```
go install github.com/dogmatiq/dogmavet/cmd/dogmavet@latest
```

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
- Ensures identity keys are UUIDs, formatted as per RFC 4122

## Fixes

The checkers are able to provide fixes to common problems, however it seems that
`go vet` does not yet expose this information for consumption by IDEs.

The fixers can be run manually using:

```
dogmavet -fix ./...
```

Please note that **all** of the following fixes will be applied:

- Rename the configurer parameter name to `c`
- Replace non-UUID identity keys with UUIDs
- Reformat non-standard UUID representations using the RFC 4122 grammar

## Caveats

The `go vet` command does not allow more than one `-vettool` argument to be
provided. Furthermore, setting `-vettool` disables all of the standard checks.

See also [vscode-go/issues#3219](https://github.com/microsoft/vscode-go/issues/3219).
