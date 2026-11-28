---
type: Concept
title: Name binding
description: Every secret is sealed under a name used as additional authenticated data; the name is not stored and a value only opens under the same name.
tags: [aad, design, usage]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: wiki-start
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/wiki/Getting-Started.md
    title: "Wiki: Getting started"
  - id: code
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/panzerschrank.go
    title: panzerschrank.go
---

# Overview

The `name` passed to `Seal`/`Open` is GCM additional authenticated data. It is not stored; opening under another name fails with `ErrDecrypt`. So a sealed value copied into another row or field does not open there.[^code]

# Choosing names

| Case | Name |
|---|---|
| build variable, names unique | `"API_TOKEN"` |
| secret per row | `"dns-token:" + strconv.Itoa(zoneID)` |
| secret per person | `"twofactor:" + userID` |

Keep names stable: changing how a name is built makes existing values unreadable.[^wiki-start] Related: [Box](/components/box.md), [Threat model](/concepts/threat-model.md).

[^wiki-start]: Wiki: Getting started
[^code]: panzerschrank.go
