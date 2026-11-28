---
type: Playbook
title: Development
description: How panzerschrank is built and tested locally and in CI, and which tests pin behaviour that must not change.
tags: [testing, ci, fuzzing]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: wiki-dev
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/wiki/Development.md
    title: "Wiki: Development"
  - id: ci
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/.github/workflows/ci.yml
    title: ci.yml
  - id: agents
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/AGENTS.md
    title: AGENTS.md
---

# Commands

| Target | What it does |
|---|---|
| `make test` | `go vet` and `go test -race` |
| `make fuzz` | `FuzzOpen` for a minute (`FUZZTIME=5m`) |
| `make lint` | shellcheck, `scripts/commit-type_test.sh`, actionlint |

# Tests that matter

- `TestKnownVector`, `TestCompatibleWithEarlierFormats`: pin the [format](/concepts/sealed-format.md). Never change their expected values.
- `TestKeyFileConcurrentCreate`: many concurrent creations, 50 times; all get the same key, no temp file stays ([Key file](/components/key-file.md)).
- `FuzzOpen`: `Open` returns only `ErrDecrypt` or a value that seals back the same. Failing inputs go to `testdata/fuzz/`.
- `example_test.go`: shown on pkg.go.dev and run as tests.[^wiki-dev]

# CI (`ci.yml`)

A `changes` job decides what runs on pull requests; pushes to `main` run everything, and a change to `ci.yml` runs every job.[^ci]

| Job | Runs when | Steps |
|---|---|---|
| `go` | `*.go`, `go.mod`/`go.sum` changed | gofmt, vet, `test -race -count=1`, 30 s fuzz, [no dependencies](/decisions/stdlib-only.md) |
| `lint` | `*.sh` or workflows changed | shellcheck, script tests, actionlint |

# Done means tested

Every change comes with tests. API changes are tried in an app (gpanel or static-hosting, with a `replace` directive) by storing and reading a secret there; whatever could not be tested end to end is said explicitly.[^agents] Docs live in `wiki/` and are synced to the GitHub wiki; update them in the same PR as the behaviour.

[^wiki-dev]: Wiki: Development
[^ci]: ci.yml
[^agents]: AGENTS.md
