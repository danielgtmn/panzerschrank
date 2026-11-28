---
type: Decision
title: Standard library only
description: No module dependencies and no go.sum; cryptography only from crypto/aes, crypto/cipher and crypto/rand; CI fails otherwise.
tags: [dependencies, ci]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: agents
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/AGENTS.md
    title: AGENTS.md
  - id: ci
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/.github/workflows/ci.yml
    title: ci.yml
---

# Decision

No dependencies, no `go.sum`, no own cryptography: only `crypto/aes`, `crypto/cipher`, `crypto/rand`.[^agents] The `go` job fails if `go.sum` exists or `go.mod` has a `require`.[^ci]

# Why

A small security library should have nothing to audit beyond the standard library. Tools like actionlint run via `go run` in CI and `make lint`, not as module dependencies.

[^agents]: AGENTS.md
[^ci]: ci.yml
