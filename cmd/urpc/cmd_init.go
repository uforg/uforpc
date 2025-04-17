package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type cmdInitArgs struct {
	Path string `arg:"positional" help:"The file or directory path to initialize the URPC schema in, if no file name is provided, the file will be named 'schema.urpc'"`
}

func cmdInit(args *cmdInitArgs) {
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
