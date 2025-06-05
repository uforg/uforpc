package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

//go:embed cmd_init_schema.urpc
var initSchema []byte

//go:embed cmd_init_config.toml
var initConfig []byte

type cmdInitArgs struct {
	Path string `arg:"positional" help:"The directory path where URPC schema and config files will be created. Defaults to the current directory."`
}

func cmdInit(args *cmdInitArgs) {
	if args.Path == "" {
		args.Path = "."
	}

	// Validate that path is a directory
	if info, err := os.Stat(args.Path); err == nil && !info.IsDir() {
		log.Fatalf("path must be a directory, not a file: %s", args.Path)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(args.Path, 0755); err != nil {
		log.Fatalf("failed to create directory: %s", err)
	}

	// Generate unique filenames
	schemaPath, configPath := generateUniqueFilenames(args.Path)

	// Write both files
	if err := os.WriteFile(schemaPath, initSchema, 0644); err != nil {
		log.Fatalf("failed to write schema file: %s", err)
	}

	if err := os.WriteFile(configPath, initConfig, 0644); err != nil {
		log.Fatalf("failed to write config file: %s", err)
	}

	fmt.Printf("URPC files initialized:\n- %s\n- %s\n", schemaPath, configPath)
}

func generateUniqueFilenames(dir string) (string, string) {
	schemaName := "schema.urpc"
	configName := "uforpc.toml"

	schemaPath := filepath.Join(dir, schemaName)
	configPath := filepath.Join(dir, configName)

	// Check if files exist
	if _, err := os.Stat(schemaPath); os.IsNotExist(err) {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return schemaPath, configPath
		}
	}

	// Generate unique suffix using unix timestamp
	timestamp := time.Now().Unix()
	schemaName = fmt.Sprintf("schema-%d.urpc", timestamp)
	configName = fmt.Sprintf("uforpc-%d.toml", timestamp)

	return filepath.Join(dir, schemaName), filepath.Join(dir, configName)
}
