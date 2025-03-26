package parser

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/urpc/lexer"
)

func TestParser(t *testing.T) {
	t.Run("Parse version", func(t *testing.T) {
		lexer := lexer.NewLexer("test.urpc", "version 2")
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Version: ast.Version{
				IsSet: true,
				Value: 2,
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
	})

	t.Run("Parse version invalid", func(t *testing.T) {
		lexer := lexer.NewLexer("test.urpc", "version foobar")
		parser := New(lexer)
		_, _, err := parser.Parse()
		require.Error(t, err)
		require.Contains(t, err.Error(), "version expected to be an integer")
	})

	t.Run("Parse version already set", func(t *testing.T) {
		lexer := lexer.NewLexer("test.urpc", "version 2 version 3")
		parser := New(lexer)
		_, _, err := parser.Parse()
		require.Error(t, err)
		require.Contains(t, err.Error(), "version already set")
	})

	t.Run("Parse type declaration basic", func(t *testing.T) {
		input := `
			type User {}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
	})

	t.Run("Parse type declaration with docstring", func(t *testing.T) {
		input := `
			""" Product type documentation """
			type Product {}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDeclaration{
				{
					Doc:  "Product type documentation",
					Name: "Product",
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
	})

	t.Run("Parse type declaration with primitive type fields", func(t *testing.T) {
		input := `
			type User {
				name: string
				age?: int
				height: float
				isActive: boolean
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "name",
							Optional: false,
							Type:     &ast.TypeString{},
						},
						{
							Name:     "age",
							Optional: true,
							Type:     &ast.TypeInt{},
						},
						{
							Name:     "height",
							Optional: false,
							Type:     &ast.TypeFloat{},
						},
						{
							Name:     "isActive",
							Optional: false,
							Type:     &ast.TypeBoolean{},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
	})

	t.Run("Parse type declaration with field using other custom type", func(t *testing.T) {
		input := `
			type User {
				name: string
				address?: Address
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Types: []ast.TypeDeclaration{
				{
					Name: "User",
					Fields: []ast.Field{
						{
							Name:     "name",
							Optional: false,
							Type:     &ast.TypeString{},
						},
						{
							Name:     "address",
							Optional: true,
							Type:     &ast.TypeCustom{Name: "Address"},
						},
					},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
	})

	t.Run("Parse procedure declaration basic", func(t *testing.T) {
		input := `
			proc CreateUser {}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Procedures: []ast.ProcDeclaration{
				{
					Name:     "CreateUser",
					Input:    ast.Input{},
					Output:   ast.Output{},
					Metadata: ast.Metadata{},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
	})

	t.Run("Parse procedure declaration with docstring", func(t *testing.T) {
		input := `
			""" Create user procedure documentation """
			proc CreateUser {}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Procedures: []ast.ProcDeclaration{
				{
					Doc:      "Create user procedure documentation",
					Name:     "CreateUser",
					Input:    ast.Input{},
					Output:   ast.Output{},
					Metadata: ast.Metadata{},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
	})

	t.Run("Parse procedure declaration with empty input, output and meta", func(t *testing.T) {
		input := `
			proc CreateUser {
				input {}
				output {}
				meta {}
			}
		`

		lexer := lexer.NewLexer("test.urpc", input)
		parser := New(lexer)
		schema, _, err := parser.Parse()

		expected := ast.Schema{
			Procedures: []ast.ProcDeclaration{
				{
					Name:     "CreateUser",
					Input:    ast.Input{},
					Output:   ast.Output{},
					Metadata: ast.Metadata{},
				},
			},
		}

		require.NoError(t, err)
		require.Equal(t, expected, schema)
	})
}
