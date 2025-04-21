package pieces

import (
	"encoding/json"
	"testing"
	"time"
)

// TestStructure includes all nullable types
type TestStructure struct {
	Text    StringOptional  `json:"text,omitempty"`
	Number  IntOptional     `json:"number,omitempty"`
	Decimal Float64Optional `json:"decimal,omitempty"`
	Flag    BoolOptional    `json:"flag,omitempty"`
	Generic Optional[[]int] `json:"generic,omitempty"`
}

func TestStringOptional(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected StringOptional
		wantErr  bool
	}{
		{
			name:     "normal value",
			json:     `"hello"`,
			expected: StringOptional{Value: "hello", Present: true},
		},
		{
			name:     "empty value",
			json:     `""`,
			expected: StringOptional{Value: "", Present: true},
		},
		{
			name:     "null value",
			json:     `null`,
			expected: StringOptional{Present: false},
		},
		{
			name:    "invalid value",
			json:    `123`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got StringOptional
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && (got.Present != tt.expected.Present || got.Value != tt.expected.Value) {
				t.Errorf("got %+v, want %+v", got, tt.expected)
			}
		})
	}
}

func TestIntOptional(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected IntOptional
		wantErr  bool
	}{
		{
			name:     "positive number",
			json:     `42`,
			expected: IntOptional{Value: 42, Present: true},
		},
		{
			name:     "negative number",
			json:     `-42`,
			expected: IntOptional{Value: -42, Present: true},
		},
		{
			name:     "zero",
			json:     `0`,
			expected: IntOptional{Value: 0, Present: true},
		},
		{
			name:     "null value",
			json:     `null`,
			expected: IntOptional{Present: false},
		},
		{
			name:    "invalid string value",
			json:    `"123"`,
			wantErr: true,
		},
		{
			name:    "invalid float value",
			json:    `123.45`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got IntOptional
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && (got.Present != tt.expected.Present || got.Value != tt.expected.Value) {
				t.Errorf("got %+v, want %+v", got, tt.expected)
			}
		})
	}
}

func TestFloat64Optional(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected Float64Optional
		wantErr  bool
	}{
		{
			name:     "integer number",
			json:     `42`,
			expected: Float64Optional{Value: 42.0, Present: true},
		},
		{
			name:     "decimal number",
			json:     `42.5`,
			expected: Float64Optional{Value: 42.5, Present: true},
		},
		{
			name:     "negative number",
			json:     `-42.5`,
			expected: Float64Optional{Value: -42.5, Present: true},
		},
		{
			name:     "zero",
			json:     `0`,
			expected: Float64Optional{Value: 0, Present: true},
		},
		{
			name:     "null value",
			json:     `null`,
			expected: Float64Optional{Present: false},
		},
		{
			name:    "invalid value",
			json:    `"123.45"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Float64Optional
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && (got.Present != tt.expected.Present || got.Value != tt.expected.Value) {
				t.Errorf("got %+v, want %+v", got, tt.expected)
			}
		})
	}
}

func TestBoolOptional(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected BoolOptional
		wantErr  bool
	}{
		{
			name:     "true",
			json:     `true`,
			expected: BoolOptional{Value: true, Present: true},
		},
		{
			name:     "false",
			json:     `false`,
			expected: BoolOptional{Value: false, Present: true},
		},
		{
			name:     "null value",
			json:     `null`,
			expected: BoolOptional{Present: false},
		},
		{
			name:    "invalid number value",
			json:    `1`,
			wantErr: true,
		},
		{
			name:    "invalid string value",
			json:    `"true"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got BoolOptional
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && (got.Present != tt.expected.Present || got.Value != tt.expected.Value) {
				t.Errorf("got %+v, want %+v", got, tt.expected)
			}
		})
	}
}

func TestTimeOptional(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected TimeOptional
		wantErr  bool
	}{
		{
			name:     "valid time",
			json:     `"2023-05-01T12:34:56Z"`,
			expected: TimeOptional{Value: time.Date(2023, 5, 1, 12, 34, 56, 0, time.UTC), Present: true},
		},
		{
			name:     "null value",
			json:     `null`,
			expected: TimeOptional{Present: false},
		},
		{
			name:    "invalid time format",
			json:    `"2023-05-01"`,
			wantErr: true,
		},
		{
			name:    "invalid type",
			json:    `123`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TimeOptional
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && (got.Present != tt.expected.Present || !got.Value.Equal(tt.expected.Value)) {
				t.Errorf("got %+v, want %+v", got, tt.expected)
			}
		})
	}
}

func TestComplexStructure(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected TestStructure
		wantErr  bool
	}{
		{
			name: "all fields with value",
			json: `{
                "text": "hello",
                "number": 42,
                "decimal": 42.5,
                "flag": true,
                "generic": [1,2,3]
            }`,
			expected: TestStructure{
				Text:    StringOptional{Value: "hello", Present: true},
				Number:  IntOptional{Value: 42, Present: true},
				Decimal: Float64Optional{Value: 42.5, Present: true},
				Flag:    BoolOptional{Value: true, Present: true},
				Generic: Optional[[]int]{Value: []int{1, 2, 3}, Present: true},
			},
		},
		{
			name: "all fields null",
			json: `{
                "text": null,
                "number": null,
                "decimal": null,
                "flag": null,
                "generic": null
            }`,
			expected: TestStructure{
				Text:    StringOptional{Present: false},
				Number:  IntOptional{Present: false},
				Decimal: Float64Optional{Present: false},
				Flag:    BoolOptional{Present: false},
				Generic: Optional[[]int]{Present: false},
			},
		},
		{
			name: "some fields null",
			json: `{
                "text": "hello",
                "number": null,
                "decimal": 42.5,
                "flag": null,
                "generic": [1,2,3]
            }`,
			expected: TestStructure{
				Text:    StringOptional{Value: "hello", Present: true},
				Number:  IntOptional{Present: false},
				Decimal: Float64Optional{Value: 42.5, Present: true},
				Flag:    BoolOptional{Present: false},
				Generic: Optional[[]int]{Value: []int{1, 2, 3}, Present: true},
			},
		},
		{
			name: "missing fields",
			json: `{
                "text": "hello",
                "decimal": 42.5
            }`,
			expected: TestStructure{
				Text:    StringOptional{Value: "hello", Present: true},
				Number:  IntOptional{Present: false},
				Decimal: Float64Optional{Value: 42.5, Present: true},
				Flag:    BoolOptional{Present: false},
				Generic: Optional[[]int]{Present: false},
			},
		},
		{
			name: "empty json",
			json: `{}`,
			expected: TestStructure{
				Text:    StringOptional{Present: false},
				Number:  IntOptional{Present: false},
				Decimal: Float64Optional{Present: false},
				Flag:    BoolOptional{Present: false},
				Generic: Optional[[]int]{Present: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got TestStructure
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if got.Text.Present != tt.expected.Text.Present || got.Text.Value != tt.expected.Text.Value {
					t.Errorf("Text: got %+v, want %+v", got.Text, tt.expected.Text)
				}
				if got.Number.Present != tt.expected.Number.Present || got.Number.Value != tt.expected.Number.Value {
					t.Errorf("Number: got %+v, want %+v", got.Number, tt.expected.Number)
				}
				if got.Decimal.Present != tt.expected.Decimal.Present || got.Decimal.Value != tt.expected.Decimal.Value {
					t.Errorf("Decimal: got %+v, want %+v", got.Decimal, tt.expected.Decimal)
				}
				if got.Flag.Present != tt.expected.Flag.Present || got.Flag.Value != tt.expected.Flag.Value {
					t.Errorf("Flag: got %+v, want %+v", got.Flag, tt.expected.Flag)
				}
			}
		})
	}
}

func TestMarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    TestStructure
		expected string
	}{
		{
			name: "all fields with value",
			input: TestStructure{
				Text:    StringOptional{Value: "hello", Present: true},
				Number:  IntOptional{Value: 42, Present: true},
				Decimal: Float64Optional{Value: 42.5, Present: true},
				Flag:    BoolOptional{Value: true, Present: true},
				Generic: Optional[[]int]{Value: []int{1, 2, 3}, Present: true},
			},
			expected: `{"text":"hello","number":42,"decimal":42.5,"flag":true,"generic":[1,2,3]}`,
		},
		{
			name: "all fields null",
			input: TestStructure{
				Text:    StringOptional{Present: false},
				Number:  IntOptional{Present: false},
				Decimal: Float64Optional{Present: false},
				Flag:    BoolOptional{Present: false},
				Generic: Optional[[]int]{Present: false},
			},
			expected: `{"text":null,"number":null,"decimal":null,"flag":null,"generic":null}`,
		},
		{
			name: "mix of values and null",
			input: TestStructure{
				Text:    StringOptional{Value: "hello", Present: true},
				Number:  IntOptional{Present: false},
				Decimal: Float64Optional{Value: 42.5, Present: true},
				Flag:    BoolOptional{Present: false},
				Generic: Optional[[]int]{Value: []int{1, 2, 3}, Present: true},
			},
			expected: `{"text":"hello","number":null,"decimal":42.5,"flag":null,"generic":[1,2,3]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.input)
			if err != nil {
				t.Errorf("Marshal() error = %v", err)
				return
			}

			if string(got) != tt.expected {
				t.Errorf("got %v, want %v", string(got), tt.expected)
			}

			// Verify that it can be deserialized back
			var roundTrip TestStructure
			err = json.Unmarshal(got, &roundTrip)
			if err != nil {
				t.Errorf("Unmarshal() error in roundtrip = %v", err)
			}
		})
	}
}
