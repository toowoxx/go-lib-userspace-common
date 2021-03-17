package xdelta

import (
	"fmt"
	"log"
	"os"
	"path"
	"testing"
)

const testDataDir = "./test_data"

func inTestDir(file string) string {
	return path.Join(testDataDir, file)
}

func TestAll(t *testing.T) {
	defer func() {
		_ = os.RemoveAll(testDataDir)
	}()

	original := "abcdefghijklmnopqrstuvwxyz0123456789876543210zyxwvutsrqponmlkjihgfedcba"
	modified := "acdefghijklmnop123tuv lo 0123456789876543210zyxfedghicba"

	var err error
	err = os.MkdirAll(testDataDir, 0750)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(inTestDir("test.original"), []byte(original), 0640)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile(inTestDir("test.target"), []byte(modified), 0640)
	if err != nil {
		t.Fatal(err)
	}

	foundDiff, err := Delta(inTestDir("test.original"), inTestDir("test.target"), inTestDir("test.delta"))
	if err != nil {
		t.Fatal(err)
	}
	if !foundDiff {
		t.Fatal("expected xdelta to find a difference")
	}

	output, err := Info(inTestDir("test.delta"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(output)

	err = Patch(inTestDir("test.original"), inTestDir("test.delta"), inTestDir("test.patched"))
	if err != nil {
		t.Fatal(err)
	}

	byt, _ := os.ReadFile(inTestDir("test.patched"))
	content := string(byt)

	if content != modified {
		log.Println(content, "!=", modified)
		t.Fatal("patched file mismatch")
	}
}
