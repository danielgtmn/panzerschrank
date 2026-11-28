---
type: Playbook
title: Keys
description: Where the 32-byte key lives (key file or environment), how to back it up apart from the database, and how to rotate it by hand.
tags: [keys, backup, rotation, operations]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: wiki-keys
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/wiki/Keys.md
    title: "Wiki: Keys"
  - id: readme
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/README.md
    title: README
---

# Key file or environment

| | Key file (`OpenKeyFile`) | Environment (`NewFromBase64`) |
|---|---|---|
| Setup | none: created on first start, mode 0600 | once: `openssl rand -base64 32` |
| Good for | one server with a data directory (gpanel) | containers, platforms with a secret store |
| Backup | the file, apart from the database | the secret store |

Details of the file: [Key file](/components/key-file.md).[^wiki-keys]

# Backups

Back up the key **apart from the database**. Both together protect nothing; the database without the key cannot be restored.

# Rotation

No built-in rotation; for small numbers of secrets it is a loop in the app:

1. Open the old box, create a new one.
2. In one transaction: `Open` with the old box, `Seal` with the new one under the same name, write back.
3. After commit, switch the app to the new key; keep the old one until no backup needs it.

Rotate well before 2^32 seals per key (random 96-bit nonces), which small secrets do not reach in practice.[^readme]

[^wiki-keys]: Wiki: Keys
[^readme]: README
