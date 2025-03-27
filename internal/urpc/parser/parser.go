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
	p.readNextToken()
	p.skipNewlines()

	var fields []ast.Field
	for {
		if p.currentToken.Type == token.RBRACE {
			break
		}
		if p.currentToken.Type == token.EOF {
			p.appendError("missing type closing brace, unexpected EOF while parsing type fields")
			return nil
		}
		if !p.expectToken(token.IDENT, "missing field name") {
			return nil
		}

		field := p.parseField()
		if field != nil {
			fields = append(fields, *field)
		}
		p.skipNewlines()
	}

	return &ast.TypeDeclaration{
		Name:   typeName,
		Doc:    docstring,
		Fields: fields,
	}
}

func (p *Parser) parseField() *ast.Field {
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

	p.readNextToken()
	if p.currentToken.Type == token.LBRACKET {
		fieldType = &ast.TypeArray{
			ArrayType: fieldType,
		}
		p.readNextToken()

		if !p.expectToken(token.RBRACKET, "missing array closing bracket") {
			return nil
		}
		p.readNextToken()
	}

	// Parse field rules
	var fieldValidationRules []ast.ValidationRule
	p.skipNewlines()

	for p.currentToken.Type == token.AT {
		rule := p.parseFieldRule()
		if rule != nil {
			fieldValidationRules = append(fieldValidationRules, *rule)
		}
		p.skipNewlines()
	}

	return &ast.Field{
		Name:            fieldName,
		Optional:        isOptional,
		Type:            fieldType,
		ValidationRules: fieldValidationRules,
	}
}

func (p *Parser) parseFieldRule() *ast.ValidationRule {
	p.skipNewlines()
	if !p.expectToken(token.AT, "missing field rule at") {
		return nil
	}

	p.readNextToken()
	if !p.expectToken(token.IDENT, "missing field rule name") {
		return nil
	}
	ruleName := p.currentToken.Literal

	// Default to simple rule with no parameters
	var rule ast.ValidationRule = &ast.ValidationRuleSimple{
		RuleName:     ruleName,
		ErrorMessage: "",
	}

	// Check if there are parameters (starting with parenthesis)
	p.readNextToken()
	if p.currentToken.Type != token.LPAREN {
		return &rule
	}

	// Process rule parameters
	p.readNextToken()

	// Handle different parameter types
	switch p.currentToken.Type {
	case token.RPAREN:
		// Empty parentheses, still a simple rule
		p.readNextToken()
		return &rule

	case token.STRING:
		// String value
		valueStr := p.currentToken.Literal
		rule = &ast.ValidationRuleWithValue{
			RuleName:     ruleName,
			Value:        valueStr,
			ValueType:    ast.ValidationRuleValueTypeString,
			ErrorMessage: "",
		}
		p.readNextToken()

	case token.INT:
		// Integer value
		valueStr := p.currentToken.Literal
		rule = &ast.ValidationRuleWithValue{
			RuleName:     ruleName,
			Value:        valueStr,
			ValueType:    ast.ValidationRuleValueTypeInt,
			ErrorMessage: "",
		}
		p.readNextToken()

	case token.FLOAT:
		// Float value
		valueStr := p.currentToken.Literal
		rule = &ast.ValidationRuleWithValue{
			RuleName:     ruleName,
			Value:        valueStr,
			ValueType:    ast.ValidationRuleValueTypeFloat,
			ErrorMessage: "",
		}
		p.readNextToken()

	case token.TRUE, token.FALSE:
		// Boolean value
		valueStr := p.currentToken.Literal
		rule = &ast.ValidationRuleWithValue{
			RuleName:     ruleName,
			Value:        valueStr,
			ValueType:    ast.ValidationRuleValueTypeBoolean,
			ErrorMessage: "",
		}
		p.readNextToken()

	case token.LBRACKET:
		// Array values
		var values []string
		var valueType ast.ValidationRuleValueType = ast.ValidationRuleValueTypeString // Default

		p.readNextToken()
		for p.currentToken.Type != token.RBRACKET {
			switch p.currentToken.Type {
			case token.STRING:
				values = append(values, p.currentToken.Literal)
				valueType = ast.ValidationRuleValueTypeString
			case token.INT:
				values = append(values, p.currentToken.Literal)
				valueType = ast.ValidationRuleValueTypeInt
			case token.FLOAT:
				values = append(values, p.currentToken.Literal)
				valueType = ast.ValidationRuleValueTypeFloat
			case token.TRUE, token.FALSE:
				values = append(values, p.currentToken.Literal)
				valueType = ast.ValidationRuleValueTypeBoolean
			}

			p.readNextToken()
			if p.currentToken.Type == token.COMMA {
				p.readNextToken()
			} else if p.currentToken.Type != token.RBRACKET {
				p.appendError("expected comma or closing bracket in enum values")
				return nil
			}
		}

		rule = &ast.ValidationRuleWithArray{
			RuleName:     ruleName,
			Values:       values,
			ValueType:    valueType,
			ErrorMessage: "",
		}
		p.readNextToken()

	default:
		p.appendError(fmt.Sprintf("unexpected token %s in validation rule parameters", p.currentToken.Type))
		return nil
	}

	// Look for error message after the initial parameter
	if p.currentToken.Type == token.COMMA {
		p.readNextToken()
		if !p.expectToken(token.IDENT, "missing error keyword after comma in validation rule") {
			return nil
		}

		if p.currentToken.Literal != "error" {
			p.appendError("expected 'error' keyword after comma in validation rule")
			return nil
		}

		p.readNextToken()
		if !p.expectToken(token.COLON, "missing colon after 'error' keyword in validation rule") {
			return nil
		}

		p.readNextToken()
		if !p.expectToken(token.STRING, "missing error message string in validation rule") {
			return nil
		}

		errorMsg := p.currentToken.Literal

		// Update the error message based on rule type
		switch r := rule.(type) {
		case *ast.ValidationRuleSimple:
			r.ErrorMessage = errorMsg
		case *ast.ValidationRuleWithValue:
			r.ErrorMessage = errorMsg
		case *ast.ValidationRuleWithArray:
			r.ErrorMessage = errorMsg
		}

		p.readNextToken()
	}

	// Expect closing parenthesis
	if !p.expectToken(token.RPAREN, "missing closing parenthesis in validation rule") {
		return nil
	}
	p.readNextToken()

	return &rule
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
	procName := p.currentToken.Literal
	p.readNextToken()
	if !p.expectToken(token.LBRACE, "missing procedure opening brace") {
		return nil
	}
	p.skipNewlines()

	input := ast.ProcInput{}
	inputSet := false
	output := ast.ProcOutput{}
	outputSet := false
	meta := ast.ProcMeta{}
	metaSet := false

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
		if p.currentToken.Type == token.INPUT {
			if inputSet {
				p.appendError("input already set for procedure " + procName)
				continue
			}

			inputRes := p.parseProcInput()
			if inputRes != nil {
				input = *inputRes
				inputSet = true
			}
		}
		if p.currentToken.Type == token.OUTPUT {
			if outputSet {
				p.appendError("output already set for procedure " + procName)
				continue
			}

			outputRes := p.parseProcOutput()
			if outputRes != nil {
				output = *outputRes
				outputSet = true
			}
		}
		if p.currentToken.Type == token.META {
			if metaSet {
				p.appendError("meta already set for procedure " + procName)
				continue
			}

			metaRes := p.parseProcMeta()
			if metaRes != nil {
				meta = *metaRes
				metaSet = true
			}
		}
	}

	return &ast.ProcDeclaration{
		Name:     procName,
		Doc:      docstring,
		Input:    input,
		Output:   output,
		Metadata: meta,
	}
}

func (p *Parser) parseProcInput() *ast.ProcInput {
	if !p.expectToken(token.INPUT, "missing input keyword") {
		return nil
	}

	p.readNextToken()
	if !p.expectToken(token.LBRACE, "missing input opening brace") {
		return nil
	}
	p.readNextToken()
	p.skipNewlines()

	var fields []ast.Field
	for {
		if p.currentToken.Type == token.RBRACE {
			break
		}
		if p.currentToken.Type == token.EOF {
			p.appendError("missing input closing brace, unexpected EOF while parsing input fields")
			return nil
		}
		if !p.expectToken(token.IDENT, "missing field name") {
			return nil
		}

		field := p.parseField()
		if field != nil {
			fields = append(fields, *field)
		}
		p.skipNewlines()
	}

	return &ast.ProcInput{
		Fields: fields,
	}
}

func (p *Parser) parseProcOutput() *ast.ProcOutput {
	if !p.expectToken(token.OUTPUT, "missing output keyword") {
		return nil
	}

	p.readNextToken()
	if !p.expectToken(token.LBRACE, "missing output opening brace") {
		return nil
	}
	p.readNextToken()
	p.skipNewlines()

	var fields []ast.Field
	for {
		if p.currentToken.Type == token.RBRACE {
			break
		}
		if p.currentToken.Type == token.EOF {
			p.appendError("missing output closing brace, unexpected EOF while parsing output fields")
			return nil
		}
		if !p.expectToken(token.IDENT, "missing field name") {
			return nil
		}

		field := p.parseField()
		if field != nil {
			fields = append(fields, *field)
		}
		p.skipNewlines()
	}

	return &ast.ProcOutput{
		Fields: fields,
	}
}

func (p *Parser) parseProcMeta() *ast.ProcMeta {
	if !p.expectToken(token.META, "missing meta keyword") {
		return nil
	}

	p.readNextToken()
	if !p.expectToken(token.LBRACE, "missing meta opening brace") {
		return nil
	}
	p.skipNewlines()

	var entries []ast.ProcMetaKV
	for {
		p.readNextToken()
		p.skipNewlines()

		if p.currentToken.Type == token.RBRACE {
			break
		}
		if p.currentToken.Type == token.EOF {
			p.appendError("missing meta closing brace, unexpected EOF while parsing meta entries")
			return nil
		}

		entry := p.parseProcMetaEntry()
		if entry != nil {
			entries = append(entries, *entry)
		}
	}

	return &ast.ProcMeta{
		Entries: entries,
	}
}

func (p *Parser) parseProcMetaEntry() *ast.ProcMetaKV {
	p.skipNewlines()
	if !p.expectToken(token.IDENT, "missing meta key") {
		return nil
	}

	key := p.currentToken.Literal
	p.readNextToken()

	if !p.expectToken(token.COLON, "missing meta key colon for "+key) {
		return nil
	}
	p.readNextToken()

	var procMetaType ast.ProcMetaValueTypeName
	switch p.currentToken.Type {
	case token.STRING:
		procMetaType = ast.ProcMetaValueTypeString
	case token.INT:
		procMetaType = ast.ProcMetaValueTypeInt
	case token.FLOAT:
		procMetaType = ast.ProcMetaValueTypeFloat
	case token.TRUE, token.FALSE:
		procMetaType = ast.ProcMetaValueTypeBoolean
	default:
		p.appendError(fmt.Sprintf("invalid meta type %s for key %s", p.currentToken.Type, key))
		return nil
	}

	return &ast.ProcMetaKV{
		Type:  procMetaType,
		Key:   key,
		Value: p.currentToken.Literal,
	}
}
