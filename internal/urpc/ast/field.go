package ast

// Field represents a field in a custom type or procedure input/output.
type Field struct {
	Name     string       `parser:"@Ident"`
	Optional bool         `parser:"@(Question)?"`
	Type     FieldType    `parser:"Colon @@"`
	Rules    []*FieldRule `parser:"@@*"`
}

// FieldType represents the type of a field. If the field is an array, the Depth
// represents the depth of the array otherwise it is 0.
type FieldType struct {
	Base  *FieldTypeBase `parser:"@@"`
	Depth FieldTypeDepth `parser:"@((LBracket RBracket)*)"`
}

// FieldTypeDepth represents the depth of an array.
type FieldTypeDepth int

func (a *FieldTypeDepth) Capture(values []string) error {
	count := 0
	for i := range len(values) {
		if values[i] == "[" && values[i+1] == "]" {
			count++
		}
	}

	*a = FieldTypeDepth(count)
	return nil
}

// FieldTypeBase represents the base type of a field. If the field is a primitive
// or custom type, the Named field is set. If the field is an inline object, the Object
// field is set.
type FieldTypeBase struct {
	Named  *string          `parser:"@(Ident | String | Int | Float | Boolean | Datetime)"`
	Object *FieldTypeObject `parser:"| @@"`
}

// FieldTypeObject represents an inline object type.
type FieldTypeObject struct {
	Fields []*Field `parser:"LBrace @@+ RBrace"`
}

// FieldRule represents a rule applied to a field.
type FieldRule struct {
	Name string        `parser:"At @Ident"`
	Body FieldRuleBody `parser:"(LParen @@ RParen)?"`
}

// FieldRuleBody represents the body of a rule applied to a field.
type FieldRuleBody struct {
	ParamSingle *string  `parser:"@(StringLiteral | IntLiteral | FloatLiteral | TrueLiteral | FalseLiteral)?"`
	ParamList   []string `parser:"(LBracket @(StringLiteral | IntLiteral | FloatLiteral | TrueLiteral | FalseLiteral) (Comma @(StringLiteral | IntLiteral | FloatLiteral | TrueLiteral | FalseLiteral))* RBracket)?"`
	Error       string   `parser:"(Comma? Error Colon @StringLiteral)?"`
}
