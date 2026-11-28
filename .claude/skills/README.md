# Agent skills

Project skills for Claude Code (and other agents that read `.claude/skills/`).

## OKF: `okf`, `okf-validate`, `okf-visualize`

Write, check and draw the knowledge bundle in [`.okf/`](../../.okf/index.md).
Copied from [scaccogatto/okf-skills](https://github.com/scaccogatto/okf-skills)
0.6.0 (commit `cb5c0a9`), MIT licence in [`okf/LICENSE`](okf/LICENSE); the
vendored OKF spec in `okf/reference/SPEC.md` is Apache-2.0 (Google LLC), see
[`okf/NOTICE`](okf/NOTICE).

Local changes against upstream, so the names do not clash with other skills:

- `validate` is `okf-validate`, `visualize` is `okf-visualize` (directory and
  `name:`); `okf/SKILL.md` and `okf/scripts/okf_init.py` point at
  `okf-validate`.

The scripts need [uv](https://docs.astral.sh/uv/) (or Python 3.11+ with
PyYAML). They are tooling only and not part of the Go module.

To update: copy the three `skills/` directories from a newer upstream release
over these, apply the changes above again, and update the version here.
