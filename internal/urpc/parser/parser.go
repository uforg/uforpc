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
		}

		p.readNextToken()
	}

	if len(p.errors) > 0 {
		return schema, p.errors, p.errors[0]
	}
	return schema, nil, nil
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
