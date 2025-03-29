package ast

// CustomRuleDecl represents a custom validation rule declaration in the URPC schema.
type CustomRuleDecl struct {
	Doc   string
	Name  string
	For   Type
	Param CustomRuleDeclParamType
	Error string
}

func (c *CustomRuleDecl) NodeType() NodeType { return NodeTypeCustomRuleDecl }

// CustomRuleDeclParamType represents the allowed parameter types for custom rules.
type CustomRuleDeclParamType struct {
	IsArray bool
	Type    PrimitiveType // Only primitive types are allowed
}
