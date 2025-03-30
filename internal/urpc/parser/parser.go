package parser

import (
	"fmt"
	"strconv"

	"github.com/uforg/uforpc/internal/urpc/ast"
	"github.com/uforg/uforpc/internal/urpc/lexer"
	"github.com/uforg/uforpc/internal/urpc/token"
	"github.com/uforg/uforpc/internal/util/strutil"
)

// ParserError is an error type that contains a message and a position from the parser.
type ParserError struct {
	Pos ast.Position
	Msg string
}

func (e ParserError) String() string {
	return fmt.Sprintf("%s Ln %d Col %d: %s", e.Pos.Filename, e.Pos.StartLine, e.Pos.StartCol, e.Msg)
}

func (e ParserError) Error() string {
	return e.String()
}

// Parser analyzes tokens from a Lexer and constructs an AST representation of the URPC schema.
// It tracks errors encountered during parsing and maintains the current parsing position.
type Parser struct {
	lex               *lexer.Lexer
	tokens            []token.Token
	errors            []ParserError
	maxIndex          int
	currentIndex      int
	currentIndexIsEOF bool
	currentToken      token.Token
}

// New creates and initializes a new Parser from a Lexer.
// It reads all tokens from the lexer and sets up the initial parsing state.
func New(lex *lexer.Lexer) *Parser {
	p := &Parser{}

	p.lex = lex
	p.tokens = lex.ReadTokens()
	p.errors = []ParserError{}
	p.maxIndex = len(p.tokens) - 1

	p.currentIndex = 0
	if p.maxIndex <= 0 {
		p.currentIndexIsEOF = true
	} else {
		p.currentIndexIsEOF = false
	}
	if p.currentIndexIsEOF {
		p.currentToken = token.Token{
			Type:        token.EOF,
			Literal:     "",
			FileName:    p.lex.FileName,
			LineStart:   p.lex.CurrentLine,
			ColumnStart: p.lex.CurrentColumn,
			LineEnd:     p.lex.CurrentLine,
			ColumnEnd:   p.lex.CurrentColumn,
		}
	} else {
		p.currentToken = p.tokens[p.currentIndex]
	}

	return p
}

// readNextToken advances the parser to the next token in the token list.
// If the end of the token list is reached, it sets the current token to EOF.
func (p *Parser) readNextToken() {
	if p.currentIndexIsEOF {
		return
	}

	p.currentIndex++
	if p.currentIndex > p.maxIndex {
		p.currentIndexIsEOF = true
		p.currentToken = token.Token{
			Type:        token.EOF,
			Literal:     "",
			FileName:    p.lex.FileName,
			LineStart:   p.lex.CurrentLine,
			ColumnStart: p.lex.CurrentColumn,
			LineEnd:     p.lex.CurrentLine,
			ColumnEnd:   p.lex.CurrentColumn,
		}
	} else {
		p.currentToken = p.tokens[p.currentIndex]
	}
}

// appendError adds a new error message to the parser's error list.
// The error message is formatted with the current token's location information.
func (p *Parser) appendError(pos ast.Position, message string) {
	p.errors = append(p.errors, ParserError{
		Pos: pos,
		Msg: message,
	})
}

// appendErrorForToken adds a new error message to the parser's error list.
// The error message is formatted with the given token's location information.
func (p *Parser) appendErrorForToken(tok token.Token, message string) {
	p.appendError(p.posFromToken(tok), message)
}

// appendErrorForCurrentToken adds a new error message to the parser's error list.
// The error message is formatted with the current token's location information.
func (p *Parser) appendErrorForCurrentToken(message string) {
	p.appendErrorForToken(p.currentToken, message)
}

// expectToken checks if the current token matches the expected type.
// If not, it adds an error to the parser's error list and returns false.
func (p *Parser) expectToken(expectedType token.TokenType, message ...string) bool {
	if p.currentToken.Type != expectedType {
		msg := fmt.Sprintf("expected token \"%s\", got \"%s\"", expectedType, p.currentToken.Type)
		if len(message) > 0 {
			msg += fmt.Sprintf(": %s", message[0])
		}
		p.appendErrorForCurrentToken(msg)
		return false
	}
	return true
}

// posFromToken creates a Position from a single token with the same start and end.
func (p *Parser) posFromToken(tok token.Token) ast.Position {
	return ast.Position{
		Filename:  tok.FileName,
		StartLine: tok.LineStart,
		StartCol:  tok.ColumnStart,
		EndLine:   tok.LineEnd,
		EndCol:    tok.ColumnEnd,
	}
}

// posFromTokenRange creates a Position from a range of tokens.
func (p *Parser) posFromTokenRange(startToken token.Token, endToken token.Token) ast.Position {
	return ast.Position{
		Filename:  startToken.FileName,
		StartLine: startToken.LineStart,
		StartCol:  startToken.ColumnStart,
		EndLine:   endToken.LineEnd,
		EndCol:    endToken.ColumnEnd,
	}
}

// Parse parses the input schema and returns an AST representation.
// It validates the syntax and reports any errors encountered.
func (p *Parser) Parse() (ast.Schema, []ParserError, error) {
	schema := ast.Schema{
		Pos: ast.Position{
			Filename:  p.lex.FileName,
			StartLine: 1, // Schema always starts at line 1
			StartCol:  1, // Schema always starts at column 1
			EndLine:   p.lex.CurrentLine,
			EndCol:    p.lex.CurrentColumn,
		},
	}

	for p.currentToken.Type != token.EOF {
		switch p.currentToken.Type {
		case token.VERSION:
			schema.Version = p.parseVersion(schema)
		case token.DOCSTRING:
			ruleDecl, typeDecl, procDecl := p.parseDocstring()
			if ruleDecl != nil {
				schema.CustomRules = append(schema.CustomRules, *ruleDecl)
				// Since parseCustomRuleDeclaration doesn't advance past RBRACE, we need to do it here
				p.readNextToken()
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
				// Since parseCustomRuleDeclaration doesn't advance past RBRACE, we need to do it here
				p.readNextToken()
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
		default:
			// Skip unknown tokens
			p.readNextToken()
		}
	}

	// Update the end position of the schema to the EOF token
	schema.Pos.EndLine = p.currentToken.LineStart
	schema.Pos.EndCol = p.currentToken.ColumnStart

	if len(p.errors) > 0 {
		return schema, p.errors, p.errors[0]
	}
	return schema, nil, nil
}

// parseVersion processes a version declaration in the schema.
// It validates that the version is a valid integer and is only declared once.
func (p *Parser) parseVersion(currSchema ast.Schema) ast.Version {
	if !p.expectToken(token.VERSION) {
		return ast.Version{}
	}

	startToken := p.currentToken
	if currSchema.Version.IsSet {
		p.appendErrorForCurrentToken("version already set")
		return ast.Version{}
	}

	p.readNextToken()

	if p.currentToken.Type != token.INT {
		p.appendErrorForCurrentToken("version expected to be an integer")
		return ast.Version{}
	}

	versionToken := p.currentToken
	versionNumber, err := strconv.Atoi(versionToken.Literal)
	if err != nil {
		p.appendErrorForCurrentToken(fmt.Sprintf("version number is not a valid integer: %s", err.Error()))
		return ast.Version{}
	}

	return ast.Version{
		Pos:   p.posFromTokenRange(startToken, versionToken),
		Value: versionNumber,
		IsSet: true,
	}
}

// parseDocstring handles a documentation string followed by a rule, type, or procedure declaration.
// It routes to the appropriate parser function based on what follows the docstring.
func (p *Parser) parseDocstring() (*ast.CustomRuleDecl, *ast.TypeDecl, *ast.ProcDecl) {
	if !p.expectToken(token.DOCSTRING) {
		return nil, nil, nil
	}

	docstringToken := p.currentToken
	docstring := docstringToken.Literal
	p.readNextToken()

	if p.currentToken.Type == token.RULE {
		customRule := p.parseCustomRuleDeclaration(docstring)
		return customRule, nil, nil
	}

	if p.currentToken.Type == token.TYPE {
		typeDecl := p.parseTypeDeclaration(docstring)
		return nil, typeDecl, nil
	}

	if p.currentToken.Type == token.PROC {
		procDecl := p.parseProcDeclaration(docstring)
		return nil, nil, procDecl
	}

	p.appendErrorForToken(docstringToken, "docstring can be only added to custom rule, type or procedure declaration")
	return nil, nil, nil
}

// parseCustomRuleDeclaration parses a custom validation rule declaration.
// It validates that the rule follows the correct syntax and that any referenced types exist.
func (p *Parser) parseCustomRuleDeclaration(docstring string) *ast.CustomRuleDecl {
	if !p.expectToken(token.RULE, "missing rule keyword") {
		return nil
	}

	startToken := p.currentToken
	p.readNextToken()
	if !p.expectToken(token.AT, "missing @ in rule name") {
		return nil
	}

	p.readNextToken()
	if !p.expectToken(token.IDENT, "missing rule name") {
		return nil
	}

	ruleName := p.currentToken.Literal
	if !strutil.IsCamelCase(ruleName) {
		p.appendErrorForCurrentToken(fmt.Sprintf("rule name '%s' must be in camelCase", ruleName))
	}

	p.readNextToken()
	if !p.expectToken(token.LBRACE, "missing rule opening brace") {
		return nil
	}
	p.readNextToken()

	// Initialize defaults
	var forType ast.Type
	var paramType ast.CustomRuleDeclParamType
	var errorMsg string

	// Parse rule fields
	for p.currentToken.Type != token.RBRACE {
		if p.currentToken.Type == token.EOF {
			p.appendErrorForCurrentToken("missing rule closing brace, unexpected EOF while parsing rule fields")
			return nil
		}

		switch p.currentToken.Type {
		case token.FOR:
			p.readNextToken()
			if !p.expectToken(token.COLON, "missing colon after 'for' keyword") {
				return nil
			}
			p.readNextToken()

			typeName := p.currentToken.Literal
			typePos := p.posFromToken(p.currentToken)

			switch typeName {
			case "string":
				forType = ast.TypePrimitive{
					Name: ast.PrimitiveTypeString,
					Pos:  typePos,
				}
			case "int":
				forType = ast.TypePrimitive{
					Name: ast.PrimitiveTypeInt,
					Pos:  typePos,
				}
			case "float":
				forType = ast.TypePrimitive{
					Name: ast.PrimitiveTypeFloat,
					Pos:  typePos,
				}
			case "boolean":
				forType = ast.TypePrimitive{
					Name: ast.PrimitiveTypeBoolean,
					Pos:  typePos,
				}
			default:
				if !strutil.IsPascalCase(typeName) {
					p.appendErrorForCurrentToken(fmt.Sprintf("custom type name '%s' must be in PascalCase", typeName))
				}
				forType = ast.TypeCustom{
					Name: typeName,
					Pos:  typePos,
				}
			}

			p.readNextToken()

			// Check if it's an array type
			if p.currentToken.Type == token.LBRACKET {
				arrayTokenPos := p.posFromToken(p.currentToken)
				p.readNextToken()
				if !p.expectToken(token.RBRACKET, "missing closing bracket in type") {
					return nil
				}

				// Get next token position for exact end position
				p.readNextToken()
				arrayTokenPos.EndLine = p.currentToken.LineStart
				arrayTokenPos.EndCol = p.currentToken.ColumnStart

				// Move back so current token isn't changed
				p.currentIndex--
				p.currentToken = p.tokens[p.currentIndex]

				forType = ast.TypeArray{
					ElementsType: forType,
					Pos:          arrayTokenPos,
				}
				p.readNextToken()
			}

		case token.PARAM:
			paramTokenPos := p.posFromToken(p.currentToken)
			p.readNextToken()
			if !p.expectToken(token.COLON, "missing colon after 'param' keyword") {
				return nil
			}
			p.readNextToken()

			// Check if it's an array type
			isArray := false
			var primitiveType ast.PrimitiveType

			// Parse the type
			switch p.currentToken.Literal {
			case "string":
				primitiveType = ast.PrimitiveTypeString
			case "int":
				primitiveType = ast.PrimitiveTypeInt
			case "float":
				primitiveType = ast.PrimitiveTypeFloat
			case "boolean":
				primitiveType = ast.PrimitiveTypeBoolean
			default:
				p.appendErrorForCurrentToken(fmt.Sprintf(`invalid "%s" param type, must be one of "string", "int", "float", "boolean" or array of one of them`, p.currentToken.Literal))
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
				// Get next token position for exact end
				p.readNextToken()
				paramTokenPos.EndLine = p.currentToken.LineStart
				paramTokenPos.EndCol = p.currentToken.ColumnStart

				// Move back to stay on current token
				p.currentIndex--
				p.currentToken = p.tokens[p.currentIndex]

				p.readNextToken()
			} else {
				// Update end position for non-array types
				paramTokenPos.EndLine = p.currentToken.LineStart
				paramTokenPos.EndCol = p.currentToken.ColumnStart
			}

			paramType = ast.CustomRuleDeclParamType{
				IsArray: isArray,
				Type:    primitiveType,
				Pos:     paramTokenPos,
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
			p.appendErrorForCurrentToken(fmt.Sprintf("unexpected token %s in custom rule declaration, expected 'for', 'param' or 'error'", p.currentToken.Type))
			return nil
		}
	}

	// Include the closing brace in the token range
	braceToken := p.currentToken

	rulePos := p.posFromTokenRange(startToken, braceToken)

	// Return the result but don't advance the token again, leave on the RBRACE
	// The caller is responsible for reading the next token after RBRACE
	customRule := &ast.CustomRuleDecl{
		Doc:   docstring,
		Name:  ruleName,
		For:   forType,
		Param: paramType,
		Error: errorMsg,
		Pos:   rulePos,
	}

	return customRule
}

// parseTypeDeclaration processes a type declaration in the schema.
// It validates the type name, fields, and their validation rules.
func (p *Parser) parseTypeDeclaration(docstring string) *ast.TypeDecl {
	if !p.expectToken(token.TYPE, "missing type keyword") {
		return nil
	}

	startToken := p.currentToken
	p.readNextToken()
	if !p.expectToken(token.IDENT, "missing type name") {
		return nil
	}

	typeName := p.currentToken.Literal
	if !strutil.IsPascalCase(typeName) {
		p.appendErrorForCurrentToken(fmt.Sprintf("type name '%s' must be in PascalCase", typeName))
		return nil
	}

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
			p.appendErrorForCurrentToken("missing type closing brace, unexpected EOF while parsing type fields")
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

	// Include the closing brace in the position
	braceToken := p.currentToken
	p.readNextToken()

	typePos := p.posFromTokenRange(startToken, braceToken)

	return &ast.TypeDecl{
		Name:   typeName,
		Doc:    docstring,
		Fields: fields,
		Pos:    typePos,
	}
}

// parseField processes a field declaration within a type, input, or output block.
// It handles the field name, type, optional flag, and validation rules.
func (p *Parser) parseField() *ast.Field {
	if !p.expectToken(token.IDENT, "missing field name") {
		return nil
	}

	startToken := p.currentToken
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
	var lastToken token.Token

	// Handle object type
	if p.currentToken.Type == token.LBRACE {
		objType := p.parseObjectType()
		if len(objType.Fields) == 0 {
			return nil
		}
		fieldType = objType
		lastToken = p.currentToken
	} else {
		typePos := p.posFromToken(p.currentToken)
		typeLiteral := p.currentToken.Literal
		lastToken = p.currentToken

		switch typeLiteral {
		case "string":
			fieldType = ast.TypePrimitive{
				Name: ast.PrimitiveTypeString,
				Pos:  typePos,
			}
		case "int":
			fieldType = ast.TypePrimitive{
				Name: ast.PrimitiveTypeInt,
				Pos:  typePos,
			}
		case "float":
			fieldType = ast.TypePrimitive{
				Name: ast.PrimitiveTypeFloat,
				Pos:  typePos,
			}
		case "boolean":
			fieldType = ast.TypePrimitive{
				Name: ast.PrimitiveTypeBoolean,
				Pos:  typePos,
			}
		default:
			if !strutil.IsPascalCase(typeLiteral) {
				p.appendErrorForCurrentToken(fmt.Sprintf("custom type name '%s' must be in PascalCase", typeLiteral))
				return nil
			}
			fieldType = ast.TypeCustom{
				Name: typeLiteral,
				Pos:  typePos,
			}
		}
		p.readNextToken()
	}

	for p.currentToken.Type == token.LBRACKET {
		arrayPos := p.posFromToken(p.currentToken)
		p.readNextToken()
		if !p.expectToken(token.RBRACKET, "missing array closing bracket") {
			return nil
		}

		// Read next token for exact end position
		p.readNextToken()
		arrayPos.EndLine = p.currentToken.LineStart
		arrayPos.EndCol = p.currentToken.ColumnStart
		lastToken = p.currentToken

		fieldType = ast.TypeArray{
			ElementsType: fieldType,
			Pos:          arrayPos,
		}
	}

	// Parse field rules
	var fieldValidationRules []ast.ValidationRule

	for p.currentToken.Type == token.AT {
		rule := p.parseFieldRule()
		if rule != nil {
			fieldValidationRules = append(fieldValidationRules, *rule)
			lastToken = p.currentToken
		}
	}

	fieldPos := p.posFromTokenRange(startToken, lastToken)

	return &ast.Field{
		Name:            fieldName,
		Optional:        isOptional,
		Type:            fieldType,
		ValidationRules: fieldValidationRules,
		Pos:             fieldPos,
	}
}

// parseObjectType processes an inline object type declaration.
// It handles the object's fields and their validation rules.
func (p *Parser) parseObjectType() ast.TypeObject {
	if !p.expectToken(token.LBRACE, "missing object type opening brace") {
		return ast.TypeObject{}
	}

	startToken := p.currentToken
	p.readNextToken()

	var fields []ast.Field
	for {
		if p.currentToken.Type == token.RBRACE {
			break
		}
		if p.currentToken.Type == token.EOF {
			p.appendErrorForCurrentToken("missing object type closing brace, unexpected EOF while parsing object fields")
			return ast.TypeObject{}
		}
		if !p.expectToken(token.IDENT, "missing field name") {
			return ast.TypeObject{}
		}

		field := p.parseField()
		if field != nil {
			fields = append(fields, *field)
		}
	}

	// Include the closing brace in the position
	braceToken := p.currentToken
	p.readNextToken()

	objPos := p.posFromTokenRange(startToken, braceToken)

	return ast.TypeObject{
		Fields: fields,
		Pos:    objPos,
	}
}

// parseFieldRule processes a validation rule applied to a field.
// It handles simple rules, rules with values, and rules with arrays of values.
func (p *Parser) parseFieldRule() *ast.ValidationRule {
	if !p.expectToken(token.AT, "missing field rule at") {
		return nil
	}

	startToken := p.currentToken
	p.readNextToken()
	if !p.expectToken(token.IDENT, "missing field rule name") {
		return nil
	}
	ruleName := p.currentToken.Literal
	if !strutil.IsCamelCase(ruleName) {
		p.appendErrorForCurrentToken(fmt.Sprintf("rule name '%s' must be in camelCase", ruleName))
	}

	// Default to simple rule with no parameters
	var rule ast.ValidationRule

	// Check if there are parameters (starting with parenthesis)
	p.readNextToken()
	if p.currentToken.Type != token.LPAREN {
		// No parameters, simple rule
		rule = &ast.ValidationRuleSimple{
			Name:  ruleName,
			Error: "",
			Pos:   p.posFromTokenRange(startToken, p.currentToken),
		}
		return &rule
	}

	// Process rule parameters
	p.readNextToken()

	// Special case for error-only rules
	if p.currentToken.Type == token.ERROR {
		p.readNextToken()
		if !p.expectToken(token.COLON, "missing colon after 'error' keyword") {
			return nil
		}
		p.readNextToken()
		if !p.expectToken(token.STRING, "missing error message string") {
			return nil
		}

		// Create simple rule with just an error message
		errorMsg := p.currentToken.Literal
		p.readNextToken()

		// Expect closing parenthesis
		if !p.expectToken(token.RPAREN, "missing closing parenthesis in validation rule") {
			return nil
		}
		closeParenToken := p.currentToken
		p.readNextToken()

		rule = &ast.ValidationRuleSimple{
			Name:  ruleName,
			Error: errorMsg,
			Pos:   p.posFromTokenRange(startToken, closeParenToken),
		}

		return &rule
	}

	// Handle different parameter types
	switch p.currentToken.Type {
	case token.RPAREN:
		// Empty parentheses, still a simple rule
		closeParenToken := p.currentToken
		p.readNextToken()

		rule = &ast.ValidationRuleSimple{
			Name:  ruleName,
			Error: "",
			Pos:   p.posFromTokenRange(startToken, closeParenToken),
		}

		return &rule

	case token.STRING:
		// String value
		valueStr := p.currentToken.Literal
		rule = &ast.ValidationRuleWithValue{
			Name:      ruleName,
			Value:     valueStr,
			ValueType: ast.ValidationRuleValueTypeString,
			Error:     "",
		}
		p.readNextToken()

	case token.INT:
		// Integer value
		valueStr := p.currentToken.Literal
		rule = &ast.ValidationRuleWithValue{
			Name:      ruleName,
			Value:     valueStr,
			ValueType: ast.ValidationRuleValueTypeInt,
			Error:     "",
		}
		p.readNextToken()

	case token.FLOAT:
		// Float value
		valueStr := p.currentToken.Literal
		rule = &ast.ValidationRuleWithValue{
			Name:      ruleName,
			Value:     valueStr,
			ValueType: ast.ValidationRuleValueTypeFloat,
			Error:     "",
		}
		p.readNextToken()

	case token.TRUE, token.FALSE:
		// Boolean value
		valueStr := p.currentToken.Literal
		rule = &ast.ValidationRuleWithValue{
			Name:      ruleName,
			Value:     valueStr,
			ValueType: ast.ValidationRuleValueTypeBoolean,
			Error:     "",
		}
		p.readNextToken()

	case token.LBRACKET:
		// Array of values
		var values []string
		var valueType ast.ValidationRuleValueType

		p.readNextToken()

		// Parse array values
		firstValue := true
		for p.currentToken.Type != token.RBRACKET {
			if !firstValue {
				if !p.expectToken(token.COMMA, "missing comma between array values") {
					return nil
				}
				p.readNextToken()
			}

			switch p.currentToken.Type {
			case token.STRING:
				if firstValue {
					valueType = ast.ValidationRuleValueTypeString
				} else if valueType != ast.ValidationRuleValueTypeString {
					p.appendErrorForCurrentToken("mixed types in validation rule array")
					return nil
				}
			case token.INT:
				if firstValue {
					valueType = ast.ValidationRuleValueTypeInt
				} else if valueType != ast.ValidationRuleValueTypeInt {
					p.appendErrorForCurrentToken("mixed types in validation rule array")
					return nil
				}
			case token.FLOAT:
				if firstValue {
					valueType = ast.ValidationRuleValueTypeFloat
				} else if valueType != ast.ValidationRuleValueTypeFloat {
					p.appendErrorForCurrentToken("mixed types in validation rule array")
					return nil
				}
			case token.TRUE, token.FALSE:
				if firstValue {
					valueType = ast.ValidationRuleValueTypeBoolean
				} else if valueType != ast.ValidationRuleValueTypeBoolean {
					p.appendErrorForCurrentToken("mixed types in validation rule array")
					return nil
				}
			default:
				p.appendErrorForCurrentToken(fmt.Sprintf("unexpected token %s in validation rule array", p.currentToken.Type))
				return nil
			}

			values = append(values, p.currentToken.Literal)
			firstValue = false
			p.readNextToken()
		}

		// Include the closing bracket in the position
		rule = &ast.ValidationRuleWithArray{
			Name:      ruleName,
			Values:    values,
			ValueType: valueType,
			Error:     "",
		}
		p.readNextToken()

	default:
		p.appendErrorForCurrentToken(fmt.Sprintf("unexpected token %s in validation rule parameter", p.currentToken.Type))
		return nil
	}

	// Check for comma followed by error message
	if p.currentToken.Type == token.COMMA {
		p.readNextToken()
		if p.currentToken.Type == token.ERROR {
			p.readNextToken()
			if !p.expectToken(token.COLON, "missing colon after 'error' keyword") {
				return nil
			}
			p.readNextToken()
			if !p.expectToken(token.STRING, "missing error message string") {
				return nil
			}

			errorMsg := p.currentToken.Literal

			// Update the error message based on rule type
			switch r := rule.(type) {
			case *ast.ValidationRuleSimple:
				r.Error = errorMsg
			case *ast.ValidationRuleWithValue:
				r.Error = errorMsg
			case *ast.ValidationRuleWithArray:
				r.Error = errorMsg
			}

			p.readNextToken()
		} else {
			p.appendErrorForCurrentToken(fmt.Sprintf("unexpected token %s after comma in validation rule, expected 'error'", p.currentToken.Type))
			return nil
		}
	}

	// Expect closing parenthesis
	if !p.expectToken(token.RPAREN, "missing closing parenthesis in validation rule") {
		return nil
	}

	closeParenToken := p.currentToken
	p.readNextToken()

	// Set the position for the rule
	rulePos := p.posFromTokenRange(startToken, closeParenToken)

	switch r := rule.(type) {
	case *ast.ValidationRuleSimple:
		r.Pos = rulePos
	case *ast.ValidationRuleWithValue:
		r.Pos = rulePos
	case *ast.ValidationRuleWithArray:
		r.Pos = rulePos
	}

	return &rule
}

// parseProcDeclaration processes a procedure declaration in the schema.
// It validates the procedure name, input, output, and metadata sections.
func (p *Parser) parseProcDeclaration(docstring string) *ast.ProcDecl {
	if !p.expectToken(token.PROC, "missing proc keyword") {
		return nil
	}

	startToken := p.currentToken
	p.readNextToken()
	if !p.expectToken(token.IDENT, "missing procedure name") {
		return nil
	}

	procName := p.currentToken.Literal
	if !strutil.IsPascalCase(procName) {
		p.appendErrorForCurrentToken(fmt.Sprintf("procedure name '%s' must be in PascalCase", procName))
	}

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

	// Read next token to enter the loop
	p.readNextToken()

	for {
		if p.currentToken.Type == token.RBRACE {
			break
		}
		if p.currentToken.Type == token.EOF {
			p.appendErrorForCurrentToken("missing procedure closing brace, unexpected EOF while parsing procedure children nodes")
			return nil
		}
		if p.currentToken.Type == token.INPUT {
			if inputSet {
				p.appendErrorForCurrentToken("input already set for procedure " + procName)
				continue
			}

			inputRes := p.parseProcInput()
			if inputRes != nil {
				input = *inputRes
				inputSet = true
			}
		} else if p.currentToken.Type == token.OUTPUT {
			if outputSet {
				p.appendErrorForCurrentToken("output already set for procedure " + procName)
				continue
			}

			outputRes := p.parseProcOutput()
			if outputRes != nil {
				output = *outputRes
				outputSet = true
			}
		} else if p.currentToken.Type == token.META {
			if metaSet {
				p.appendErrorForCurrentToken("meta already set for procedure " + procName)
				continue
			}

			metaRes := p.parseProcMeta()
			if metaRes != nil {
				meta = *metaRes
				metaSet = true
			}
		} else {
			// Skip unknown tokens and continue
			p.readNextToken()
		}
	}

	// Include the closing brace in the position
	braceToken := p.currentToken

	procPos := p.posFromTokenRange(startToken, braceToken)

	return &ast.ProcDecl{
		Name:     procName,
		Doc:      docstring,
		Input:    input,
		Output:   output,
		Metadata: meta,
		Pos:      procPos,
	}
}

// parseProcInput processes an input block within a procedure declaration.
// It handles the input fields and their validation rules.
func (p *Parser) parseProcInput() *ast.ProcInput {
	if !p.expectToken(token.INPUT, "missing input keyword") {
		return nil
	}

	startToken := p.currentToken
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
			p.appendErrorForCurrentToken("missing input closing brace, unexpected EOF while parsing input fields")
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

	// Include the closing brace in the position
	braceToken := p.currentToken
	p.readNextToken()

	inputPos := p.posFromTokenRange(startToken, braceToken)

	return &ast.ProcInput{
		Fields: fields,
		Pos:    inputPos,
	}
}

// parseProcOutput processes an output block within a procedure declaration.
// It handles the output fields and their validation rules.
func (p *Parser) parseProcOutput() *ast.ProcOutput {
	if !p.expectToken(token.OUTPUT, "missing output keyword") {
		return nil
	}

	startToken := p.currentToken
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
			p.appendErrorForCurrentToken("missing output closing brace, unexpected EOF while parsing output fields")
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

	// Include the closing brace in the position
	braceToken := p.currentToken
	p.readNextToken()

	outputPos := p.posFromTokenRange(startToken, braceToken)

	return &ast.ProcOutput{
		Fields: fields,
		Pos:    outputPos,
	}
}

// parseProcMeta processes a metadata block within a procedure declaration.
// It handles the key-value pairs that define the procedure's metadata.
func (p *Parser) parseProcMeta() *ast.ProcMeta {
	if !p.expectToken(token.META, "missing meta keyword") {
		return nil
	}

	startToken := p.currentToken
	p.readNextToken()
	if !p.expectToken(token.LBRACE, "missing meta opening brace") {
		return nil
	}

	// Read next token to enter the loop
	p.readNextToken()

	var entries []ast.ProcMetaKV
	for {
		if p.currentToken.Type == token.RBRACE {
			break
		}
		if p.currentToken.Type == token.EOF {
			p.appendErrorForCurrentToken("missing meta closing brace, unexpected EOF while parsing meta entries")
			return nil
		}

		entry := p.parseProcMetaEntry()
		if entry != nil {
			entries = append(entries, *entry)
		} else {
			// Skip unknown tokens
			p.readNextToken()
		}
	}

	// Include the closing brace in the position
	braceToken := p.currentToken
	p.readNextToken()

	metaPos := p.posFromTokenRange(startToken, braceToken)

	return &ast.ProcMeta{
		Entries: entries,
		Pos:     metaPos,
	}
}

// parseProcMetaEntry processes a single key-value pair in a procedure's metadata block.
// It validates the key and value types.
func (p *Parser) parseProcMetaEntry() *ast.ProcMetaKV {
	if !p.expectToken(token.IDENT, "missing meta key") {
		return nil
	}

	startToken := p.currentToken
	key := p.currentToken.Literal
	p.readNextToken()

	if !p.expectToken(token.COLON, "missing meta key colon for "+key) {
		return nil
	}
	p.readNextToken()

	var procMetaType ast.PrimitiveType
	valueToken := p.currentToken

	switch p.currentToken.Type {
	case token.STRING:
		procMetaType = ast.PrimitiveTypeString
	case token.INT:
		procMetaType = ast.PrimitiveTypeInt
	case token.FLOAT:
		procMetaType = ast.PrimitiveTypeFloat
	case token.TRUE, token.FALSE:
		procMetaType = ast.PrimitiveTypeBoolean
	default:
		p.appendErrorForCurrentToken(fmt.Sprintf("invalid meta type %s for key %s", p.currentToken.Type, key))
		return nil
	}

	// Advance to the next token after processing the value
	valueText := valueToken.Literal
	p.readNextToken()

	entryPos := p.posFromTokenRange(startToken, valueToken)

	return &ast.ProcMetaKV{
		Type:  procMetaType,
		Key:   key,
		Value: valueText,
		Pos:   entryPos,
	}
}
