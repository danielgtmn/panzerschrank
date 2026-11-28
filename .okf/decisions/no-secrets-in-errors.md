---
type: Decision
title: No secret material in errors
description: Every failure to open is the single ErrDecrypt without detail, and errors about keys never quote the key.
tags: [security, errors]
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

- `Open`/`OpenString` return `ErrDecrypt` for every failure: too short, wrong key, wrong name, tampered, invalid base64.[^code]
- `New` names the wrong length, never the key. `NewFromBase64` replaces the decoder's message with a generic one, because the input string *is* the key.[^code]

# Why

Errors end up in logs and responses. No detail means no oracle and no leak.[^agents] Related: [Threat model](/concepts/threat-model.md).

[^agents]: AGENTS.md
[^code]: panzerschrank.go
