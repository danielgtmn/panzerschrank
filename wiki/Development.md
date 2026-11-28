# Development

## Prerequisites

- **Go:** the version in `go.mod`.
- **shellcheck** for `make lint`.

## Make targets

| Target | What it does |
|---|---|
| `make test` | `go vet` and `go test -race` |
| `make fuzz` | fuzzes `Open` for a minute (`FUZZTIME=5m` for longer) |
| `make lint` | shellcheck, the commit message check's tests, actionlint |

CI runs the tests, 30 seconds of fuzzing, and checks that the module has no dependencies.

## Tests

- `TestKnownVector` and `TestCompatibleWithEarlierFormats` pin the format. Never change their expected values.
- `TestKeyFileConcurrentCreate` starts many key file creations at once, 50 times, and checks that all get the same key and no temporary file stays behind.
- `FuzzOpen` feeds random data to `Open`: it must only ever return `ErrDecrypt` or a value that seals back the same.
- The examples in `example_test.go` are shown on pkg.go.dev and run as tests.

## Releases

Pull requests go into `main` with a merge commit; titles and commits follow Conventional Commits. A release is made by hand from `main`: in the Actions tab, run **Release** (or `gh workflow run release.yml --ref main`). It tags the next version from the commits since the last tag (`fix:` patch, `feat:` minor, `!` major) on `main`'s head and publishes a GitHub release; without such commits it releases nothing. Apps then use `go get github.com/danielgtmn/panzerschrank@vX.Y.Z`, or `@main` for the latest commit.

A Go version is just a git tag, and the Go module proxy keeps every version it has seen for good: a published version cannot be changed or deleted, only marked as retracted in `go.mod`. So release on purpose, when `main` is in a state worth keeping.

## Conventions

See [AGENTS.md](https://github.com/danielgtmn/panzerschrank/blob/main/AGENTS.md): the format never changes, no secret material in errors, standard library only, nonces are always random, the key file is never overwritten.
