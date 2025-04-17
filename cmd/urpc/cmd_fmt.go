package main

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/uforg/uforpc/internal/urpc/formatter"
)

type cmdFmtArgs struct {
	Pattern string `arg:"positional" help:"The file pattern to format (support globs e.g. './rpc/**/*.urpc')"`
	Verbose bool   `arg:"-v,--verbose" help:"Verbose output prints all formatted files"`
}

func cmdFmt(args *cmdFmtArgs) {
	var matches []string
	var err error

	startTime := time.Now()
	defer func() {
		log.Printf("formatted %d files in %s", len(matches), time.Since(startTime))
	}()

	matches, err = filepath.Glob(args.Pattern)
	if err != nil {
		log.Fatalf("failed to glob pattern: %s", err)
	}

	for _, match := range matches {
		fileBytes, err := os.ReadFile(match)
		if err != nil {
			log.Fatalf("failed to read file: %s", err)
		}

		formatted, err := formatter.Format(match, string(fileBytes))
		if err != nil {
			log.Fatalf("failed to format file: %s", err)
		}

		if err := os.WriteFile(match, []byte(formatted), 0644); err != nil {
			log.Fatalf("failed to write file: %s", err)
		}

		if args.Verbose {
			log.Println("formatted", match)
		}
	}
}
