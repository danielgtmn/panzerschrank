# AGENTS.md

Guidance for AI coding agents (and people) working on panzerschrank. User
documentation is the [README](README.md) and the package docs (shown on
pkg.go.dev).

## What this is

panzerschrank encrypts small secrets for storage in an app's database with
AES-256-GCM: the key outside the database, the secret's name as additional
authenticated data. It replaces the two near-identical copies in gpanel
([ownpanel](https://github.com/danielgtmn/ownpanel), `internal/secrets`) and
static-hosting (`backend/internal/secrets`), so the code exists once and is
tested once.

## Layout

```
panzerschrank.go       Box: New, NewFromBase64, Seal/Open, SealString/OpenString, GenerateKey
keyfile.go             OpenKeyFile, LoadOrCreateKey: the key in a file, created on first start
*_test.go              tests, the fixed test vector, FuzzOpen, examples (shown on pkg.go.dev)
scripts/               commit-type.sh (Conventional Commits check) and its test
```

## Rules that are easy to break

- **The format never changes.** Sealed data is `nonce(12) || ciphertext ||
  tag(16)`, the name is the AAD, base64 is `StdEncoding`. Apps have stored
  it. `TestKnownVector` and `TestCompatibleWithEarlierFormats` pin it;
  never edit their expected values. A new format would need a new function
  and a way to tell the formats apart, and a major version.
- **No secret material in errors or logs.** Every failure to open is
  `ErrDecrypt`, without detail. Errors about keys never quote the key.
- **Standard library only.** No dependencies, no `go.sum`; CI fails
  otherwise. No own cryptography: only `crypto/aes`, `crypto/cipher`,
  `crypto/rand`.
- **Nonces are random** from `crypto/rand`, fresh for every seal. Never
  accept a nonce from the caller.
- **The key file is never overwritten.** A key file of the wrong size is an
  error, not a reason to create a new key: that would lose every secret.
- Pin GitHub Actions to commit SHAs with a version comment. Apache-2.0: no
  GPL/AGPL code.

## Commands

```sh
make test   # go vet + go test -race
make fuzz   # FuzzOpen for a minute (FUZZTIME=5m)
make lint   # shellcheck, script tests, actionlint
```

## Done means tested

Every change comes with tests. For changes to the API, use the new version
in an app (gpanel or static-hosting, with a `replace` directive) and go
through storing and reading a secret there. If something could not be
tested end to end, say so explicitly.

## Git, pull requests, releases

- `main` is protected: pull requests only, merged with a **merge commit**.
  Required checks: `go`, `lint`, `commits`. CI skips jobs whose files did
  not change.
- **Conventional Commits** for every commit on `main` and every PR title
  (`feat:`, `fix:`, `perf:`, `docs:`, `chore:`, `ci:`, `refactor:`,
  `test:`, `build:`, optional scope and `!`). The `commits` check enforces it.
- Merging `main` into `release` tags a new version computed from the
  commits (`fix:` patch, `feat:` minor, `!` major) and publishes a GitHub
  release. Apps then `go get github.com/danielgtmn/panzerschrank@vX.Y.Z`.
