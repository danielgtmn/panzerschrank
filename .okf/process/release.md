---
type: Playbook
title: Release
description: Changes reach main through pull requests with merge commits and Conventional Commits; versions are tagged by a manual Release workflow run on main.
tags: [release, git, conventional-commits]
generated: { by: claude-code/claude-opus-5-5, at: 2026-11-28T21:30:00Z }
sources:
  - id: agents
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/AGENTS.md
    title: AGENTS.md
  - id: release-yml
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/.github/workflows/release.yml
    title: release.yml
  - id: commits-yml
    resource: https://github.com/danielgtmn/panzerschrank/blob/main/.github/workflows/commits.yml
    title: commits.yml
---

# Pull requests

- `main` is protected: pull requests only, merged with a **merge commit**. Required checks: `go`, `lint`, `commits`.[^agents]
- Conventional Commits for every commit and every PR title (`feat:`, `fix:`, `perf:`, `docs:`, `chore:`, `ci:`, `refactor:`, `test:`, `build:`, optional scope and `!`). Because of merge commits, every commit in a PR must follow it; `commits.yml` checks the title via `scripts/commit-type.sh`.[^commits-yml]
- GitHub Actions are pinned to commit SHAs with a version comment. License Apache-2.0: no GPL/AGPL code.

# Steps to release

1. Make sure `main` is in a state worth keeping.
2. `gh workflow run release.yml --ref main` (or Actions → Release).
3. The tag job computes the next version from commits since the last `v*` tag (`fix:`/`perf:` patch, `feat:` minor, `!`/`BREAKING CHANGE:` major), pushes `v<version>` on `main`'s head; the release job publishes a GitHub release. Without such commits nothing is released.
4. Apps: `go get github.com/danielgtmn/panzerschrank@vX.Y.Z` (or `@main` for a pseudo-version).
5. If the release job fails after tagging: "Re-run failed jobs".[^release-yml]

A pushed version is permanent in the Go module proxy. Why: [Release by hand from main](/decisions/release-from-main.md).

[^agents]: AGENTS.md
[^release-yml]: release.yml
[^commits-yml]: commits.yml
