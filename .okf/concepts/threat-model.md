---
type: Concept
title: Threat model
description: panzerschrank protects a leaked database, tampered values, values moved between rows and leaks through errors, but not a compromised server or a lost key.
tags: [security]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: wiki-security
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/wiki/Security.md
    title: "Wiki: Security"
  - id: security-md
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/SECURITY.md
    title: SECURITY.md
---

# Protects against

- A copy of the database without the key (backup, dump, copied file, leaked replica).
- Changed, truncated or swapped values: they fail to open.
- Values moved between rows or fields ([Name binding](/concepts/name-binding.md)).
- Leaks through errors: every failure is the same `ErrDecrypt` ([decision](/decisions/no-secrets-in-errors.md)).[^wiki-security]

# Does not protect against

- An attacker on the server who can read the key file, the environment or memory.
- Losing the key: no recovery ([Keys](/concepts/keys.md)).
- Who may read a secret inside the app.

# The app must

- Keep the key outside the database and back it up separately.
- Treat `ErrDecrypt` as an internal error: log it, show the effect, not the cause.
- Never log or return secrets after storing them.
- Not use it for sign-in passwords (hash with argon2id) or session tokens (SHA-256).

Vulnerabilities are reported privately through GitHub security advisories.[^security-md]

[^wiki-security]: Wiki: Security
[^security-md]: SECURITY.md
