package panzerschrank

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"time"
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
// with the key of the one that created the file.
//
// Losing this file means losing every secret sealed with it: back it up,
// separately from the database.
func LoadOrCreateKey(file string) ([]byte, error) {
	key, err := os.ReadFile(file)
	if errors.Is(err, fs.ErrNotExist) {
		cerr := createKey(file)
		if cerr != nil && !errors.Is(cerr, fs.ErrExist) {
			return nil, fmt.Errorf("panzerschrank: create key file: %w", cerr)
		}
		// Read back what is on disk. If another process won the race, it
		// may not have written its key yet: wait a moment for it.
		key, err = os.ReadFile(file)
		for i := 0; cerr != nil && err == nil && len(key) < KeySize && i < 40; i++ {
			time.Sleep(25 * time.Millisecond)
			key, err = os.ReadFile(file)
		}
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
// is there already. On any other failure the partial file is removed.
func createKey(file string) (err error) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			os.Remove(file)
		}
	}()
	if _, err := f.Write(GenerateKey()); err != nil {
		return err
	}
	return f.Sync()
}
