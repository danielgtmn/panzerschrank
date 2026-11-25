# Getting started

```sh
go get github.com/danielgtmn/panzerschrank@latest
```

## Open a box at startup

A `Box` holds the key and is safe for concurrent use. Create it once and pass it to whatever stores secrets.

With a key file next to the database, created on the first start:

```go
box, err := panzerschrank.OpenKeyFile(filepath.Join(dataDir, panzerschrank.KeyFile))
if err != nil {
    log.Fatal(err)
}
```

Or with a key from the environment, in base64:

```go
box, err := panzerschrank.NewFromBase64(os.Getenv("APP_SECRET_KEY"))
```

See [[Keys|Keys]] for which to choose.

## Seal and open

```go
// Store this in a text column.
sealed := box.SealString(token, "dns-token:42")

// Later:
token, err := box.OpenString(sealed, "dns-token:42")
if errors.Is(err, panzerschrank.ErrDecrypt) {
    // Wrong key, wrong name or damaged data. Log it; do not show details.
}
```

For binary columns use `Seal` and `Open`, which work on `[]byte`.

## Choosing the name

The name is not stored, but a sealed value only opens under the name it was sealed with. Pick a name that says what the secret is and where it belongs:

- `"API_TOKEN"` for a build variable, if names are unique.
- `"dns-token:" + strconv.Itoa(zoneID)` for a secret per row, so a value copied to another row does not open there.
- `"twofactor:" + userID` for a secret per person.

Keep names stable: changing how a name is built makes existing values unreadable.

## What not to use it for

Passwords that people sign in with are hashed (argon2id), not encrypted: they never need to be read back. The same goes for session and reset tokens (a SHA-256 hash is enough).
