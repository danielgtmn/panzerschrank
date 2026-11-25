# Security

## What panzerschrank protects against

- **A copy of the database** (backup, dump, a copied file, a leaked replica) without the key: the secrets cannot be read.
- **Changed values:** a modified, truncated or swapped value fails to open instead of returning wrong data.
- **Values moved between rows or fields:** each value is bound to its name and only opens under it.
- **Leaks through errors:** every failure to open is the same `ErrDecrypt` without details; errors about keys never contain the key.

## What it does not protect against

- **An attacker on the server** who can read the key file, the environment or the memory of the app can open everything. That is by design: the app must be able to read its secrets.
- **Losing the key.** There is no recovery. See [[Keys|Keys]].
- **Who may read a secret** inside the app: that is the app's decision.

## What the app must do

- Keep the key outside the database and back it up separately.
- Choose names that bind each value to where it belongs (see [[Getting Started|Getting-Started]]).
- Treat `ErrDecrypt` as an internal error: log it, show people a sentence about the effect ("The DNS token cannot be read; enter it again"), never the technical cause.
- Never log or return secrets after they were stored.

## Reporting a vulnerability

Privately, through [Security → Report a vulnerability](https://github.com/danielgtmn/panzerschrank/security/advisories/new). See [SECURITY.md](https://github.com/danielgtmn/panzerschrank/blob/main/SECURITY.md).
