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
	}

	return g.String(), nil
}

// function createProcedureTypesTemplate() {
//   return `
//     // -----------------------------------------------------------------------------
//     // Procedure Types
//     // -----------------------------------------------------------------------------

//     {{#each procedures}}

//     // P{{name}}Input represents the input parameters for the {{name}} procedure.
//     {{#if input}}
//     type P{{name}}Input struct {
//       {{renderGoFields input}}
//     }
//     {{else}}
//     type P{{name}}Input struct{}
//     {{/if}}

//     // P{{name}}Output represents the output results for the {{name}} procedure.
//     {{#if output}}
//     type P{{name}}Output struct {
//       {{renderGoFields output}}
//     }
//     {{else}}
//     type P{{name}}Output struct{}
//     {{/if}}

//     {{/each}}

//     // ProcedureTypes defines the interface for all procedure types.
//     type ProcedureTypes interface {
//       {{#each procedures}}
//         // {{name}} implements the {{name}} procedure.
//         {{name}}(input P{{name}}Input) (UFOResponse[P{{name}}Output], error)
//       {{/each}}
//     }

//     type UFOProcedureName string

//     var UFOProcedureNames = struct {
//       {{#each procedures}}
//         {{name}} UFOProcedureName
//       {{/each}}
//     }{
//       {{#each procedures}}
//         {{name}}: "{{name}}",
//       {{/each}}
//     }
//   `;
// }
