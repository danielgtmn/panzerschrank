// Package panzerschrank encrypts small secrets (API tokens, passwords for
// other services, secret environment variables) for storage in a database,
// with AES-256-GCM.
//
// The key is kept outside the database (a file next to it, or an
// environment variable), so a copy of the database alone reveals nothing.
// Every secret is bound to a name, such as the column, the variable name or
// "dns-token:42", as additional authenticated data: the name is not stored,
// and a sealed value only opens under the name it was sealed with, so a
// secret cannot be copied into another row or field and used there.
//
//	box, err := panzerschrank.OpenKeyFile(filepath.Join(dataDir, panzerschrank.KeyFile))
//	sealed := box.Seal([]byte(token), "dns-token")
//	token, err := box.Open(sealed, "dns-token")
//
// Sealed data is nonce || ciphertext || tag (12 + len(plaintext) + 16
// bytes). SealString and OpenString wrap it in standard base64 for text
// columns and environment variables.
package panzerschrank

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

// KeySize is the length of a key in bytes (AES-256).
const KeySize = 32

// ErrDecrypt means sealed data could not be opened: it was changed, it was
// sealed with another key, or under another name. The message never
// contains secret material.
var ErrDecrypt = errors.New("panzerschrank: cannot decrypt (wrong key, wrong name or damaged data)")

// Box seals and opens secrets with one key. It is safe for concurrent use.
type Box struct {
	aead cipher.AEAD
}

// New returns a Box for a 32-byte key. The key is not retained beyond the
// cipher state; the caller may clear its copy afterwards.
func New(key []byte) (*Box, error) {
	if len(key) != KeySize {
		return nil, fmt.Errorf("panzerschrank: the key must be %d bytes, not %d", KeySize, len(key))
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("panzerschrank: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("panzerschrank: %w", err)
	}
	return &Box{aead: aead}, nil
}

// NewFromBase64 returns a Box for a key given as standard base64 (with
// padding), for example from an environment variable. Surrounding
// whitespace is ignored. Create such a key with GenerateKeyBase64 or
// `openssl rand -base64 32`.
func NewFromBase64(encoded string) (*Box, error) {
	key, err := base64.StdEncoding.DecodeString(strings.TrimSpace(encoded))
	if err != nil {
		// The decoder's message quotes the input position, not the key, but
		// stay generic: this string is the key.
		return nil, errors.New("panzerschrank: the key is not valid base64")
	}
	return New(key)
}

// GenerateKey returns a new random 32-byte key.
func GenerateKey() []byte {
	key := make([]byte, KeySize)
	// crypto/rand.Read never fails; on a broken system it crashes the
	// program instead of returning a weak key.
	rand.Read(key)
	return key
}

// GenerateKeyBase64 returns a new random key in standard base64, the format
// NewFromBase64 reads.
func GenerateKeyBase64() string {
	return base64.StdEncoding.EncodeToString(GenerateKey())
}

// Overhead is how many bytes Seal adds to the plaintext (nonce and tag).
func (b *Box) Overhead() int {
	return b.aead.NonceSize() + b.aead.Overhead()
}

// Seal encrypts plaintext with a fresh random nonce. name is authenticated
// but not stored: Open must be given the same name, so a secret sealed for
// one use cannot be passed off as another. Sealing the same plaintext twice
// gives different results.
func (b *Box) Seal(plaintext []byte, name string) []byte {
	n := b.aead.NonceSize()
	nonce := make([]byte, n, n+len(plaintext)+b.aead.Overhead())
	rand.Read(nonce)
	return b.aead.Seal(nonce, nonce, plaintext, []byte(name))
}

// Open decrypts what Seal returned. Any change to sealed, another key or
// another name yields ErrDecrypt.
func (b *Box) Open(sealed []byte, name string) ([]byte, error) {
	n := b.aead.NonceSize()
	if len(sealed) < n+b.aead.Overhead() {
		return nil, ErrDecrypt
	}
	out, err := b.aead.Open(nil, sealed[:n], sealed[n:], []byte(name))
	if err != nil {
		return nil, ErrDecrypt
	}
	return out, nil
}

// SealString is Seal for text: it returns the sealed value as standard
// base64, ready for a text column.
func (b *Box) SealString(plaintext, name string) string {
	return base64.StdEncoding.EncodeToString(b.Seal([]byte(plaintext), name))
}

// OpenString decrypts what SealString returned. Invalid base64 is reported
// as ErrDecrypt as well.
func (b *Box) OpenString(sealed, name string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(sealed)
	if err != nil {
		return "", ErrDecrypt
	}
	out, err := b.Open(raw, name)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
