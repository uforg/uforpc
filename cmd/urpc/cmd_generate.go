package main

import (
	"log"
	"time"

	"github.com/uforg/uforpc/internal/codegen"
)

type cmdGenerateArgs struct {
	ConfigPath string `arg:"positional" help:"The urpc.toml config file path (default: ./urpc.toml)"`
}

func cmdGenerate(args *cmdGenerateArgs) {
	startTime := time.Now()

	if args.ConfigPath == "" {
		args.ConfigPath = "./urpc.toml"
	}

	if err := codegen.Run(args.ConfigPath); err != nil {
		log.Fatalf("failed to run code generator: %s", err)
	}

	log.Printf("code generation finished in %s", time.Since(startTime))
}
