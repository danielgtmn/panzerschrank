# panzerschrank

panzerschrank encrypts small secrets for your Go app's database: API tokens, passwords for other services, secret environment variables. Instead of storing them in plain text, you store a sealed value; only your app, holding the key, can open it again.

- **The key stays out of the database.** A copy of the database alone, such as a backup or a dump, reveals nothing.
- **Every secret is bound to its name.** A sealed value copied into another row or field does not open there.
- **Changes are detected.** Changed data, another key or another name give an error, never wrong data.
- **Standard library only:** AES-256-GCM from Go's `crypto` packages, no other dependencies.

It is used by gpanel and planned for sesam's two-factor sign-in.

## Pages

- [[Getting Started|Getting-Started]]: add panzerschrank to an app
- [[Keys|Keys]]: where the key lives, backups, changing the key
- [[Format|Format]]: what is stored, and data from earlier copies
- [[Security|Security]]: what it protects and what is up to you
- [[Development|Development]]: build, test, release

API reference with examples: [pkg.go.dev](https://pkg.go.dev/github.com/danielgtmn/panzerschrank). The source is at <https://github.com/danielgtmn/panzerschrank>.
