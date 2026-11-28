---
type: Decision
title: Release by hand from main
description: Releases are tagged by running the Release workflow on main by hand; the earlier release branch was dropped.
tags: [release, ci]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: commit-release-main
    resource: https://github.com/danielgtmn/panzerschrank/commit/4d70f15
    title: "ci: release by hand from main instead of a release branch"
    last_modified: 2026-11-28
  - id: commit-release-branch
    resource: https://github.com/danielgtmn/panzerschrank/commit/03cf4bf
    title: "ci: release only from the release branch"
    last_modified: 2026-11-25
  - id: release-yml
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/.github/workflows/release.yml
    title: release.yml
---

# Decision

`gh workflow run release.yml --ref main`. The workflow's tag job runs only on `refs/heads/main`.[^release-yml]

# History

1. First setup: releases tagged from a release branch.
2. 2026-11-25: restricted to the release branch, because a manual run could tag any branch and skip review and required checks.[^commit-release-branch]
3. 2026-11-28: release branch dropped; releases are made by hand from `main`, which is protected and only gets reviewed, checked PRs.[^commit-release-main]

# Why by hand

The Go module proxy keeps every published version forever; a version can only be retracted. So release on purpose, not on every merge. Process: [Release](/process/release.md).

[^commit-release-main]: ci: release by hand from main instead of a release branch
[^commit-release-branch]: ci: release only from the release branch
[^release-yml]: release.yml
