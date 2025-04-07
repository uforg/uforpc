package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/alexflint/go-arg"
	"github.com/uforg/uforpc/internal/urpc/lsp"
)

//go:embed init-schema.urpc
var initSchema []byte

type allArgs struct {
	Init *initCmd `arg:"subcommand:init" help:"Initialize a new URPC schema in the specified path"`
	LSP  *lspCmd  `arg:"subcommand:lsp" help:"Start the LSP server"`
}

type initCmd struct {
	Path string `arg:"positional" help:"The file or directory path to initialize the URPC schema in, if no file name is provided, the file will be named 'schema.urpc'"`
}

type lspCmd struct{}

func main() {
	// If the LSP is called, then omit the arg parser to avoid taking
	// control of the stdin/stdout because the LSP will need it.
	if len(os.Args) > 1 && os.Args[1] == "lsp" {
		lspCmdFn(nil)
		return
	}

	var args allArgs
	arg.MustParse(&args)

	if args.Init != nil {
		initCmdFn(args.Init)
	}
}

func initCmdFn(args *initCmd) {
	if args.Path == "" || args.Path == "." {
		args.Path = "./schema.urpc"
	}

	if !strings.HasSuffix(args.Path, ".urpc") {
		args.Path = path.Join(args.Path, "schema.urpc")
	}

	if err := os.WriteFile(args.Path, initSchema, 0644); err != nil {
		log.Fatalf("failed to write init schema: %s", err)
	}

	fmt.Printf("URPC schema initialized in %s\n", args.Path)
}

func lspCmdFn(_ *lspCmd) {
	lspInstance := lsp.New(os.Stdin, os.Stdout)
	if err := lspInstance.Run(); err != nil {
		log.Fatalf("failed to run lsp: %s", err)
	}
}
