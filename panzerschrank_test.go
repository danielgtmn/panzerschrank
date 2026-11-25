package panzerschrank

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"testing"
)

func newBox(t *testing.T) *Box {
	t.Helper()
	b, err := New(GenerateKey())
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestRoundTrip(t *testing.T) {
	b := newBox(t)
	for _, plain := range []string{"", "hetzner-token-123", strings.Repeat("x", 5000)} {
		sealed := b.Seal([]byte(plain), "dns")
		if len(sealed) != len(plain)+b.Overhead() {
			t.Fatalf("sealed length %d, want %d", len(sealed), len(plain)+b.Overhead())
		}
		if plain != "" && bytes.Contains(sealed, []byte(plain)) {
			t.Fatal("sealed data contains the plaintext")
		}
		got, err := b.Open(sealed, "dns")
		if err != nil || string(got) != plain {
			t.Fatalf("got %q, %v", got, err)
		}
	}
	// A fresh nonce every time.
	if bytes.Equal(b.Seal([]byte("same"), "dns"), b.Seal([]byte("same"), "dns")) {
		t.Fatal("two seals of the same plaintext are equal")
	}
}

func TestTamperDetection(t *testing.T) {
	b := newBox(t)
	sealed := b.Seal([]byte("token"), "dns")
	for i := range sealed {
		bad := bytes.Clone(sealed)
		bad[i] ^= 0x01
		if _, err := b.Open(bad, "dns"); !errors.Is(err, ErrDecrypt) {
			t.Fatalf("byte %d flipped: err = %v", i, err)
		}
	}
	for _, bad := range [][]byte{nil, sealed[:len(sealed)-1], sealed[:b.Overhead()-1], append(bytes.Clone(sealed), 0)} {
		if _, err := b.Open(bad, "dns"); !errors.Is(err, ErrDecrypt) {
			t.Fatalf("%d bytes: err = %v", len(bad), err)
		}
	}
}

func TestNameIsBound(t *testing.T) {
	b := newBox(t)
	sealed := b.Seal([]byte("token"), "API_TOKEN")
	for _, other := range []string{"", "api_token", "API_TOKEN ", "DB_PASSWORD"} {
		if _, err := b.Open(sealed, other); !errors.Is(err, ErrDecrypt) {
			t.Fatalf("opened under name %q: %v", other, err)
		}
	}
}

func TestWrongKey(t *testing.T) {
	sealed := newBox(t).Seal([]byte("token"), "dns")
	_, err := newBox(t).Open(sealed, "dns")
	if !errors.Is(err, ErrDecrypt) {
		t.Fatalf("err = %v", err)
	}
	if strings.Contains(err.Error(), "token") {
		t.Fatal("the error reveals the plaintext")
	}
}

func TestNewRefusesBadKeys(t *testing.T) {
	for _, n := range []int{0, 16, 24, 31, 33, 64} {
		if _, err := New(make([]byte, n)); err == nil {
			t.Fatalf("a %d-byte key must be refused", n)
		}
	}
}

func TestNewFromBase64(t *testing.T) {
	enc := GenerateKeyBase64()
	b1, err := NewFromBase64(enc)
	if err != nil {
		t.Fatal(err)
	}
	// Trailing newline from a file or `openssl rand -base64 32`.
	b2, err := NewFromBase64("  " + enc + "\n")
	if err != nil {
		t.Fatal(err)
	}
	if got, err := b2.Open(b1.Seal([]byte("x"), "p"), "p"); err != nil || string(got) != "x" {
		t.Fatalf("same key gives another box: %q %v", got, err)
	}

	short := base64.StdEncoding.EncodeToString(make([]byte, 16))
	secretish := "not-base64-SECRET!"
	for _, bad := range []string{"", short, secretish, enc[:len(enc)-1]} {
		_, err := NewFromBase64(bad)
		if err == nil {
			t.Fatalf("key %q must be refused", bad)
		}
		if bad != "" && strings.Contains(err.Error(), bad) {
			t.Fatalf("the error reveals the key: %v", err)
		}
	}
}

func TestStrings(t *testing.T) {
	b := newBox(t)
	sealed := b.SealString("s3cret", "DB_PASSWORD")
	if _, err := base64.StdEncoding.DecodeString(sealed); err != nil {
		t.Fatalf("not base64: %q", sealed)
	}
	got, err := b.OpenString(sealed, "DB_PASSWORD")
	if err != nil || got != "s3cret" {
		t.Fatalf("got %q, %v", got, err)
	}
	for _, bad := range []string{"", "%%%", sealed[:len(sealed)-4]} {
		if _, err := b.OpenString(bad, "DB_PASSWORD"); !errors.Is(err, ErrDecrypt) {
			t.Fatalf("%q: err = %v", bad, err)
		}
	}
	if _, err := b.OpenString(sealed, "OTHER"); !errors.Is(err, ErrDecrypt) {
		t.Fatalf("other name: %v", err)
	}
}

// Data sealed by the earlier copies in gpanel (nonce || ciphertext, purpose
// as AAD) and static-hosting (ciphertext and nonce stored apart, variable
// name as AAD) must open. Both used plain AES-256-GCM with a 12-byte nonce.
func TestCompatibleWithEarlierFormats(t *testing.T) {
	key, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	nonce, _ := hex.DecodeString("a0a1a2a3a4a5a6a7a8a9aaab")
	block, _ := aes.NewCipher(key)
	aead, _ := cipher.NewGCM(block)
	ct := aead.Seal(nil, nonce, []byte("legacy"), []byte("API_TOKEN"))

	b, err := New(key)
	if err != nil {
		t.Fatal(err)
	}
	// gpanel: one blob, nonce first.
	if got, err := b.Open(append(bytes.Clone(nonce), ct...), "API_TOKEN"); err != nil || string(got) != "legacy" {
		t.Fatalf("gpanel format: %q %v", got, err)
	}
	// static-hosting: base64 key, base64 nonce and ciphertext in two columns.
	b64, err := NewFromBase64(base64.StdEncoding.EncodeToString(key))
	if err != nil {
		t.Fatal(err)
	}
	ctCol := base64.StdEncoding.EncodeToString(ct)
	nonceCol := base64.StdEncoding.EncodeToString(nonce)
	rawCt, _ := base64.StdEncoding.DecodeString(ctCol)
	rawNonce, _ := base64.StdEncoding.DecodeString(nonceCol)
	if got, err := b64.Open(append(rawNonce, rawCt...), "API_TOKEN"); err != nil || string(got) != "legacy" {
		t.Fatalf("static-hosting format: %q %v", got, err)
	}
}

// A fixed vector pins the format: if this changes, stored secrets no longer
// open after an update.
func TestKnownVector(t *testing.T) {
	key, _ := hex.DecodeString("000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f")
	sealed, _ := hex.DecodeString("a0a1a2a3a4a5a6a7a8a9aaab" + knownCiphertext)
	b, _ := New(key)
	got, err := b.Open(sealed, "API_TOKEN")
	if err != nil || string(got) != "legacy" {
		t.Fatalf("got %q, %v", got, err)
	}
}

// AES-256-GCM of "legacy" under the key and nonce above with AAD
// "API_TOKEN", computed once with the standard library.
const knownCiphertext = "8a7d1b4c26b2c40de994f32e82c9c1a995d51d1227d3"

func FuzzOpen(f *testing.F) {
	b, _ := New(make([]byte, KeySize))
	f.Add(b.Seal([]byte("token"), "n"), "n")
	f.Add([]byte{}, "")
	f.Fuzz(func(t *testing.T, sealed []byte, name string) {
		out, err := b.Open(sealed, name)
		if err != nil && !errors.Is(err, ErrDecrypt) {
			t.Fatalf("unexpected error type: %v", err)
		}
		if err == nil {
			// Only real seals open; re-sealing must round-trip.
			if got, err := b.Open(b.Seal(out, name), name); err != nil || !bytes.Equal(got, out) {
				t.Fatalf("round trip: %v", err)
			}
		}
	})
}
