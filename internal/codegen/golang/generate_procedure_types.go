package golang

import (
	"fmt"

	"github.com/uforg/uforpc/internal/genkit"
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/util/strutil"
)

func generateProcedureTypes(sch schema.Schema, _ Config) (string, error) {
	g := genkit.NewGenKit().WithTabs()

	g.Line("// -----------------------------------------------------------------------------")
	g.Line("// Procedure Types")
	g.Line("// -----------------------------------------------------------------------------")
	g.Break()

	for _, procNode := range sch.GetProcNodes() {
		namePascal := strutil.ToPascalCase(procNode.Name)
		inputName := fmt.Sprintf("P%sInput", namePascal)
		outputName := fmt.Sprintf("P%sOutput", namePascal)
		responseName := fmt.Sprintf("P%sResponse", namePascal)

		inputType := generateCommonRenderStructFromFieldSlice(procNode.Input)
		if inputType == "" {
			inputType = "struct{}"
		}

		outputType := generateCommonRenderStructFromFieldSlice(procNode.Output)
		if outputType == "" {
			outputType = "struct{}"
		}

		g.Linef("// %s represents the input parameters for the %s procedure.", inputName, namePascal)
		g.Linef("type %s = %s", inputName, inputType)
		g.Break()

		g.Linef("// %s represents the output results for the %s procedure.", outputName, namePascal)
		g.Linef("type %s = %s", outputName, outputType)
		g.Break()

		g.Linef("// %s represents the response for the %s procedure.", responseName, namePascal)
		g.Linef("type %s = Response[%s]", responseName, outputName)
		g.Break()
	}

	g.Line("// ProcedureTypes defines the interface for all procedure types.")
	g.Line("type ProcedureTypes interface {")
	g.Block(func() {
		for _, procNode := range sch.GetProcNodes() {
			name := procNode.Name

			inputName := fmt.Sprintf("P%sInput", strutil.ToPascalCase(name))
			responseName := fmt.Sprintf("P%sResponse", strutil.ToPascalCase(name))

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
