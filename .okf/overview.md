---
type: Library
title: panzerschrank
description: Go library that encrypts small secrets for an app's database with AES-256-GCM, the key outside the database and the secret's name as additional authenticated data.
resource: https://github.com/danielgtmn/panzerschrank
tags: [go, crypto, secrets, aes-gcm]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: readme
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/README.md
    title: README
  - id: agents
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/AGENTS.md
    title: AGENTS.md
  - id: commit-feat
    resource: https://github.com/danielgtmn/panzerschrank/commit/7bd6404
    title: "feat: encrypt small secrets with AES-256-GCM"
    last_modified: 2026-11-25
---

# Overview

panzerschrank seals small secrets (API tokens, passwords for other services, secret environment variables) so an app can store them in its database. A copy of the database alone reveals nothing, because the key lives outside it.[^readme] Module path: `github.com/danielgtmn/panzerschrank`, Go 1.26, standard library only.

It replaces two near-identical copies that lived in gpanel (`internal/secrets`) and static-hosting (`backend/internal/secrets`), so the code exists and is tested once.[^agents][^commit-feat]

# Parts

| Part | File | Concept |
|---|---|---|
| Box: seal and open | `panzerschrank.go` | [Box](/components/box.md) |
| Key file | `keyfile.go` | [Key file](/components/key-file.md) |
| Tests, vector, fuzzing, examples | `*_test.go` | [Development](/process/development.md) |
| User docs | `wiki/` (synced to the GitHub wiki) | [Development](/process/development.md) |
| Commit-type check | `scripts/commit-type.sh` | [Release](/process/release.md) |

# Key ideas

- [Stored format](/concepts/sealed-format.md): `nonce(12) || ciphertext || tag(16)`, frozen.
- [Name binding](/concepts/name-binding.md): the name is the AAD, not stored.
- [Keys](/concepts/keys.md): key file or base64 from the environment; back up apart from the database.
- [Threat model](/concepts/threat-model.md): protects a leaked database, not a compromised server.
- [Rules that are easy to break](/decisions/index.md).
- [Apps using it](/integrations/consuming-apps.md).

[^readme]: README
[^agents]: AGENTS.md
[^commit-feat]: feat: encrypt small secrets with AES-256-GCM
