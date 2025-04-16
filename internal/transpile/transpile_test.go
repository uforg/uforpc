package transpile

import (
	"embed"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/schema"
	"github.com/uforg/uforpc/internal/urpc/parser"
	"github.com/uforg/uforpc/internal/util/testutil"
)

func TestTranspile(t *testing.T) {
	testFiles := collectTestFiles()

	for _, file := range testFiles {
		assertTranspile(t, file)
	}
}

//////////////////
// TEST HELPERS //
//////////////////

//go:embed tests/*
var urpcTestFiles embed.FS

// testFile represents a URPC schema and its expected JSON representation.
type testFile struct {
	name string
	urpc string
	json string
}

// collectTestFiles collects all test files from the embedded filesystem.
func collectTestFiles() []testFile {
	files, err := urpcTestFiles.ReadDir("tests")
	if err != nil {
		panic(err)
	}

	testFiles := []testFile{}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".urpc") {
			content, err := urpcTestFiles.ReadFile("tests/" + file.Name())
			if err != nil {
				panic(err)
			}

			testFiles = append(testFiles, testFile{
				name: strings.TrimSuffix(file.Name(), ".urpc"),
				urpc: string(content),
			})
		}
	}

	// Populate the JSON content for each test file
	for i, file := range testFiles {
		jsonContent, err := urpcTestFiles.ReadFile("tests/" + file.name + ".json")
		if err != nil {
			panic(err)
		}
		testFiles[i].json = string(jsonContent)
	}

	return testFiles
}

// assertTranspile asserts that the transpilation of a URPC Schema AST <> JSON
// is correct in back and forth.
func assertTranspile(t *testing.T, file testFile) {
	t.Helper()

	astSchema, err := parser.ParserInstance.ParseString("", file.urpc)
	require.NoError(t, err, "error parsing URPC schema")

	jsonSchema, err := schema.ParseSchema(file.json)
	require.NoError(t, err, "error parsing JSON schema")

	// 1. Test AST > JSON > AST
	ast2json, err := ToJSON(*astSchema)
	require.NoError(t, err, "error transpiling AST to JSON")
	require.Equal(t, jsonSchema, ast2json, "incorrect JSON schema")

	ast2json2ast, err := ToURPC(ast2json)
	require.NoError(t, err, "error transpiling JSON to AST")
	testutil.ASTEqualNoPos(t, astSchema, &ast2json2ast)

	// 2. Test JSON > AST > JSON
	json2ast, err := ToURPC(jsonSchema)
	require.NoError(t, err, "error transpiling JSON to AST")
	testutil.ASTEqualNoPos(t, astSchema, &json2ast)

	json2ast2json, err := ToJSON(json2ast)
	require.NoError(t, err, "error transpiling AST to JSON")
	require.Equal(t, jsonSchema, json2ast2json, "incorrect JSON schema")
}
