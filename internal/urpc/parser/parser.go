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
			ruleDecl, typeDecl, procDecl := p.parseDocstring()
			if ruleDecl != nil {
				schema.CustomRules = append(schema.CustomRules, *ruleDecl)
			}
			if typeDecl != nil {
				schema.Types = append(schema.Types, *typeDecl)
			}
			if procDecl != nil {
				schema.Procedures = append(schema.Procedures, *procDecl)
			}
		case token.RULE:
			ruleDecl := p.parseCustomRuleDeclaration("")
			if ruleDecl != nil {
				schema.CustomRules = append(schema.CustomRules, *ruleDecl)
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
func (p *Parser) parseDocstring() (*ast.CustomRuleDeclaration, *ast.TypeDeclaration, *ast.ProcDeclaration) {
	if !p.expectToken(token.DOCSTRING) {
		return nil, nil, nil
	}

	docstring := p.currentToken.Literal
	p.readNextToken()

	if p.currentToken.Type == token.RULE {
		return p.parseCustomRuleDeclaration(docstring), nil, nil
	}

	if p.currentToken.Type == token.TYPE {
		return nil, p.parseTypeDeclaration(docstring), nil
	}

	if p.currentToken.Type == token.PROC {
		return nil, nil, p.parseProcDeclaration(docstring)
	}

	p.appendError("docstring can be only added to custom rule, type or procedure declaration")
	return nil, nil, nil
}

// parseCustomRuleDeclaration parses a custom validation rule declaration.
func (p *Parser) parseCustomRuleDeclaration(docstring string) *ast.CustomRuleDeclaration {
	if !p.expectToken(token.RULE, "missing rule keyword") {
		return nil
	}

	p.readNextToken()
	if !p.expectToken(token.AT, "missing @ in rule name") {
		return nil
	}

	p.readNextToken()
	if !p.expectToken(token.IDENT, "missing rule name") {
		return nil
	}

	ruleName := p.currentToken.Literal
	p.readNextToken()
	if !p.expectToken(token.LBRACE, "missing rule opening brace") {
		return nil
	}
	p.readNextToken()

	// Initialize defaults
	var forType ast.TypeName
	var paramType ast.CustomRuleParamType
	var errorMsg string

	// Parse rule fields
	for p.currentToken.Type != token.RBRACE {
		if p.currentToken.Type == token.EOF {
			p.appendError("missing rule closing brace, unexpected EOF while parsing rule fields")
			return nil
		}

		switch p.currentToken.Type {
		case token.FOR:
			p.readNextToken()
			if !p.expectToken(token.COLON, "missing colon after 'for' keyword") {
				return nil
			}
			p.readNextToken()

			switch p.currentToken.Literal {
			case "string":
				forType = ast.TypeNameString
			case "int":
				forType = ast.TypeNameInt
			case "float":
				forType = ast.TypeNameFloat
			case "boolean":
				forType = ast.TypeNameBoolean
			default:
				p.appendError(fmt.Sprintf("invalid 'for' type: %s", p.currentToken.Literal))
				return nil
			}
			p.readNextToken()

		case token.PARAM:
			p.readNextToken()
			if !p.expectToken(token.COLON, "missing colon after 'param' keyword") {
				return nil
			}
			p.readNextToken()

			// Check if it's an array type
			isArray := false
			var primitiveType ast.CustomRulePrimitiveType

			// Parse the type
			switch p.currentToken.Literal {
			case "string":
				primitiveType = ast.CustomRulePrimitiveTypeString
			case "int":
				primitiveType = ast.CustomRulePrimitiveTypeInt
			case "float":
				primitiveType = ast.CustomRulePrimitiveTypeFloat
			case "boolean":
				primitiveType = ast.CustomRulePrimitiveTypeBoolean
			default:
				p.appendError(fmt.Sprintf("invalid param type: %s", p.currentToken.Literal))
				return nil
			}

			p.readNextToken()

			// Check for array brackets
			if p.currentToken.Type == token.LBRACKET {
				isArray = true
				p.readNextToken()
				if !p.expectToken(token.RBRACKET, "missing closing bracket in param type") {
					return nil
				}
				p.readNextToken()
			}

			paramType = ast.CustomRuleParamType{
				IsArray: isArray,
				Type:    primitiveType,
			}

		case token.ERROR:
			p.readNextToken()
			if !p.expectToken(token.COLON, "missing colon after 'error' keyword") {
				return nil
			}
			p.readNextToken()
			if !p.expectToken(token.STRING, "missing default error message string") {
				return nil
			}

			errorMsg = p.currentToken.Literal
			p.readNextToken()

		default:
			p.appendError(fmt.Sprintf("unexpected token %s in custom rule declaration, expected 'for', 'param' or 'error'", p.currentToken.Type))
			return nil
		}
	}

	// Validate required fields
	if forType == "" {
		p.appendError("missing required 'for' field in custom rule declaration")
		return nil
	}

	return &ast.CustomRuleDeclaration{
		Doc:      docstring,
		Name:     ruleName,
		For:      forType,
		Param:    paramType,
		ErrorMsg: errorMsg,
	}
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

	var fieldType ast.Type

	// Handle object type
	if p.currentToken.Type == token.LBRACE {
		fieldType = p.parseObjectType()
		if fieldType == nil {
			return nil
		}
	} else {
		typeLiteral := p.currentToken.Literal

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
	}

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

	for p.currentToken.Type == token.AT {
		rule := p.parseFieldRule()
		if rule != nil {
			fieldValidationRules = append(fieldValidationRules, *rule)
		}
	}

	return &ast.Field{
		Name:            fieldName,
		Optional:        isOptional,
		Type:            fieldType,
		ValidationRules: fieldValidationRules,
	}
}

// parseObjectType parses an inline object type.
//
// The parser expects the current token to be an opening brace.
// Returns a TypeObject that contains the fields of the inline object.
func (p *Parser) parseObjectType() ast.Type {
	if !p.expectToken(token.LBRACE, "missing object type opening brace") {
		return nil
	}
	p.readNextToken()

	var fields []ast.Field
	for {
		if p.currentToken.Type == token.RBRACE {
			break
		}
		if p.currentToken.Type == token.EOF {
			p.appendError("missing object type closing brace, unexpected EOF while parsing object fields")
			return nil
		}
		if !p.expectToken(token.IDENT, "missing field name") {
			return nil
		}

		field := p.parseField()
		if field != nil {
			fields = append(fields, *field)
		}
	}

	// Skip the closing brace
	p.readNextToken()

	return &ast.TypeObject{
		Fields: fields,
	}
}

func (p *Parser) parseFieldRule() *ast.ValidationRule {
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

	case token.IDENT, token.ERROR:
		// Check if it's directly starting with an "error" parameter
		isErrorParam := (p.currentToken.Type == token.ERROR) ||
			(p.currentToken.Type == token.IDENT && p.currentToken.Literal == "error")

		if isErrorParam {
			p.readNextToken()
			if !p.expectToken(token.COLON, "missing colon after 'error' keyword in validation rule") {
				return nil
			}

			p.readNextToken()
			if !p.expectToken(token.STRING, "missing error message string in validation rule") {
				return nil
			}

			errorMsg := p.currentToken.Literal
			rule = &ast.ValidationRuleSimple{
				RuleName:     ruleName,
				ErrorMessage: errorMsg,
			}
			p.readNextToken()
		} else {
			p.appendError(fmt.Sprintf("unexpected identifier %s in validation rule parameters", p.currentToken.Literal))
			return nil
		}

	default:
		p.appendError(fmt.Sprintf("unexpected token %s in validation rule parameters", p.currentToken.Type))
		return nil
	}

	// Look for additional parameters (error message) after comma
	if p.currentToken.Type == token.COMMA {
		p.readNextToken()

		// Check for error parameter
		isErrorParam := (p.currentToken.Type == token.ERROR) ||
			(p.currentToken.Type == token.IDENT && p.currentToken.Literal == "error")

		if isErrorParam {
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
		} else {
			p.appendError(fmt.Sprintf("unexpected token %s after comma in validation rule, expected 'error'", p.currentToken.Type))
			return nil
		}
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

	input := ast.ProcInput{}
	inputSet := false
	output := ast.ProcOutput{}
	outputSet := false
	meta := ast.ProcMeta{}
	metaSet := false

	for {
		p.readNextToken()

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

	var entries []ast.ProcMetaKV
	for {
		p.readNextToken()

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
