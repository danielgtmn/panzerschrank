package panzerschrank_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/danielgtmn/panzerschrank"
)

// The key lives in a file next to the database and is created on first
// start.
func Example() {
	dataDir, _ := os.MkdirTemp("", "app")
	defer os.RemoveAll(dataDir)

	box, err := panzerschrank.OpenKeyFile(filepath.Join(dataDir, panzerschrank.KeyFile))
	if err != nil {
		log.Fatal(err)
	}

	// Store sealed in the database, bound to what it is.
	sealed := box.SealString("hetzner-token-123", "dns-token")

	token, err := box.OpenString(sealed, "dns-token")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(token)

	// The same value under another name does not open.
	_, err = box.OpenString(sealed, "smtp-password")
	fmt.Println(err != nil)
	// Output:
	// hetzner-token-123
	// true
}

// The key comes from an environment variable in base64, e.g. created with
// `openssl rand -base64 32`.
func ExampleNewFromBase64() {
	key := panzerschrank.GenerateKeyBase64() // os.Getenv("APP_SECRET_KEY")
	box, err := panzerschrank.NewFromBase64(key)
	if err != nil {
		log.Fatal(err)
	}
	sealed := box.Seal([]byte("s3cret"), "DB_PASSWORD")
	plain, _ := box.Open(sealed, "DB_PASSWORD")
	fmt.Println(string(plain))
	// Output: s3cret
}
