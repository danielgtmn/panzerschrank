---
type: Decision
title: The key file is never overwritten
description: A key file of the wrong size is an error, never a reason to make a new key; a new file is created atomically via temp file and hard link.
tags: [keys, filesystem, durability]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: agents
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/AGENTS.md
    title: AGENTS.md
  - id: commit-keyfile
    resource: https://github.com/danielgtmn/panzerschrank/commit/d2e422f
    title: "fix(keyfile): never expose a partly written key file"
    last_modified: 2026-11-25
---

# Decision

`LoadOrCreateKey` only creates a key when the file does not exist. Wrong size is an error.[^agents]

# Why

Replacing the key would lose every secret sealed with it.

# How creation stays safe

The first version created the file with `O_EXCL` and then wrote the key. A process starting at the same time could read the empty file and fail (27 of 300 test runs). Now the key goes to a synced temp file that is hard-linked into place, and the directory is synced, so the file appears complete or not at all, also after a crash.[^commit-keyfile] Details: [Key file](/components/key-file.md).

[^agents]: AGENTS.md
[^commit-keyfile]: fix(keyfile): never expose a partly written key file
