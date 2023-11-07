package tests

import (
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"os"
	"strings"
	"testing"
)

func TestAddFileToIPFS(t *testing.T) {
	sh := shell.NewShell("localhost:5001")

	cid, err := sh.Add(strings.NewReader("first file!"), shell.CidVersion(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}

	fmt.Printf("added %s", cid)
}
