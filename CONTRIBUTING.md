# Contributing to panzerschrank

Thanks for helping. panzerschrank is small on purpose, so a short heads-up before larger work saves everyone time.

## Before you start

- **Bugs and concrete features:** [issues](https://github.com/danielgtmn/panzerschrank/issues/new/choose). For anything bigger than a small fix, comment on the issue (or open one) first, so we can agree on the approach.
- **Security problems:** never in public, see [SECURITY.md](SECURITY.md).

## Setting up

You need Go (the version in `go.mod`). For `make lint`, shellcheck.

```sh
make test
make fuzz
make lint
```

## Making a change

1. Fork, then branch from `main`.
2. Read [AGENTS.md](AGENTS.md): the rules that are easy to break (the format never changes, no secret material in errors, standard library only).
3. Add tests with the change, and update the documentation in `wiki/` (and the README and package docs, where they cover the same thing) in the same pull request.
4. Run `make test` and `make lint`, and `make fuzz` for changes to opening data.

## Pull requests

- The title follows [Conventional Commits](https://www.conventionalcommits.org): `feat: …`, `fix: …`, `docs: …`, `chore: …`, with an optional scope like `fix(keyfile): …`. The same applies to every commit. The `commits` check enforces it, and release versions are computed from it.
- Reference the issue (`Closes #12`) and fill in the template.
- Pull requests are merged with a merge commit, so keep the commits readable.

## License

panzerschrank is licensed under the [Apache License 2.0](LICENSE). By contributing, you agree that your contribution is licensed under the same terms, and that you have the right to contribute it.

Everyone taking part follows the [Code of Conduct](CODE_OF_CONDUCT.md).
