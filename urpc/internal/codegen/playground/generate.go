package playground

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/uforg/uforpc/embedplayground"
	"github.com/uforg/uforpc/urpc/internal/urpc/ast"
	"github.com/uforg/uforpc/urpc/internal/urpc/formatter"
)

// Generate takes a schema and a config and generates the playground for the schema.
func Generate(absConfigDir string, sch *ast.Schema, config Config) error {
	outputDir := filepath.Join(absConfigDir, config.OutputDir)

	err := extractEmbedFS(embedplayground.BuildFS, "build", outputDir)
	if err != nil {
		return fmt.Errorf("error extracting embedded filesystem: %w", err)
	}

	formattedSchema := formatter.FormatSchema(sch)
	formattedSchemaPath := filepath.Join(outputDir, "schema.urpc")
	if err := os.WriteFile(formattedSchemaPath, []byte(formattedSchema), 0644); err != nil {
		return fmt.Errorf("error writing formatted schema to %s: %w", formattedSchemaPath, err)
	}

	return nil
}

func extractEmbedFS(embedFS embed.FS, rootDir string, destDir string) error {
	return fs.WalkDir(embedFS, rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0o700)
		}

		data, err := fs.ReadFile(embedFS, path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, 0o644)
	})
}
