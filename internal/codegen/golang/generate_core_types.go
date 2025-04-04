package golang

import (
	_ "embed"
	"fmt"

	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/util/strutil"
)

//go:embed pieces/coretypes.go
var coreTypesRawPiece string

func generateCoreTypes(_ schema.Schema, _ Config) (string, error) {
	piece := strutil.GetStrAfter(coreTypesRawPiece, "/** START FROM HERE **/")
	if piece == "" {
		return "", fmt.Errorf("coretypes.go: could not find start delimiter")
	}
	return piece, nil
}
