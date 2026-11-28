---
type: Data Format
title: Stored format
description: Sealed data is nonce(12) || ciphertext || tag(16) with AES-256-GCM, the name as AAD and standard base64 for strings; apps have stored it, so it never changes.
tags: [format, compatibility, aes-gcm]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: wiki-format
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/wiki/Format.md
    title: "Wiki: Format"
  - id: agents
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/AGENTS.md
    title: AGENTS.md
  - id: tests
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/panzerschrank_test.go
    title: panzerschrank_test.go
---

# Schema

| Bytes | Field |
|---|---|
| 12 | nonce, random per seal |
| len(plaintext) | ciphertext |
| 16 | GCM tag |

- Cipher: AES-256-GCM (`crypto/aes`, `crypto/cipher`).
- AAD: the [name](/concepts/name-binding.md), not stored.
- `SealString`: the same bytes, `base64.StdEncoding` with padding.

# Compatibility

The format never changes within a major version. `TestKnownVector` and `TestCompatibleWithEarlierFormats` pin it; their expected values are never edited.[^agents][^tests] A new format would need a new function, a way to tell formats apart, and a major version. See [the decision](/decisions/format-is-frozen.md).

Data from the earlier copies:[^wiki-format]

- **gpanel**: same format and key file, opens as is.
- **static-hosting**: ciphertext and nonce in two base64 columns, key in base64. Join them first:

```go
box, _ := panzerschrank.NewFromBase64(key)
ct, _ := base64.StdEncoding.DecodeString(row.ValueEnc)
nonce, _ := base64.StdEncoding.DecodeString(row.Nonce)
plain, err := box.Open(append(nonce, ct...), row.Name)
```

[^wiki-format]: Wiki: Format
[^agents]: AGENTS.md
[^tests]: panzerschrank_test.go
