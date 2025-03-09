package schema_test

import (
	"embed"
	"testing"

	_ "embed"

	"github.com/stretchr/testify/assert"
	"github.com/uforg/uforpc/internal/schema"
)

func TestValidateSchema(t *testing.T) {
	t.Run("valid schemas", func(t *testing.T) {
		testCases, err := getSchemasFromFS(validSchemasFS, "examples/valid")
		assert.NoError(t, err)

		for _, testCase := range testCases {
			t.Run(testCase.fileName, func(t *testing.T) {
				err := schema.ValidateSchema(testCase.schema)
				assert.NoError(t, err)
			})
		}
	})

	t.Run("invalid schemas", func(t *testing.T) {
		testCases, err := getSchemasFromFS(invalidSchemasFS, "examples/invalid")
		assert.NoError(t, err)

		for _, testCase := range testCases {
			t.Run(testCase.fileName, func(t *testing.T) {
				err := schema.ValidateSchema(testCase.schema)
				assert.Error(t, err)
			})
		}
	})
}

/**********
** HELPERS
**********/

//go:embed examples/valid/*.json
var validSchemasFS embed.FS

//go:embed examples/invalid/*.json
var invalidSchemasFS embed.FS

type schemaTestCase struct {
	fileName string
	schema   string
}

func getSchemasFromFS(fs embed.FS, directory string) ([]schemaTestCase, error) {
	files, err := fs.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	schemas := make([]schemaTestCase, 0, len(files))
	for _, file := range files {
		schemaBytes, err := fs.ReadFile(directory + "/" + file.Name())
		if err != nil {
			return nil, err
		}

		schemas = append(schemas, schemaTestCase{
			fileName: file.Name(),
			schema:   string(schemaBytes),
		})
	}

	return schemas, nil
}
