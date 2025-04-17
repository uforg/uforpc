package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alexflint/go-arg"
)

type allArgs struct {
	Init      *cmdInitArgs      `arg:"subcommand:init" help:"Initialize a new URPC schema in the specified path"`
	Fmt       *cmdFmtArgs       `arg:"subcommand:fmt" help:"Format the URPC schema in the specified path"`
	Transpile *cmdTranspileArgs `arg:"subcommand:transpile" help:"Transpile a URPC schema to JSON and vice versa, the result will be printed to stdout"`
	LSP       *cmdLSPArgs       `arg:"subcommand:lsp" help:"Start the LSP server"`
}

func main() {
	// If the LSP is called, then omit the arg parser to avoid taking
	// control of the stdin/stdout because the LSP will need it.
	if len(os.Args) > 1 && os.Args[1] == "lsp" {
		cmdLSP(nil)
		return
	}

	var args allArgs
	p, err := arg.NewParser(arg.Config{}, &args)
	if err != nil {
		log.Fatalf("failed to create arg parser: %s", err)
	}

	err = p.Parse(os.Args[1:])
	switch {
	case err == arg.ErrHelp: // indicates that user wrote "--help" on command line
		p.WriteHelp(os.Stdout)
		os.Exit(0)
	case err != nil:
		fmt.Printf("error: %v\n", err)
		p.WriteUsage(os.Stdout)
		os.Exit(1)
	}

	if args.Init != nil {
		cmdInit(args.Init)
	}

	if args.Fmt != nil {
		cmdFmt(args.Fmt)
	}

	if args.Transpile != nil {
		cmdTranspile(args.Transpile)
	}
}
