package main

import (
	_ "embed"
	"os"

	"github.com/alexflint/go-arg"
)

//go:embed init-schema.urpc
var initSchema []byte

type allArgs struct {
	Init *cmdInitArgs `arg:"subcommand:init" help:"Initialize a new URPC schema in the specified path"`
	Fmt  *cmdFmtArgs  `arg:"subcommand:fmt" help:"Format the URPC schema in the specified path"`
	LSP  *cmdLSPArgs  `arg:"subcommand:lsp" help:"Start the LSP server"`
}

func main() {
	// If the LSP is called, then omit the arg parser to avoid taking
	// control of the stdin/stdout because the LSP will need it.
	if len(os.Args) > 1 && os.Args[1] == "lsp" {
		cmdLSP(nil)
		return
	}

	var args allArgs
	arg.MustParse(&args)

	if args.Init != nil {
		cmdInit(args.Init)
	}

	if args.Fmt != nil {
		cmdFmt(args.Fmt)
	}
}
