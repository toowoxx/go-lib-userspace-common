package tar

import (
	"io"
	"log"
	"testing"
)

func TestCreate(t *testing.T) {
	reader, errChan, err := Create(CreateOptions{
		Archive:            true,
		Directory:          "../",
		ExcludePatterns:    []string{"go.mod", "*.sum"},
		ExcludeVCS:         true,
		ExcludeVCSIgnores:  true,
		NumericOwner:       true,
		ACLs:               true,
		SELinux:            true,
		ExtendedAttributes: true,
		FullTime:           true,
		Verbose:            true,
	})
	if err != nil {
		t.Fatal(err)
	}

	written, err := io.Copy(io.Discard, reader)
	if err != nil {
		t.Fatal(err)
	}
	if err, hasErr := <-errChan; hasErr {
		t.Fatal(err)
	}

	log.Println("Bytes written:", written)
}
