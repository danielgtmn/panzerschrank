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

Pull requests go into `main` with a merge commit; titles and commits follow Conventional Commits. Merging `main` into `release` tags the next version from the commits (`fix:` patch, `feat:` minor, `!` major) and publishes a GitHub release. Apps then use `go get github.com/danielgtmn/panzerschrank@vX.Y.Z`.

## Conventions

See [AGENTS.md](https://github.com/danielgtmn/panzerschrank/blob/main/AGENTS.md): the format never changes, no secret material in errors, standard library only, nonces are always random, the key file is never overwritten.
