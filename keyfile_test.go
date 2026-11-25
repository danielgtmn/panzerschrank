package panzerschrank

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

func TestKeyFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), KeyFile)
	b1, err := OpenKeyFile(file)
	if err != nil {
		t.Fatal(err)
	}
	fi, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Size() != KeySize {
		t.Fatalf("size %d", fi.Size())
	}
	if runtime.GOOS != "windows" && fi.Mode().Perm() != 0o600 {
		t.Fatalf("mode %v", fi.Mode().Perm())
	}
	// Opening again reuses the key.
	b2, err := OpenKeyFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := b2.Open(b1.Seal([]byte("x"), "p"), "p"); err != nil || string(got) != "x" {
		t.Fatalf("second open uses another key: %q %v", got, err)
	}
}

func TestKeyFileRefusesWrongSize(t *testing.T) {
	for _, content := range []string{"", "too short", string(make([]byte, KeySize+1))} {
		file := filepath.Join(t.TempDir(), KeyFile)
		if err := os.WriteFile(file, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := OpenKeyFile(file); err == nil {
			t.Fatalf("a key file of %d bytes must be refused", len(content))
		}
		// And it is left alone, not replaced by a new key.
		if got, _ := os.ReadFile(file); string(got) != content {
			t.Fatal("the key file was changed")
		}
	}
}

func TestKeyFileMissingDirectory(t *testing.T) {
	file := filepath.Join(t.TempDir(), "missing", KeyFile)
	if _, err := OpenKeyFile(file); err == nil {
		t.Fatal("a missing directory must be an error")
	}
}

// Several processes (here goroutines) starting at once all get the same key.
func TestKeyFileConcurrentCreate(t *testing.T) {
	file := filepath.Join(t.TempDir(), KeyFile)
	const n = 16
	keys := make([][]byte, n)
	errs := make([]error, n)
	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			keys[i], errs[i] = LoadOrCreateKey(file)
		}()
	}
	wg.Wait()
	for i := range n {
		if errs[i] != nil {
			t.Fatalf("goroutine %d: %v", i, errs[i])
		}
		if string(keys[i]) != string(keys[0]) {
			t.Fatalf("goroutine %d got another key", i)
		}
	}
}
