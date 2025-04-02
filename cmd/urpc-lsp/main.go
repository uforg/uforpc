package main

import (
	"log"
	"os"

	"github.com/uforg/uforpc/internal/urpc/lsp"
)

func main() {
	lspInstance := lsp.New(os.Stdin, os.Stdout)
	if err := lspInstance.Run(); err != nil {
		log.Fatalf("failed to run lsp: %s", err)
	}
}
