package golang

import (
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

func generateStreamTypes(sch schema.Schema, _ Config) (string, error) {
	g := genkit.NewGenKit().WithTabs()

	g.Line("// -----------------------------------------------------------------------------")
	g.Line("// Stream Types")
	g.Line("// -----------------------------------------------------------------------------")
	g.Break()

	for _, streamNode := range sch.GetStreamNodes() {
		namePascal := strutil.ToPascalCase(streamNode.Name)
		inputName := fmt.Sprintf("%sInput", namePascal)
		outputName := fmt.Sprintf("%sOutput", namePascal)
		responseName := fmt.Sprintf("%sResponse", namePascal)

		inputDesc := fmt.Sprintf("%s represents the input parameters for the %s stream.", inputName, namePascal)
		outputDesc := fmt.Sprintf("%s represents the output parameters for the %s stream.", outputName, namePascal)
		responseDesc := fmt.Sprintf("%s represents the response for the %s stream.", responseName, namePascal)

		g.Line(renderType("", inputName, inputDesc, streamNode.Input))
		g.Break()

		g.Line(renderPreType("", inputName, streamNode.Input))
		g.Break()

		g.Line(renderType("", outputName, outputDesc, streamNode.Output))
		g.Break()

		g.Linef("// %s", responseDesc)
		g.Linef("type %s = Response[%s]", responseName, outputName)
		g.Break()
	}

	g.Line("// StreamTypes defines the interface for all stream types.")
	g.Line("type StreamTypes interface {")
	g.Block(func() {
		for _, streamNode := range sch.GetStreamNodes() {
			name := streamNode.Name

			inputName := fmt.Sprintf("%sInput", strutil.ToPascalCase(name))
			responseName := fmt.Sprintf("%sResponse", strutil.ToPascalCase(name))

			g.Linef("// %s implements the %s stream.", name, name)
			g.Linef("%s(input %s) %s", name, inputName, responseName)
		}
	})
	g.Line("}")
	g.Break()

	g.Line("// ufoStreamNames is a list of all stream names.")
	g.Line("var ufoStreamNames = []string{")
	g.Block(func() {
		for _, streamNode := range sch.GetStreamNodes() {
			g.Linef("\"%s\",", streamNode.Name)
		}
	})
	g.Line("}")
	g.Break()

	return g.String(), nil
}
