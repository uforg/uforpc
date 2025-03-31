package parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uforg/uforpc/internal/urpc/ast"
)

// setEmptyPos uses reflection to check if the node has Pos or EndPos fields
// and sets them to the empty position (emptyPos).
//
// It handles pointers and non-pointers correctly.
func setEmptyPos[T any](node T) T {
	valueOfNode, valueOfEmptyPos := reflect.ValueOf(node), reflect.ValueOf(ast.Position{})
	fix := func(s reflect.Value) {
		if f := s.FieldByName("Pos"); f.IsValid() && f.CanSet() && f.Type() == valueOfEmptyPos.Type() {
			f.Set(valueOfEmptyPos)
		}
		if f := s.FieldByName("EndPos"); f.IsValid() && f.CanSet() && f.Type() == valueOfEmptyPos.Type() {
			f.Set(valueOfEmptyPos)
		}
	}
	switch valueOfNode.Kind() {
	case reflect.Ptr:
		if !valueOfNode.IsNil() && valueOfNode.Elem().Kind() == reflect.Struct {
			fix(valueOfNode.Elem())
		}
		return node
	case reflect.Struct:
		s := reflect.New(valueOfNode.Type()).Elem()
		s.Set(valueOfNode)
		fix(s)
		return s.Interface().(T)
	default:
		return node
	}
}

// equal compares two URPC structs and fails if they are not equal.
// The validation includes the positions of the AST nodes.
func equal(t *testing.T, expected, actual *ast.URPCSchema) {
	t.Helper()
	require.Equal(t, expected, actual)
}

// equalNoPos compares two URPC structs and fails if they are not equal.
// It ignores the positions of any nested AST nodes.
func equalNoPos(t *testing.T, expected, actual *ast.URPCSchema) {
	t.Helper()

	cleanPositions := func(ast *ast.URPCSchema) *ast.URPCSchema {
		ast = setEmptyPos(ast)

		if ast.Version != nil {
			ast.Version = setEmptyPos(ast.Version)
		}

		for _, rule := range ast.Rules {
			if rule.Docstring != nil {
				rule.Docstring = setEmptyPos(rule.Docstring)
			}
			rule = setEmptyPos(rule)
			rule.Body = setEmptyPos(rule.Body)
		}

		return ast
	}

	expected = cleanPositions(expected)
	actual = cleanPositions(actual)
	equal(t, expected, actual)
}

func TestParserPositions(t *testing.T) {
	t.Run("Version position", func(t *testing.T) {
		input := `version: 1`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)
		require.NotNil(t, parsed)

		expected := &ast.URPCSchema{
			Pos: ast.Position{
				Filename: "schema.urpc",
				Line:     1,
				Offset:   0,
				Column:   1,
			},
			EndPos: ast.Position{
				Filename: "schema.urpc",
				Line:     1,
				Offset:   10,
				Column:   11,
			},
			Version: &ast.Version{
				Pos: ast.Position{
					Filename: "schema.urpc",
					Line:     1,
					Offset:   0,
					Column:   1,
				},
				EndPos: ast.Position{
					Filename: "schema.urpc",
					Line:     1,
					Offset:   10,
					Column:   11,
				},
				Number: 1,
			},
		}

		equal(t, expected, parsed)
	})
}

func TestParserVersion(t *testing.T) {
	t.Run("Correct version parsing", func(t *testing.T) {
		input := `version: 1`
		parsed, err := Parser.ParseString("schema.urpc", input)

		require.NoError(t, err)
		require.NotNil(t, parsed)

		expected := &ast.URPCSchema{
			Version: &ast.Version{
				Number: 1,
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("More than one version should fail", func(t *testing.T) {
		input := `version: 1 version: 2`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Version as float should fail", func(t *testing.T) {
		input := `version: 1.0`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Version as identifier should fail", func(t *testing.T) {
		input := `version: version`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Version as string should fail", func(t *testing.T) {
		input := `version: "1"`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})
}

func TestParserRuleDecl(t *testing.T) {
	t.Run("Minimum rule declaration parsing", func(t *testing.T) {
		input := `
			rule @myRule {
				for: string
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Rules: []*ast.RuleDecl{
				{
					Name: "myRule",
					Body: ast.RuleDeclBody{
						For: "string",
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Empty rule not allowed", func(t *testing.T) {
		input := `rule @myRule {}`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("For in rule body is required", func(t *testing.T) {
		input := `rule @myRule { param: string }`
		_, err := Parser.ParseString("schema.urpc", input)
		require.Error(t, err)
	})

	t.Run("Rule with docstring", func(t *testing.T) {
		input := `
			"""
			My rule description
			"""
			rule @myRule {
				for: string
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Rules: []*ast.RuleDecl{
				{
					Docstring: &ast.Docstring{Content: "My rule description"},
					Name:      "myRule",
					Body: ast.RuleDeclBody{
						For: "string",
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})

	t.Run("Rule with all options", func(t *testing.T) {
		input := `
			"""
			My rule description
			"""
			rule @myRule {
				for: string
				param: string
				error: "My error message"
			}
		`
		parsed, err := Parser.ParseString("schema.urpc", input)
		require.NoError(t, err)

		expected := &ast.URPCSchema{
			Rules: []*ast.RuleDecl{
				{
					Docstring: &ast.Docstring{Content: "My rule description"},
					Name:      "myRule",
					Body: ast.RuleDeclBody{
						For:   "string",
						Param: "string",
						Error: "My error message",
					},
				},
			},
		}

		equalNoPos(t, expected, parsed)
	})
}
