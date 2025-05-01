package schema

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"

	"github.com/orsinium-labs/enum"
)

///////////
// ENUMS //
///////////

// PrimitiveType represents the primitive type names defined in the URPC specification.
type PrimitiveType enum.Member[string]

func (f PrimitiveType) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}

func (f *PrimitiveType) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &f.Value)
}

var (
	PrimitiveTypeString   = PrimitiveType{Value: "string"}
	PrimitiveTypeInt      = PrimitiveType{Value: "int"}
	PrimitiveTypeFloat    = PrimitiveType{Value: "float"}
	PrimitiveTypeBoolean  = PrimitiveType{Value: "boolean"}
	PrimitiveTypeDatetime = PrimitiveType{Value: "datetime"}
)

// ParamPrimitiveType represents the primitive type names allowed in rule parameters.
type ParamPrimitiveType enum.Member[string]

func (f ParamPrimitiveType) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}

func (f *ParamPrimitiveType) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &f.Value)
}

var (
	ParamPrimitiveTypeString  = ParamPrimitiveType{Value: "string"}
	ParamPrimitiveTypeInt     = ParamPrimitiveType{Value: "int"}
	ParamPrimitiveTypeFloat   = ParamPrimitiveType{Value: "float"}
	ParamPrimitiveTypeBoolean = ParamPrimitiveType{Value: "boolean"}
)

////////////////////
// Main Structure //
////////////////////

// Node represents a generic node in the URPC schema structure.
// All specific node types (DocNode, RuleNode, etc.) implement this interface.
type Node interface {
	NodeKind() string
}

// Schema represents the root of the intermediate JSON structure.
type Schema struct {
	// Version is the URPC specification version (always 1 according to the schema).
	Version int `json:"version"`
	// Nodes contains the ordered list of declared elements in the schema.
	Nodes []Node `json:"nodes"`
}

// UnmarshalJSON implements custom JSON unmarshalling for Schema to handle the polymorphic Nodes array.
func (s *Schema) UnmarshalJSON(data []byte) error {
	// 1. Unmarshal into a temporary struct to get Version and raw Nodes data.
	var rawSchema struct {
		Version int               `json:"version"`
		Nodes   []json.RawMessage `json:"nodes"`
	}
	if err := json.Unmarshal(data, &rawSchema); err != nil {
		return fmt.Errorf("failed to unmarshal raw schema: %w", err)
	}

	// 2. Assign Version.
	s.Version = rawSchema.Version
	if s.Version != 1 {
		return fmt.Errorf("unsupported schema version: %d", s.Version)
	}

	// 3. Process each raw node message.
	s.Nodes = make([]Node, 0, len(rawSchema.Nodes))
	var nodeKind struct {
		Kind string `json:"kind"`
	}

	for i, rawNode := range rawSchema.Nodes {
		// 3.1. Peek at the kind field.
		if err := json.Unmarshal(rawNode, &nodeKind); err != nil {
			return fmt.Errorf("failed to determine kind of node at index %d: %w", i, err)
		}

		// 3.2. Unmarshal into the specific node type based on kind.
		var node Node
		var err error
		switch nodeKind.Kind {
		case "doc":
			var docNode NodeDoc
			err = json.Unmarshal(rawNode, &docNode)
			node = &docNode
		case "rule":
			var ruleNode NodeRule
			err = json.Unmarshal(rawNode, &ruleNode)
			node = &ruleNode
		case "type":
			var typeNode NodeType
			err = json.Unmarshal(rawNode, &typeNode)
			node = &typeNode
		case "proc":
			var procNode NodeProc
			err = json.Unmarshal(rawNode, &procNode)
			node = &procNode
		default:
			return fmt.Errorf("unknown node kind '%s' at index %d", nodeKind.Kind, i)
		}

		if err != nil {
			return fmt.Errorf("failed to unmarshal node of kind '%s' at index %d: %w", nodeKind.Kind, i, err)
		}
		s.Nodes = append(s.Nodes, node)
	}

	return nil
}

// GetDocNodes returns all DocNode instances from the schema.
func (s *Schema) GetDocNodes() []*NodeDoc {
	docNodes := []*NodeDoc{}
	for _, node := range s.Nodes {
		if docNode, ok := node.(*NodeDoc); ok {
			docNodes = append(docNodes, docNode)
		}
	}
	return docNodes
}

// GetRuleNodes returns all RuleNode instances from the schema.
func (s *Schema) GetRuleNodes() []*NodeRule {
	ruleNodes := []*NodeRule{}
	for _, node := range s.Nodes {
		if ruleNode, ok := node.(*NodeRule); ok {
			ruleNodes = append(ruleNodes, ruleNode)
		}
	}
	return ruleNodes
}

// GetTypeNodes returns all TypeNode instances from the schema.
func (s *Schema) GetTypeNodes() []*NodeType {
	typeNodes := []*NodeType{}
	for _, node := range s.Nodes {
		if typeNode, ok := node.(*NodeType); ok {
			typeNodes = append(typeNodes, typeNode)
		}
	}
	return typeNodes
}

// GetProcNodes returns all ProcNode instances from the schema.
func (s *Schema) GetProcNodes() []*NodeProc {
	procNodes := []*NodeProc{}
	for _, node := range s.Nodes {
		if procNode, ok := node.(*NodeProc); ok {
			procNodes = append(procNodes, procNode)
		}
	}
	return procNodes
}

////////////////
// Node Types //
////////////////

// NodeDoc represents a standalone documentation block.
type NodeDoc struct {
	Kind    string `json:"kind"` // Always "doc"
	Content string `json:"content"`
}

func (n *NodeDoc) NodeKind() string { return n.Kind }

// NodeRule represents the definition of a custom validation rule.
type NodeRule struct {
	Kind string `json:"kind"` // Always "rule"
	Name string `json:"name"`
	// Doc is the associated documentation string (optional).
	Doc *string `json:"doc,omitempty"`
	// Deprecated indicates if the rule is deprecated and contains the message
	// associated with the deprecation.
	Deprecated *string `json:"deprecated,omitempty"`
	// For indicates the primitive or custom type name this rule applies to.
	For *ForDefinition `json:"for"`
	// Param defines the parameter structure expected by this rule (null if none).
	Param *ParamDefinition `json:"paramDef,omitempty"` // Pointer handles null
	// Error is the default error message for the rule (optional).
	Error *string `json:"error,omitempty"`
}

func (n *NodeRule) NodeKind() string { return n.Kind }

// NodeType represents the definition of a custom data type.
type NodeType struct {
	Kind string `json:"kind"` // Always "type"
	Name string `json:"name"`
	// Doc is the associated documentation string (optional).
	Doc *string `json:"doc,omitempty"`
	// Deprecated indicates if the type is deprecated and contains the message
	// associated with the deprecation.
	Deprecated *string `json:"deprecated,omitempty"`
	// Fields is the ordered list of fields within the type.
	Fields []FieldDefinition `json:"fields"`
}

func (n *NodeType) NodeKind() string { return n.Kind }

// NodeProc represents the definition of an RPC procedure.
type NodeProc struct {
	Kind string `json:"kind"` // Always "proc"
	Name string `json:"name"`
	// Doc is the associated documentation string (optional).
	Doc *string `json:"doc,omitempty"`
	// Deprecated indicates if the procedure is deprecated and contains the message
	// associated with the deprecation.
	Deprecated *string `json:"deprecated,omitempty"`
	// Input is the ordered list of input fields for the procedure.
	Input []FieldDefinition `json:"input"`
	// Output is the ordered list of output fields for the procedure.
	Output []FieldDefinition `json:"output"`
	// Meta contains optional key-value metadata.
	Meta []MetaKeyValue `json:"meta,omitempty"`
}

func (n *NodeProc) NodeKind() string { return n.Kind }

//////////////////////////
// Auxiliary Structures //
//////////////////////////

// MetaKeyValue represents a key-value pair within the NodeProc.Meta array.
type MetaKeyValue struct {
	Key   string    `json:"key"`
	Value MetaValue `json:"value"`
}

// MetaValue holds a metadata value, using mutually exclusive fields for type safety.
type MetaValue struct {
	StringVal *string  `json:"-"` // Ignored by standard JSON marshalling/unmarshalling
	IntVal    *int64   `json:"-"`
	FloatVal  *float64 `json:"-"`
	BoolVal   *bool    `json:"-"`
}

// UnmarshalJSON implements custom unmarshalling for MetaValue.
// It validates the incoming JSON type and stores it in the corresponding field.
func (mv *MetaValue) UnmarshalJSON(data []byte) error {
	var rawValue any
	if err := json.Unmarshal(data, &rawValue); err != nil {
		return fmt.Errorf("failed to unmarshal meta value: %w", err)
	}

	// Reset fields before assigning
	mv.StringVal = nil
	mv.IntVal = nil
	mv.FloatVal = nil
	mv.BoolVal = nil

	switch v := rawValue.(type) {
	case string:
		mv.StringVal = &v
	case float64: // JSON numbers are float64
		// Check if it's an integer without loss of precision
		if math.Trunc(v) == v {
			intVal := int64(v)
			mv.IntVal = &intVal
		} else {
			mv.FloatVal = &v
		}
	case bool:
		mv.BoolVal = &v
	default:
		return fmt.Errorf("invalid meta value type: expected string, number, or boolean, got %T", v)
	}
	return nil
}

// MarshalJSON implements custom marshalling for MetaValue.
// It marshals the non-nil field back to its original JSON type.
func (mv MetaValue) MarshalJSON() ([]byte, error) {
	if mv.StringVal != nil {
		return json.Marshal(mv.StringVal)
	}
	if mv.IntVal != nil {
		return json.Marshal(mv.IntVal)
	}
	if mv.FloatVal != nil {
		return json.Marshal(mv.FloatVal)
	}
	if mv.BoolVal != nil {
		return json.Marshal(mv.BoolVal)
	}
	// Should ideally not happen if unmarshalling is correct,
	// but marshal as null if no value is set.
	return json.Marshal(nil)
}

// ForDefinition describes the type and structure expected for a rule's for clause.
type ForDefinition struct {
	Type    string `json:"type"`
	IsArray bool   `json:"isArray"`
}

// ParamDefinition describes the parameter structure expected by a rule.
type ParamDefinition struct {
	Type    ParamPrimitiveType `json:"type"`
	IsArray bool               `json:"isArray"`
}

// FieldDefinition defines a field within a type or procedure input/output.
type FieldDefinition struct {
	Name string `json:"name"`
	// TypeName holds the name if the type is named (primitive or custom). Mutually exclusive with TypeInline.
	TypeName *string `json:"typeName,omitempty"`
	// TypeInline holds the definition if the type is inline. Mutually exclusive with TypeName.
	TypeInline *InlineTypeDefinition `json:"typeInline,omitempty"`
	// IsArray indicates if the field is an array.
	IsArray bool `json:"isArray"`
	// Optional indicates if the field is optional.
	Optional bool `json:"optional"`
	// Rules is the list of validation rules applied to this field.
	Rules []AppliedRule `json:"rules"`
}

// IsNamed checks if the field definition uses a named type.
func (fd *FieldDefinition) IsNamed() bool {
	return fd.TypeName != nil
}

// IsInline checks if the field definition uses an inline type.
func (fd *FieldDefinition) IsInline() bool {
	return fd.TypeInline != nil
}

// IsBuiltInType checks if the field definition uses a built-in type.
func (fd *FieldDefinition) IsBuiltInType() bool {
	return fd.IsNamed() && slices.Contains([]string{"string", "int", "float", "boolean", "datetime"}, *fd.TypeName)
}

// IsCustomType checks if the field definition uses a custom type.
func (fd *FieldDefinition) IsCustomType() bool {
	return fd.IsNamed() && !fd.IsBuiltInType()
}

// InlineTypeDefinition represents the structure of an anonymous inline object type.
// It's used within the FieldDefinition.TypeInline field.
type InlineTypeDefinition struct {
	// Fields is the ordered list of fields within the inline type.
	Fields []FieldDefinition `json:"fields"`
}

// AppliedRule represents a validation rule applied to a field.
type AppliedRule struct {
	Rule string `json:"rule"`
	// Param holds the parameter value(s) passed to the rule instance (null if none).
	Param *AppliedParam `json:"param,omitempty"` // Pointer handles null
	// Error is the custom error message overriding the rule's default (optional).
	Error *string `json:"error,omitempty"`
}

// AppliedParam holds the actual value(s) passed to a rule instance, represented as strings.
type AppliedParam struct {
	Type ParamPrimitiveType `json:"type"`
	// IsArray indicates if the parameter was passed as an array.
	IsArray bool `json:"isArray"`
	// Value holds the single parameter value as a string (used if IsArray is false).
	Value string `json:"value,omitempty,omitzero"`
	// ArrayValues holds array parameter values as strings (used if IsArray is true).
	ArrayValues []string `json:"arrayValues,omitempty"`
}
