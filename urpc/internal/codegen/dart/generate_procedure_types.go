package dart

import (
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

func generateProcedureTypes(sch schema.Schema, _ Config) (string, error) {
	g := genkit.NewGenKit().WithSpaces(2)

	g.Line("// -----------------------------------------------------------------------------")
	g.Line("// Procedure Types")
	g.Line("// -----------------------------------------------------------------------------")
	g.Break()

	for _, procNode := range sch.GetProcNodes() {
		namePascal := strutil.ToPascalCase(procNode.Name)
		inputName := fmt.Sprintf("%sInput", namePascal)
		outputName := fmt.Sprintf("%sOutput", namePascal)
		responseName := fmt.Sprintf("%sResponse", namePascal)

		inputDesc := fmt.Sprintf("%s represents the input parameters for the %s procedure.", inputName, namePascal)
		outputDesc := fmt.Sprintf("%s represents the output parameters for the %s procedure.", outputName, namePascal)
		responseDesc := fmt.Sprintf("%s represents the response for the %s procedure.", responseName, namePascal)

		g.Line(renderDartType("", inputName, inputDesc, procNode.Input))
		g.Break()

		g.Line(renderDartType("", outputName, outputDesc, procNode.Output))
		g.Break()

		g.Linef("// %s", responseDesc)
		g.Linef("typedef %s = Response<%s>;", responseName, outputName)
		g.Break()
	}

	g.Line("// __ufoProcedureNames is a list of all procedure names.")
	g.Line("const List<String> __ufoProcedureNames = [")
	g.Block(func() {
		for _, procNode := range sch.GetProcNodes() {
			g.Linef("'%s',", procNode.Name)
		}
	})
	g.Line("];")
	g.Break()

	return g.String(), nil
}
