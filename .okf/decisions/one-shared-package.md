---
type: Decision
title: One shared package
description: Replace the near-identical secrets code in gpanel and static-hosting with one tested package that opens their stored data unchanged.
tags: [history, scope]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: commit-feat
    resource: https://github.com/danielgtmn/panzerschrank/commit/7bd6404
    title: "feat: encrypt small secrets with AES-256-GCM"
    last_modified: 2026-11-25
  - id: agents
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/AGENTS.md
    title: AGENTS.md
---

# Decision

gpanel (`internal/secrets`) and static-hosting (`backend/internal/secrets`) each had their own copy of the same idea. panzerschrank is the one package for both, so the code exists and is tested once.[^agents]

# Consequences

- Data sealed by both earlier copies opens unchanged; static-hosting needs its two columns joined first ([Stored format](/concepts/sealed-format.md)).[^commit-feat]
- API changes are tried in a consuming app before release ([Apps](/integrations/consuming-apps.md)).

[^commit-feat]: feat: encrypt small secrets with AES-256-GCM
[^agents]: AGENTS.md
