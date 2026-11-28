# AGENTS.md

Guidance for AI coding agents (and people) working on panzerschrank. User
documentation lives in [`wiki/`](wiki/) and is copied to the GitHub wiki by
a workflow: edit it there. The package docs are shown on pkg.go.dev.

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
wiki/                  documentation (synced to the GitHub wiki)
scripts/               commit-type.sh (Conventional Commits check) and its test
.okf/                  knowledge bundle (Open Knowledge Format): concepts, decisions, process
.claude/skills/        okf, okf-validate, okf-visualize: agent skills for .okf/ (copied, MIT)
```

Keep `.okf/` in step with the code: when a change touches what a concept
describes, update that concept (`generated`), its `index.md` and `log.md`
in the same PR, and check it with
`uv run .claude/skills/okf-validate/scripts/okf_validate.py .okf --strict`.

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
- Update the wiki (and README where it covers the same) in the same PR as
  the behaviour it describes.
- **Releases are made by hand from `main`:** run the Release workflow
  (`gh workflow run release.yml --ref main`). It tags the next version
  computed from the commits (`fix:` patch, `feat:` minor, `!` major) and
  publishes a GitHub release. Apps then `go get github.com/danielgtmn/panzerschrank@vX.Y.Z`.
  There is no release branch. A pushed version is permanent (the Go module
  proxy keeps it), so release on purpose, not on every merge.
