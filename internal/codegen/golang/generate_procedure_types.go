package golang

import (
	"fmt"

	"github.com/uforg/uforpc/internal/codegen/genkit"
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/util/strutil"
)

func generateProcedureTypes(sch schema.Schema, _ Config) (string, error) {
	g := genkit.NewGenKit().WithTabs()

	g.Line("// -----------------------------------------------------------------------------")
	g.Line("// Procedure Types")
	g.Line("// -----------------------------------------------------------------------------")
	g.Break()

	for name, procedure := range sch.Procedures {
		namePascal := strutil.ToPascalCase(name)
		inputName := fmt.Sprintf("P%sInput", namePascal)
		outputName := fmt.Sprintf("P%sOutput", namePascal)
		responseName := fmt.Sprintf("P%sResponse", namePascal)

		inputType := generateCommonRenderField(generateCommonRenderFieldParams{
			name:     inputName,
			field:    procedure.Input,
			typeOnly: true,
			omitTag:  true,
		})
		if inputType == "" {
			inputType = "struct{}"
		}

		outputType := generateCommonRenderField(generateCommonRenderFieldParams{
			name:     outputName,
			field:    procedure.Output,
			typeOnly: true,
			omitTag:  true,
		})
		if outputType == "" {
			outputType = "struct{}"
		}

		g.Linef("// %s represents the input parameters for the %s procedure.", inputName, namePascal)
		g.Linef("type %s %s", inputName, inputType)
		g.Break()

		g.Linef("// %s represents the output results for the %s procedure.", outputName, namePascal)
		g.Linef("type %s %s", outputName, outputType)
		g.Break()

		g.Linef("// %s represents the response for the %s procedure.", responseName, namePascal)
		g.Linef("type %s = Response[%s]", responseName, outputName)
		g.Break()
	}

	g.Line("// ProcedureTypes defines the interface for all procedure types.")
	g.Line("type ProcedureTypes interface {")
	g.Block(func() {
		for name := range sch.Procedures {
			inputName := fmt.Sprintf("P%sInput", strutil.ToPascalCase(name))
			outputName := fmt.Sprintf("P%sResponse", strutil.ToPascalCase(name))

			g.Linef("// %s implements the %s procedure.", name, name)
			g.Linef("%s(input %s) (%s, error)", name, inputName, outputName)
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
		for name := range sch.Procedures {
			g.Linef("%s ProcedureName", name)
		}
	})
	g.Line("}{")
	g.Block(func() {
		for name := range sch.Procedures {
			g.Linef("%s: \"%s\",", name, name)
		}
	})
	g.Line("}")
	g.Break()

	return g.String(), nil
}
