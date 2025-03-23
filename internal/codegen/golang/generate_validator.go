package golang

import (
	_ "embed"
	"fmt"

	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/util/strutil"
)

//go:embed pieces/validator.go
var validatorRawPiece string

func generateValidator(_ schema.Schema, config Config) (string, error) {
	if config.OmitClientRequestValidation && config.OmitServerRequestValidation {
		return "", nil
	}

	piece := strutil.GetStrAfter(validatorRawPiece, "/** START FROM HERE **/")
	if piece == "" {
		return "", fmt.Errorf("validator.go: could not find start delimiter")
	}
	return piece, nil
}
