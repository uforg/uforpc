package pieces

import (
	"encoding/json"
	"time"
)

/** START FROM HERE **/

// -----------------------------------------------------------------------------
// Optional utility type
// -----------------------------------------------------------------------------

type (
	// Optional represents a value that can be null or not present in JSON
	Optional[T any] struct {
		Present bool // Whether the value is present or not
		Value   T    // The actual value
	}

	// StringOptional is a string that can be null or not present in JSON
	StringOptional = Optional[string]

	// IntOptional is an int that can be null or not present in JSON
	IntOptional = Optional[int]

	// Float64Optional is a float64 that can be null or not present in JSON
	Float64Optional = Optional[float64]

	// BoolOptional is a bool that can be null or not present in JSON
	BoolOptional = Optional[bool]

	// TimeOptional is a time.Time that can be null or not present in JSON
	TimeOptional = Optional[time.Time]
)

// UnmarshalJSON implements json.Unmarshaler for Optional[T]
func (n *Optional[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Present = false
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	n.Value = value
	n.Present = true
	return nil
}

// MarshalJSON implements json.Marshaler for Optional[T]
func (n Optional[T]) MarshalJSON() ([]byte, error) {
	if !n.Present {
		return []byte("null"), nil
	}
	return json.Marshal(n.Value)
}
