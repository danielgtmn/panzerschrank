---
type: Decision
title: Random nonces only
description: Every seal uses a fresh 12-byte nonce from crypto/rand; the API never accepts a nonce from the caller.
tags: [crypto, api]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: agents
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/AGENTS.md
    title: AGENTS.md
  - id: code
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/panzerschrank.go
    title: panzerschrank.go
---

# Decision

`Seal` draws the nonce itself; sealing the same value twice gives different results. Never accept a nonce from the caller.[^agents]

# Why

Reusing a GCM nonce under one key breaks confidentiality and authenticity. With random 96-bit nonces a key should be rotated before about 2^32 seals ([Keys](/concepts/keys.md)). `GenerateKey` relies on `crypto/rand.Read` never failing: on a broken system it crashes instead of returning weak bytes.[^code]

[^agents]: AGENTS.md
[^code]: panzerschrank.go
