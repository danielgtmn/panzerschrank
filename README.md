# panzerschrank

Encrypts small secrets for your Go app's database: API tokens, passwords
for other services, secret environment variables. AES-256-GCM from the
standard library, the key outside the database, and every secret bound to
its name.

```go
box, err := panzerschrank.OpenKeyFile(filepath.Join(dataDir, panzerschrank.KeyFile))
if err != nil { … }

sealed := box.SealString(token, "dns-token:42")       // store this
token, err := box.OpenString(sealed, "dns-token:42")  // and get it back
```

## What you get

- **The key stays out of the database.** `OpenKeyFile` keeps it in a file
  next to the database (created on first start, mode 0600). Or pass it in
  base64 from an environment variable with `NewFromBase64`. A copy of the
  database alone reveals nothing.
- **Secrets are bound to their name.** The name (a column, a variable name,
  `"dns-token:42"`) is authenticated with the secret but not stored. A
  sealed value copied into another row or field does not open there.
- **Changes are detected.** Changed data, another key or another name give
  `ErrDecrypt`, never garbage and never secret material in the message.
- **A fresh random nonce** for every seal; sealing the same value twice
  gives different results.
- **No dependencies** besides the standard library.

## Format

`Seal` returns `nonce (12 bytes) || ciphertext || tag (16 bytes)`;
`SealString` returns the same in standard base64. The format is fixed by a
test vector, so values stored today still open after an update.

Data sealed by the earlier copies in gpanel (`nonce || ciphertext`) opens
as is. Data from static-hosting (ciphertext and nonce in two base64
columns, base64 key) opens after joining the two:

```go
box, _ := panzerschrank.NewFromBase64(os.Getenv("SECRETS_KEY"))
ct, _ := base64.StdEncoding.DecodeString(row.ValueEnc)
nonce, _ := base64.StdEncoding.DecodeString(row.Nonce)
plain, err := box.Open(append(nonce, ct...), row.Name)
```

## Keep in mind

- **Back up the key, apart from the database.** Without it the secrets are
  gone; with both together, so is the protection.
- It protects a stolen or leaked database (backups, dumps, a copied file),
  not a compromised server: whoever can read the key file can open
  everything.
- One key encrypts many secrets with random nonces. Rotate it well before
  2³² seals (about 4 billion), which small secrets do not reach in practice.
- For passwords that people sign in with, use a password hash (argon2id),
  not this package: those never need to be decrypted.

## Install

```sh
go get github.com/danielgtmn/panzerschrank@latest
```

Documentation with examples: [pkg.go.dev](https://pkg.go.dev/github.com/danielgtmn/panzerschrank).

## Development

```sh
make test   # vet and tests with the race detector
make fuzz   # fuzz Open (FUZZTIME=5m)
make lint   # shellcheck, script tests, actionlint
```

See [AGENTS.md](AGENTS.md) and [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[Apache-2.0](LICENSE).
