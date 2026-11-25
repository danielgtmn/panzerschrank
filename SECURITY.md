# Security policy

panzerschrank protects secrets that apps store in their databases, so security reports are welcome and taken seriously.

## Reporting a vulnerability

**Please do not open a public issue.** Report it privately instead:

1. Go to [Security → Report a vulnerability](https://github.com/danielgtmn/panzerschrank/security/advisories/new).
2. Describe what is affected, how to reproduce it (ideally a Go test), and what an attacker gains. Add the version from your `go.mod`. Never include real keys or secrets.

You get an answer within a week. Once a fix is released, the advisory is published, and you are credited unless you prefer not to be.

## Supported versions

panzerschrank is before 1.0. Fixes go into `main` and into the next release, which is cut promptly for security fixes. Only the **latest release** receives fixes.

## Scope

In scope: the package in this repository.

Out of scope: vulnerabilities in Go's own cryptography (report them upstream), attackers who can read the key file or the memory of the running app (the key is meant to be readable by the app), and apps using panzerschrank, unless its documented use leads to the problem. What it protects against is described under "Keep in mind" in the [README](README.md).
