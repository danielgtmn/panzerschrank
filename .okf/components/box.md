---
type: Go Type
title: Box
description: Holds one AES-256-GCM key and seals and opens secrets bound to a name; safe for concurrent use.
resource: https://github.com/danielgtmn/panzerschrank/blob/main/panzerschrank.go
tags: [api, aes-gcm]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: code
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/panzerschrank.go
    title: panzerschrank.go
  - id: wiki-start
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/wiki/Getting-Started.md
    title: "Wiki: Getting started"
---

# Overview

A `Box` wraps a `cipher.AEAD` (AES-256-GCM). Create it once at startup and pass it to whatever stores secrets.[^wiki-start] The key is not kept beyond the cipher state.[^code]

# API

| Symbol | Purpose |
|---|---|
| `KeySize = 32` | key length in bytes |
| `New(key []byte)` | Box for a raw 32-byte key; other lengths are an error that names the length, never the key |
| `NewFromBase64(s)` | Box for a key in standard base64; surrounding whitespace ignored; bad base64 gives a generic error |
| `GenerateKey()` / `GenerateKeyBase64()` | fresh random key from `crypto/rand` |
| `(*Box).Seal(plaintext, name) []byte` | encrypt with a fresh random nonce; see [format](/concepts/sealed-format.md) |
| `(*Box).Open(sealed, name)` | decrypt; any failure is `ErrDecrypt` |
| `SealString` / `OpenString` | the same in standard base64 for text columns; invalid base64 is `ErrDecrypt` too |
| `(*Box).Overhead()` | 28 bytes (nonce + tag) |
| `ErrDecrypt` | the single error for wrong key, wrong name or damaged data |

For a Box backed by a file see [Key file](/components/key-file.md). The `name` argument is explained in [Name binding](/concepts/name-binding.md).

# Examples

```go
box, err := panzerschrank.NewFromBase64(os.Getenv("APP_SECRET_KEY"))
sealed := box.SealString(token, "dns-token:42")
token, err := box.OpenString(sealed, "dns-token:42")
if errors.Is(err, panzerschrank.ErrDecrypt) { /* log it, show no details */ }
```

[^code]: panzerschrank.go
[^wiki-start]: Wiki: Getting started
