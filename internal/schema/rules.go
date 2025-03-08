package schema

import (
	"encoding/json"
	"strconv"

	"github.com/orsinium-labs/enum"
)

/*
	Here we define the rules structure using the following pattern:
	https://web.archive.org/web/20250226184746/https://danielmschmidt.de/posts/2024-07-22-discriminated-union-pattern-in-go/#using-a-combined-base-type-that-derives-a-specific-subtype

	Please make sure to follow the pattern and double check the code that must be duplicated.
*/

// RuleName represents the allowed field rules
type RuleName enum.Member[string]

// MarshalJSON implements the json.Marshaler interface
func (ruleName RuleName) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ruleName.Value + `"`), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (ruleName *RuleName) UnmarshalJSON(data []byte) error {
	value := string(data[1 : len(data)-1])
	*ruleName = RuleName{value}
	return nil
}

// Allowed field rule names
var (
	RuleNameOptional  = RuleName{"optional"}
	RuleNameEquals    = RuleName{"equals"}
	RuleNameContains  = RuleName{"contains"}
	RuleNameRegex     = RuleName{"regex"}
	RuleNameLength    = RuleName{"length"}
	RuleNameMinLength = RuleName{"minLength"}
	RuleNameMaxLength = RuleName{"maxLength"}
	RuleNameEnum      = RuleName{"enum"}
	RuleNameEmail     = RuleName{"email"}
	RuleNameIso8601   = RuleName{"iso8601"}
	RuleNameUuid      = RuleName{"uuid"}
	RuleNameJson      = RuleName{"json"}
	RuleNameLowercase = RuleName{"lowercase"}
	RuleNameUppercase = RuleName{"uppercase"}
	RuleNameMin       = RuleName{"min"}
	RuleNameMax       = RuleName{"max"}
)

// Rule represents the interface that all rule types must implement
// to be considered a valid rule
type Rule interface {
	RuleName() RuleName
}

type RuleOptional struct {
	Name    RuleName `json:"name"`
	Message string   `json:"message"`
}

func (rule RuleOptional) RuleName() RuleName {
	return rule.Name
}

type RuleEquals struct {
	Name    RuleName `json:"name"`
	Value   string   `json:"value"`
	Message string   `json:"message"`
}

func (rule RuleEquals) RuleName() RuleName {
	return rule.Name
}

type RuleContains struct {
	Name    RuleName `json:"name"`
	Value   string   `json:"value"`
	Message string   `json:"message"`
}

func (rule RuleContains) RuleName() RuleName {
	return rule.Name
}

type RuleRegex struct {
	Name    RuleName `json:"name"`
	Pattern string   `json:"pattern"`
	Message string   `json:"message"`
}

func (rule RuleRegex) RuleName() RuleName {
	return rule.Name
}

type RuleLength struct {
	Name    RuleName `json:"name"`
	Value   int      `json:"value"`
	Message string   `json:"message"`
}

func (rule RuleLength) RuleName() RuleName {
	return rule.Name
}

type RuleMinLength struct {
	Name    RuleName `json:"name"`
	Value   int      `json:"value"`
	Message string   `json:"message"`
}

func (rule RuleMinLength) RuleName() RuleName {
	return rule.Name
}

type RuleMaxLength struct {
	Name    RuleName `json:"name"`
	Value   int      `json:"value"`
	Message string   `json:"message"`
}

func (rule RuleMaxLength) RuleName() RuleName {
	return rule.Name
}

type RuleEnum struct {
	Name    RuleName `json:"name"`
	Values  []any    `json:"values"`
	Message string   `json:"message"`
}

func (rule RuleEnum) RuleName() RuleName {
	return rule.Name
}

type RuleEmail struct {
	Name    RuleName `json:"name"`
	Message string   `json:"message"`
}

func (rule RuleEmail) RuleName() RuleName {
	return rule.Name
}

type RuleIso8601 struct {
	Name    RuleName `json:"name"`
	Message string   `json:"message"`
}

func (rule RuleIso8601) RuleName() RuleName {
	return rule.Name
}

type RuleUuid struct {
	Name    RuleName `json:"name"`
	Message string   `json:"message"`
}

func (rule RuleUuid) RuleName() RuleName {
	return rule.Name
}

type RuleJson struct {
	Name    RuleName `json:"name"`
	Message string   `json:"message"`
}

func (rule RuleJson) RuleName() RuleName {
	return rule.Name
}

type RuleLowercase struct {
	Name    RuleName `json:"name"`
	Message string   `json:"message"`
}

func (rule RuleLowercase) RuleName() RuleName {
	return rule.Name
}

type RuleUppercase struct {
	Name    RuleName `json:"name"`
	Message string   `json:"message"`
}

func (rule RuleUppercase) RuleName() RuleName {
	return rule.Name
}

type RuleMin struct {
	Name    RuleName `json:"name"`
	Value   float64  `json:"value"`
	Message string   `json:"message"`
}

func (rule RuleMin) RuleName() RuleName {
	return rule.Name
}

type RuleMax struct {
	Name    RuleName `json:"name"`
	Value   float64  `json:"value"`
	Message string   `json:"message"`
}

func (rule RuleMax) RuleName() RuleName {
	return rule.Name
}

// RuleCatchAll represents a catch-all rule that can be used to parse
// any custom rule and later transform it into a specific rule
type RuleCatchAll struct {
	Name    RuleName `json:"name"`
	Value   string   `json:"value"`
	Message string   `json:"message"`
}

// ToSpecificRule converts the catch-all rule into a specific rule
// to be used as needed
func (ruleCatchAll RuleCatchAll) ToSpecificRule() Rule {
	switch ruleCatchAll.Name {
	case RuleNameOptional:
		return RuleOptional{
			Name:    ruleCatchAll.Name,
			Message: ruleCatchAll.Message,
		}
	case RuleNameEquals:
		return RuleEquals{
			Name:    ruleCatchAll.Name,
			Value:   ruleCatchAll.Value,
			Message: ruleCatchAll.Message,
		}
	case RuleNameContains:
		return RuleContains{
			Name:    ruleCatchAll.Name,
			Value:   ruleCatchAll.Value,
			Message: ruleCatchAll.Message,
		}
	case RuleNameRegex:
		return RuleRegex{
			Name:    ruleCatchAll.Name,
			Pattern: ruleCatchAll.Value,
			Message: ruleCatchAll.Message,
		}
	case RuleNameLength:
		// Convert string value to int
		value, _ := strconv.Atoi(ruleCatchAll.Value)
		return RuleLength{
			Name:    ruleCatchAll.Name,
			Value:   value,
			Message: ruleCatchAll.Message,
		}
	case RuleNameMinLength:
		// Convert string value to int
		value, _ := strconv.Atoi(ruleCatchAll.Value)
		return RuleMinLength{
			Name:    ruleCatchAll.Name,
			Value:   value,
			Message: ruleCatchAll.Message,
		}
	case RuleNameMaxLength:
		// Convert string value to int
		value, _ := strconv.Atoi(ruleCatchAll.Value)
		return RuleMaxLength{
			Name:    ruleCatchAll.Name,
			Value:   value,
			Message: ruleCatchAll.Message,
		}
	case RuleNameEnum:
		// Note: This would need proper JSON parsing to handle the array
		var values []any
		_ = json.Unmarshal([]byte(ruleCatchAll.Value), &values)
		return RuleEnum{
			Name:    ruleCatchAll.Name,
			Values:  values,
			Message: ruleCatchAll.Message,
		}
	case RuleNameEmail:
		return RuleEmail{
			Name:    ruleCatchAll.Name,
			Message: ruleCatchAll.Message,
		}
	case RuleNameIso8601:
		return RuleIso8601{
			Name:    ruleCatchAll.Name,
			Message: ruleCatchAll.Message,
		}
	case RuleNameUuid:
		return RuleUuid{
			Name:    ruleCatchAll.Name,
			Message: ruleCatchAll.Message,
		}
	case RuleNameJson:
		return RuleJson{
			Name:    ruleCatchAll.Name,
			Message: ruleCatchAll.Message,
		}
	case RuleNameLowercase:
		return RuleLowercase{
			Name:    ruleCatchAll.Name,
			Message: ruleCatchAll.Message,
		}
	case RuleNameUppercase:
		return RuleUppercase{
			Name:    ruleCatchAll.Name,
			Message: ruleCatchAll.Message,
		}
	case RuleNameMin:
		// Convert string value to float64
		value, _ := strconv.ParseFloat(ruleCatchAll.Value, 64)
		return RuleMin{
			Name:    ruleCatchAll.Name,
			Value:   value,
			Message: ruleCatchAll.Message,
		}
	case RuleNameMax:
		// Convert string value to float64
		value, _ := strconv.ParseFloat(ruleCatchAll.Value, 64)
		return RuleMax{
			Name:    ruleCatchAll.Name,
			Value:   value,
			Message: ruleCatchAll.Message,
		}
	default:
		return nil
	}
}
