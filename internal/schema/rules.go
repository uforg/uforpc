package schema

import "github.com/orsinium-labs/enum"

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

// TODO: Add more rules

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
	default:
		return nil
	}
}
