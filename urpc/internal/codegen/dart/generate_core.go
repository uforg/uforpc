package dart

import (
	_ "embed"

	"github.com/uforg/uforpc/urpc/internal/schema"
)

//go:embed pieces/core.dart
var coreRawPiece string

func generateCore(_ schema.Schema, _ Config) (string, error) {
	return coreRawPiece, nil
}
