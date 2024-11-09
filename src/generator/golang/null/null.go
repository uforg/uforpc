package main

import "encoding/json"

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

	// NullBool is a bool that can be null in JSON
	NullBool = Null[bool]

	// NullInt is an int that can be null in JSON
	NullInt = Null[int]

	// NullFloat64 is a float64 that can be null in JSON
	NullFloat64 = Null[float64]
)

// UnmarshalJSON implements json.Unmarshaler
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

// MarshalJSON implements json.Marshaler
func (n Null[T]) MarshalJSON() ([]byte, error) {
	if !n.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(n.Value)
}
