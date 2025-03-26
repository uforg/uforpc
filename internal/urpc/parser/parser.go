package parser

import (
	"fmt"
	"strconv"

	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/urpc/lexer"
	"github.com/uforg/uforpc/internal/urpc/token"
)

type Parser struct {
	lex               *lexer.Lexer
	tokens            []token.Token
	errors            []error
	maxIndex          int
	currentIndex      int
	currentIndexIsEOF bool
	currentToken      token.Token
}

// New creates and initializes a new Parser from a Lexer.
func New(lex *lexer.Lexer) *Parser {
	p := &Parser{}

	p.lex = lex
	p.tokens = lex.ReadTokens()
	p.errors = []error{}
	p.maxIndex = len(p.tokens) - 1

	p.currentIndex = 0
	if p.maxIndex <= 0 {
		p.currentIndexIsEOF = true
	} else {
		p.currentIndexIsEOF = false
	}
	if p.currentIndexIsEOF {
		p.currentToken = token.Token{
			Type:     token.EOF,
			Literal:  "",
			FileName: p.lex.FileName,
			Line:     p.lex.CurrentLine,
			Column:   p.lex.CurrentColumn,
		}
	} else {
		p.currentToken = p.tokens[p.currentIndex]
	}

	return p
}

// readNextToken reads the next token from the tokens list and updates the
// current Parser state.
func (p *Parser) readNextToken() {
	if p.currentIndexIsEOF {
		return
	}

	p.currentIndex++
	if p.currentIndex > p.maxIndex {
		p.currentIndexIsEOF = true
		p.currentToken = token.Token{
			Type:     token.EOF,
			Literal:  "",
			FileName: p.lex.FileName,
			Line:     p.lex.CurrentLine,
			Column:   p.lex.CurrentColumn,
		}
	} else {
		p.currentToken = p.tokens[p.currentIndex]
	}
}

// peekToken peeks the token at the next index + depth without moving the current index.
//
// Returns the token and a boolean indicating if the EOF was reached.
func (p *Parser) peekToken(depth int) (token.Token, bool) {
	indexToPeek := p.currentIndex + depth
	if indexToPeek > p.maxIndex {
		return token.Token{
			Type:     token.EOF,
			Literal:  "",
			FileName: p.lex.FileName,
			Line:     p.lex.CurrentLine,
			Column:   p.lex.CurrentColumn,
		}, true
	}
	return p.tokens[indexToPeek], false
}

// appendError appends an error to the parser's errors list.
//
// The error message is formatted with the current token's file name, line, and column.
func (p *Parser) appendError(message string) {
	fileName := p.currentToken.FileName
	line := p.currentToken.Line
	column := p.currentToken.Column
	err := fmt.Errorf("%s Ln %d Col %d: %s", fileName, line, column, message)
	p.errors = append(p.errors, err)
}

// expectToken expects the current token to be of the given type.
//
// If the current token is not of the given type, an error is added to the parser's errors list.
//
// If a message is provided, it is appended to the error message.
//
// Returns:
//   - bool: indicating if the token was encountered.
func (p *Parser) expectToken(expectedType token.TokenType, message ...string) bool {
	if p.currentToken.Type != expectedType {
		msg := fmt.Sprintf("expected token \"%s\", got \"%s\"", expectedType, p.currentToken.Type)
		if len(message) > 0 {
			msg += fmt.Sprintf(": %s", message[0])
		}
		p.appendError(msg)
		return false
	}
	return true
}

// skipNewlines skips all newline tokens from the current index to the next non-newline token.
func (p *Parser) skipNewlines() {
	for p.currentToken.Type == token.NEWLINE {
		p.readNextToken()
	}
}

// Parse parses the tokens provided by the Lexer.
//
// Returns:
//   - ast.Schema: The parsed AST.
//   - []error: A list of all errors encountered during parsing.
//   - error: The first error encountered during parsing, or nil if no errors were encountered.
func (p *Parser) Parse() (ast.Schema, []error, error) {
	schema := ast.Schema{}

	for p.currentToken.Type != token.EOF {
		switch p.currentToken.Type {
		case token.VERSION:
			schema.Version = p.parseVersion(schema)
		case token.DOCSTRING:
			typeDecl, procDecl := p.parseDocstring()
			if typeDecl != nil {
				schema.Types = append(schema.Types, *typeDecl)
			}
			if procDecl != nil {
				schema.Procedures = append(schema.Procedures, *procDecl)
			}
		case token.TYPE:
			typeDecl := p.parseTypeDeclaration("")
			if typeDecl != nil {
				schema.Types = append(schema.Types, *typeDecl)
			}
		case token.PROC:
			procDecl := p.parseProcDeclaration("")
			if procDecl != nil {
				schema.Procedures = append(schema.Procedures, *procDecl)
			}
		}

		p.readNextToken()
	}

	if len(p.errors) > 0 {
		return schema, p.errors, p.errors[0]
	}
	return schema, nil, nil
}

// parseVersion parses the version token, should exist exactly once in the input, otherwise
// an error is added to the parser's errors list.
func (p *Parser) parseVersion(currSchema ast.Schema) ast.Version {
	if !p.expectToken(token.VERSION) {
		return ast.Version{}
	}

	if currSchema.Version.IsSet {
		p.appendError("version already set")
		return ast.Version{}
	}

	p.readNextToken()

	if p.currentToken.Type != token.INT {
		p.appendError("version expected to be an integer")
		return ast.Version{}
	}

	versionNumber, err := strconv.Atoi(p.currentToken.Literal)
	if err != nil {
		p.appendError(fmt.Sprintf("version number is not a valid integer: %s", err.Error()))
		return ast.Version{}
	}

	return ast.Version{
		IsSet: true,
		Value: versionNumber,
	}
}

// parseDocstring parses a docstring token and returns the underlying field or type
// to be added to the schema.
func (p *Parser) parseDocstring() (*ast.TypeDeclaration, *ast.ProcDeclaration) {
	if !p.expectToken(token.DOCSTRING) {
		return nil, nil
	}

	docstring := p.currentToken.Literal
	p.readNextToken()
	p.skipNewlines()

	if p.currentToken.Type == token.TYPE {
		return p.parseTypeDeclaration(docstring), nil
	}

	if p.currentToken.Type == token.PROC {
		return nil, p.parseProcDeclaration(docstring)
	}

	p.appendError("docstring can be only added to type or procedure declaration")
	return nil, nil
}

// parseTypeDeclaration parses a type declaration and returns it to be added to the schema.
func (p *Parser) parseTypeDeclaration(docstring string) *ast.TypeDeclaration {
	if !p.expectToken(token.TYPE, "missing type keyword") {
		return nil
	}

	p.readNextToken()
	if !p.expectToken(token.IDENT, "missing type name") {
		return nil
	}

	// TODO: Validate type name PascalCase
	typeName := p.currentToken.Literal
	p.readNextToken()
	if !p.expectToken(token.LBRACE, "missing type opening brace") {
		return nil
	}
	p.skipNewlines()

	var fields []ast.Field
	for {
		p.readNextToken()
		p.skipNewlines()

		if p.currentToken.Type == token.RBRACE {
			break
		}
		if p.currentToken.Type == token.EOF {
			p.appendError("missing type closing brace, unexpected EOF while parsing type fields")
			return nil
		}

		field := p.parseField()
		if field != nil {
			fields = append(fields, *field)
		}
	}

	return &ast.TypeDeclaration{
		Name:   typeName,
		Doc:    docstring,
		Fields: fields,
	}
}

func (p *Parser) parseField() *ast.Field {
	p.skipNewlines()
	if !p.expectToken(token.IDENT, "missing field name") {
		return nil
	}

	fieldName := p.currentToken.Literal
	p.readNextToken()

	isOptional := false
	if p.currentToken.Type == token.QUESTION {
		isOptional = true
		p.readNextToken()
	}

	if !p.expectToken(token.COLON, "missing field type colon for "+fieldName) {
		return nil
	}
	p.readNextToken()

	typeLiteral := p.currentToken.Literal
	var fieldType ast.Type

	switch typeLiteral {
	case "string":
		fieldType = &ast.TypeString{}
	case "int":
		fieldType = &ast.TypeInt{}
	case "float":
		fieldType = &ast.TypeFloat{}
	case "boolean":
		fieldType = &ast.TypeBoolean{}
	default:
		// TODO: Validate type name PascalCase
		fieldType = &ast.TypeCustom{
			Name: typeLiteral,
		}
	}

	nextToken, eofReached := p.peekToken(1)
	if !eofReached && nextToken.Type == token.AT {
		// TODO: Parse field rules
	}

	return &ast.Field{
		Name:     fieldName,
		Optional: isOptional,
		Type:     fieldType,
	}
}

// parseProcDeclaration parses a procedure declaration and returns it to be added to the schema.
func (p *Parser) parseProcDeclaration(docstring string) *ast.ProcDeclaration {
	if !p.expectToken(token.PROC, "missing proc keyword") {
		return nil
	}

	p.readNextToken()
	if !p.expectToken(token.IDENT, "missing procedure name") {
		return nil
	}

	// TODO: Validate procedure name PascalCase
	typeName := p.currentToken.Literal
	p.readNextToken()
	if !p.expectToken(token.LBRACE, "missing procedure opening brace") {
		return nil
	}
	p.skipNewlines()

	input := ast.Input{}
	output := ast.Output{}
	metadata := ast.Metadata{}

	for {
		p.readNextToken()
		p.skipNewlines()

		if p.currentToken.Type == token.RBRACE {
			break
		}
		if p.currentToken.Type == token.EOF {
			p.appendError("missing procedure closing brace, unexpected EOF while parsing procedure children nodes")
			return nil
		}
	}

	return &ast.ProcDeclaration{
		Name:     typeName,
		Doc:      docstring,
		Input:    input,
		Output:   output,
		Metadata: metadata,
	}
}
