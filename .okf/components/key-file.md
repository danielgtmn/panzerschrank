---
type: Go Function
title: Key file
description: OpenKeyFile and LoadOrCreateKey read the raw 32-byte key from a file and create it atomically with mode 0600 if it is missing; an existing file is never overwritten.
resource: https://github.com/danielgtmn/panzerschrank/blob/main/keyfile.go
tags: [api, keys, filesystem]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: code
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/keyfile.go
    title: keyfile.go
  - id: commit-keyfile
    resource: https://github.com/danielgtmn/panzerschrank/commit/d2e422f
    title: "fix(keyfile): never expose a partly written key file"
    last_modified: 2026-11-25
---

# Overview

| Symbol | Purpose |
|---|---|
| `KeyFile = "secret.key"` | suggested file name in the app's data directory |
| `OpenKeyFile(path)` | `LoadOrCreateKey` + `New` |
| `LoadOrCreateKey(path)` | read the key; if missing, create it; the directory must exist |

A file that exists with the wrong size is an **error**, never a reason to create a new key: that would lose every secret sealed with the old one.[^code] See [the decision](/decisions/key-file-never-overwritten.md).

# Creation

1. Write a new key to a temp file `.secret.key.tmp-*` (mode 0600) in the same directory and `fsync` it.
2. `os.Link` the temp file to the target: the link appears complete or not at all, and fails if the target exists.
3. `fsync` the directory (where the OS allows; not on Windows).
4. Read back whatever is on disk, so racing processes all end up with the winner's key.

Creating the file directly with `O_EXCL` let a second process read it empty and fail to start (27 of 300 test runs).[^commit-keyfile] `TestKeyFileConcurrentCreate` guards against this ([Development](/process/development.md)).

[^code]: keyfile.go
[^commit-keyfile]: fix(keyfile): never expose a partly written key file
