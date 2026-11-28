---
type: Integration
title: Apps using panzerschrank
description: Which apps use or are meant to use panzerschrank (gpanel, static-hosting, sesam) and how their existing data opens.
tags: [consumers, migration]
status: draft
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
stale_after: 2027-05-28
sources:
  - id: wiki-home
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/wiki/Home.md
    title: "Wiki: Home"
  - id: wiki-format
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/wiki/Format.md
    title: "Wiki: Format"
  - id: agents
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/AGENTS.md
    title: AGENTS.md
---

# Apps

| App | Earlier code | Status per the docs | Data |
|---|---|---|---|
| gpanel ([ownpanel](https://github.com/danielgtmn/ownpanel)) | `internal/secrets` | uses panzerschrank[^wiki-home] | same format and key file, opens as is |
| static-hosting | `backend/internal/secrets` | copy to be replaced[^agents] | two base64 columns + base64 key; join before `Open`[^wiki-format] |
| sesam | none | planned, for two-factor sign-in[^wiki-home] | new |

Status is as written in the repository docs on the generation date; check the apps themselves before relying on it. See [Stored format](/concepts/sealed-format.md) and [One shared package](/decisions/one-shared-package.md).

[^wiki-home]: Wiki: Home
[^wiki-format]: Wiki: Format
[^agents]: AGENTS.md
