package golang

import (
	"fmt"

	"github.com/uforg/uforpc/urpc/internal/genkit"
	"github.com/uforg/uforpc/urpc/internal/schema"
	"github.com/uforg/uforpc/urpc/internal/util/strutil"
)

func generateProcedureTypes(sch schema.Schema, _ Config) (string, error) {
	g := genkit.NewGenKit().WithTabs()

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

		g.Line(renderType("", inputName, inputDesc, procNode.Input))
		g.Break()

		g.Line(renderPreType("", inputName, procNode.Input))
		g.Break()

		g.Line(renderType("", outputName, outputDesc, procNode.Output))
		g.Break()

		g.Linef("// %s", responseDesc)
		g.Linef("type %s = Response[%s]", responseName, outputName)
		g.Break()
	}

	g.Line("// ProcedureTypes defines the interface for all procedure types.")
	g.Line("type ProcedureTypes interface {")
	g.Block(func() {
		for _, procNode := range sch.GetProcNodes() {
			name := procNode.Name

			inputName := fmt.Sprintf("%sInput", strutil.ToPascalCase(name))
			responseName := fmt.Sprintf("%sResponse", strutil.ToPascalCase(name))

			g.Linef("// %s implements the %s procedure.", name, name)
			g.Linef("%s(input %s) %s", name, inputName, responseName)
		}
	})
	g.Line("}")
	g.Break()

	g.Line("// ProcedureName represents the name of a procedure.")
	g.Line("type ProcedureName string")
	g.Break()

	g.Line("// ProcedureNames is a struct that contains all procedure names in its literal string form.")
	g.Line("var ProcedureNames = struct {")
	g.Block(func() {
		for _, procNode := range sch.GetProcNodes() {
			g.Linef("%s ProcedureName", procNode.Name)
		}
	})
	g.Line("}{")
	g.Block(func() {
		for _, procNode := range sch.GetProcNodes() {
			g.Linef("%s: \"%s\",", procNode.Name, procNode.Name)
		}
	})
	g.Line("}")
	g.Break()

	g.Line("// ProcedureNamesList is a list of all procedure names.")
	g.Line("var ProcedureNamesList = []ProcedureName{")
	g.Block(func() {
		for _, procNode := range sch.GetProcNodes() {
			g.Linef("ProcedureNames.%s,", procNode.Name)
		}
	})
	g.Line("}")
	g.Break()

	return g.String(), nil
}
