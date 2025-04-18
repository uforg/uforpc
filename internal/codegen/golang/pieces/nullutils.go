package pieces

import (
	"encoding/json"
	"time"
)

/** START FROM HERE **/

// -----------------------------------------------------------------------------
// Null utility type
// -----------------------------------------------------------------------------

type (
	// Null represents a value that can be null in JSON
	Null[T any] struct {
		Value T    // The actual value
		Valid bool // Whether the value is not null
	}

	// NullString is a string that can be null in JSON
	NullString = Null[string]

	// NullInt is an int that can be null in JSON
	NullInt = Null[int]

	// NullFloat64 is a float64 that can be null in JSON
	NullFloat64 = Null[float64]

	// NullBool is a bool that can be null in JSON
	NullBool = Null[bool]

	// NullTime is a time.Time that can be null in JSON
	NullTime = Null[time.Time]
)

// UnmarshalJSON implements json.Unmarshaler for Null[T]
func (n *Null[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Valid = false
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	n.Value = value
	n.Valid = true
	return nil
}

// MarshalJSON implements json.Marshaler for Null[T]
func (n Null[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Value)
}
