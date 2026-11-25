package panzerschrank

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// KeyFile is the suggested name of the key file inside an app's data
// directory.
const KeyFile = "secret.key"

// OpenKeyFile returns a Box for the key in file, creating the file with a
// new random key (mode 0600) if it does not exist. See LoadOrCreateKey.
func OpenKeyFile(file string) (*Box, error) {
	key, err := LoadOrCreateKey(file)
	if err != nil {
		return nil, err
	}
	return New(key)
}

// LoadOrCreateKey reads the raw 32-byte key in file. If the file does not
// exist it is created with a new random key and mode 0600; the directory
// must exist. When several processes start at once, all of them end up
// with the key of the one that created the file. The file never appears
// partly written, also not after a crash during creation.
//
// Losing this file means losing every secret sealed with it: back it up,
// separately from the database.
func LoadOrCreateKey(file string) ([]byte, error) {
	key, err := os.ReadFile(file)
	if errors.Is(err, fs.ErrNotExist) {
		if err := createKey(file); err != nil && !errors.Is(err, fs.ErrExist) {
			return nil, fmt.Errorf("panzerschrank: create key file: %w", err)
		}
		// Read back what is on disk, also when another process won the race.
		key, err = os.ReadFile(file)
	}
	if err != nil {
		return nil, fmt.Errorf("panzerschrank: read key file: %w", err)
	}
	if len(key) != KeySize {
		return nil, fmt.Errorf("panzerschrank: %s must hold exactly %d bytes, it has %d", file, KeySize, len(key))
	}
	return key, nil
}

// createKey writes a new key to file, failing with fs.ErrExist if the file
// is there already. The key is written and synced to a temporary file in
// the same directory first, then hard-linked to file: the link appears
// complete or not at all, and fails if file exists. Creating file
// directly with O_EXCL would let a concurrent reader see it empty.
func createKey(file string) error {
	dir := filepath.Dir(file)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(file)+".tmp-*") // mode 0600
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	_, werr := tmp.Write(GenerateKey())
	if werr == nil {
		werr = tmp.Sync()
	}
	if cerr := tmp.Close(); werr == nil {
		werr = cerr
	}
	if werr != nil {
		return werr
	}
	if err := os.Link(tmp.Name(), file); err != nil {
		return err
	}
	// Make the new directory entry durable, so a power loss cannot drop the
	// key after secrets were sealed with it. Not possible on every system
	// (Windows); the file itself is synced above.
	if d, err := os.Open(dir); err == nil {
		d.Sync()
		d.Close()
	}
	return nil
}
