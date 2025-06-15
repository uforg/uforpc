package golang

import (
	_ "embed"
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

//go:embed pieces/client.go
var clientRawPiece string

func generateClient(sch schema.Schema, config Config) (string, error) {
	if !config.IncludeClient {
		return "", nil
	}

	piece := strutil.GetStrAfter(clientRawPiece, "/** START FROM HERE **/")
	if piece == "" {
		return "", fmt.Errorf("client.go: could not find start delimiter")
	}

	g := genkit.NewGenKit().WithTabs()

	g.Raw(piece)
	g.Break()

	g.Line("// -----------------------------------------------------------------------------")
	g.Line("// Client generated implementation")
	g.Line("// -----------------------------------------------------------------------------")
	g.Break()

	// TODO: Generate the client fluent builder wrapper

	// TODO: Generate the client proc methods
	// for _, procNode := range sch.GetProcNodes() {
	// 	name := strutil.ToPascalCase(procNode.Name)
	// }

	// TODO: Generate the client stream methods
	// for _, streamNode := range sch.GetStreamNodes() {
	// 	name := strutil.ToPascalCase(streamNode.Name)
	// }

	return g.String(), nil
}
