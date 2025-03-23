package golang

import (
	_ "embed"
	"fmt"

	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/util/strutil"
)

//go:embed pieces/nullutils.go
var nullUtilsRawPiece string

func generateNullUtils(_ schema.Schema, _ Config) (string, error) {
	piece := strutil.GetStrAfter(nullUtilsRawPiece, "/** START FROM HERE **/")
	if piece == "" {
		return "", fmt.Errorf("nullutils.go: could not find start delimiter")
	}
	return piece, nil
}
