---
type: Decision
title: The format is frozen
description: Sealed data, AAD and base64 encoding never change, because apps have stored them; test vectors pin it.
tags: [compatibility, format]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: agents
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/AGENTS.md
    title: AGENTS.md
---

# Decision

`nonce(12) || ciphertext || tag(16)`, the name as AAD, `StdEncoding` never change. `TestKnownVector` and `TestCompatibleWithEarlierFormats` pin it; never edit their expected values.[^agents]

# Why

Apps have stored values in this format. A change would make their secrets unreadable after an update.

# If a new format is ever needed

A new function, a way to tell the formats apart, and a major version. See [Stored format](/concepts/sealed-format.md).

[^agents]: AGENTS.md
