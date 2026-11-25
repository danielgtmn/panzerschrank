# Format

`Seal` returns:

```
nonce (12 bytes) || ciphertext (as long as the plaintext) || tag (16 bytes)
```

AES-256-GCM with a fresh random nonce for every seal; the name is the additional authenticated data and is not stored. `SealString` returns the same bytes in standard base64 (with padding). `Box.Overhead()` is the 28 extra bytes.

The format never changes within a major version. A test vector in the repository pins it, so values stored today still open after an update.

## Data from earlier copies

panzerschrank replaces two copies that were built into apps before:

- **gpanel** (`internal/secrets`): the same format and the same key file. Values open as they are.
- **static-hosting** (`backend/internal/secrets`): the same cipher, but ciphertext and nonce in two base64 columns and the key in base64. Join them before opening:

```go
box, _ := panzerschrank.NewFromBase64(key)
ct, _ := base64.StdEncoding.DecodeString(row.ValueEnc)
nonce, _ := base64.StdEncoding.DecodeString(row.Nonce)
plain, err := box.Open(append(nonce, ct...), row.Name)
```
