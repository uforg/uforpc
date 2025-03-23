package pieces

import (
	"encoding/json"
	"testing"
)

// TestStructure includes all nullable types
type TestStructure struct {
	Text    NullString  `json:"text,omitempty"`
	Number  NullInt     `json:"number,omitempty"`
	Decimal NullFloat64 `json:"decimal,omitempty"`
	Flag    NullBool    `json:"flag,omitempty"`
	Generic Null[[]int] `json:"generic,omitempty"`
}

func TestNullString(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected NullString
		wantErr  bool
	}{
		{
			name:     "valor normal",
			json:     `"hello"`,
			expected: NullString{Value: "hello", Valid: true},
		},
		{
			name:     "valor vacío",
			json:     `""`,
			expected: NullString{Value: "", Valid: true},
		},
		{
			name:     "valor null",
			json:     `null`,
			expected: NullString{Valid: false},
		},
		{
			name:    "valor inválido",
			json:    `123`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got NullString
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && (got.Valid != tt.expected.Valid || got.Value != tt.expected.Value) {
				t.Errorf("got %+v, want %+v", got, tt.expected)
			}
		})
	}
}

func TestNullInt(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected NullInt
		wantErr  bool
	}{
		{
			name:     "número positivo",
			json:     `42`,
			expected: NullInt{Value: 42, Valid: true},
		},
		{
			name:     "número negativo",
			json:     `-42`,
			expected: NullInt{Value: -42, Valid: true},
		},
		{
			name:     "cero",
			json:     `0`,
			expected: NullInt{Value: 0, Valid: true},
		},
		{
			name:     "valor null",
			json:     `null`,
			expected: NullInt{Valid: false},
		},
		{
			name:    "valor inválido string",
			json:    `"123"`,
			wantErr: true,
		},
		{
			name:    "valor inválido float",
			json:    `123.45`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got NullInt
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && (got.Valid != tt.expected.Valid || got.Value != tt.expected.Value) {
				t.Errorf("got %+v, want %+v", got, tt.expected)
			}
		})
	}
}

func TestNullFloat64(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected NullFloat64
		wantErr  bool
	}{
		{
			name:     "número entero",
			json:     `42`,
			expected: NullFloat64{Value: 42.0, Valid: true},
		},
		{
			name:     "número decimal",
			json:     `42.5`,
			expected: NullFloat64{Value: 42.5, Valid: true},
		},
		{
			name:     "número negativo",
			json:     `-42.5`,
			expected: NullFloat64{Value: -42.5, Valid: true},
		},
		{
			name:     "cero",
			json:     `0`,
			expected: NullFloat64{Value: 0, Valid: true},
		},
		{
			name:     "valor null",
			json:     `null`,
			expected: NullFloat64{Valid: false},
		},
		{
			name:    "valor inválido",
			json:    `"123.45"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got NullFloat64
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && (got.Valid != tt.expected.Valid || got.Value != tt.expected.Value) {
				t.Errorf("got %+v, want %+v", got, tt.expected)
			}
		})
	}
}

func TestNullBool(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected NullBool
		wantErr  bool
	}{
		{
			name:     "true",
			json:     `true`,
			expected: NullBool{Value: true, Valid: true},
		},
		{
			name:     "false",
			json:     `false`,
			expected: NullBool{Value: false, Valid: true},
		},
		{
			name:     "valor null",
			json:     `null`,
			expected: NullBool{Valid: false},
		},
		{
			name:    "valor inválido número",
			json:    `1`,
			wantErr: true,
		},
		{
			name:    "valor inválido string",
			json:    `"true"`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got NullBool
			err := json.Unmarshal([]byte(tt.json), &got)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil && (got.Valid != tt.expected.Valid || got.Value != tt.expected.Value) {
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
			name: "todos los campos con valor",
			json: `{
                "text": "hello",
                "number": 42,
                "decimal": 42.5,
                "flag": true,
                "generic": [1,2,3]
            }`,
			expected: TestStructure{
				Text:    NullString{Value: "hello", Valid: true},
				Number:  NullInt{Value: 42, Valid: true},
				Decimal: NullFloat64{Value: 42.5, Valid: true},
				Flag:    NullBool{Value: true, Valid: true},
				Generic: Null[[]int]{Value: []int{1, 2, 3}, Valid: true},
			},
		},
		{
			name: "todos los campos null",
			json: `{
                "text": null,
                "number": null,
                "decimal": null,
                "flag": null,
                "generic": null
            }`,
			expected: TestStructure{
				Text:    NullString{Valid: false},
				Number:  NullInt{Valid: false},
				Decimal: NullFloat64{Valid: false},
				Flag:    NullBool{Valid: false},
				Generic: Null[[]int]{Valid: false},
			},
		},
		{
			name: "algunos campos null",
			json: `{
                "text": "hello",
                "number": null,
                "decimal": 42.5,
                "flag": null,
                "generic": [1,2,3]
            }`,
			expected: TestStructure{
				Text:    NullString{Value: "hello", Valid: true},
				Number:  NullInt{Valid: false},
				Decimal: NullFloat64{Value: 42.5, Valid: true},
				Flag:    NullBool{Valid: false},
				Generic: Null[[]int]{Value: []int{1, 2, 3}, Valid: true},
			},
		},
		{
			name: "campos faltantes",
			json: `{
                "text": "hello",
                "decimal": 42.5
            }`,
			expected: TestStructure{
				Text:    NullString{Value: "hello", Valid: true},
				Number:  NullInt{Valid: false},
				Decimal: NullFloat64{Value: 42.5, Valid: true},
				Flag:    NullBool{Valid: false},
				Generic: Null[[]int]{Valid: false},
			},
		},
		{
			name: "json vacío",
			json: `{}`,
			expected: TestStructure{
				Text:    NullString{Valid: false},
				Number:  NullInt{Valid: false},
				Decimal: NullFloat64{Valid: false},
				Flag:    NullBool{Valid: false},
				Generic: Null[[]int]{Valid: false},
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
				if got.Text.Valid != tt.expected.Text.Valid || got.Text.Value != tt.expected.Text.Value {
					t.Errorf("Text: got %+v, want %+v", got.Text, tt.expected.Text)
				}
				if got.Number.Valid != tt.expected.Number.Valid || got.Number.Value != tt.expected.Number.Value {
					t.Errorf("Number: got %+v, want %+v", got.Number, tt.expected.Number)
				}
				if got.Decimal.Valid != tt.expected.Decimal.Valid || got.Decimal.Value != tt.expected.Decimal.Value {
					t.Errorf("Decimal: got %+v, want %+v", got.Decimal, tt.expected.Decimal)
				}
				if got.Flag.Valid != tt.expected.Flag.Valid || got.Flag.Value != tt.expected.Flag.Value {
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
			name: "todos los campos con valor",
			input: TestStructure{
				Text:    NullString{Value: "hello", Valid: true},
				Number:  NullInt{Value: 42, Valid: true},
				Decimal: NullFloat64{Value: 42.5, Valid: true},
				Flag:    NullBool{Value: true, Valid: true},
				Generic: Null[[]int]{Value: []int{1, 2, 3}, Valid: true},
			},
			expected: `{"text":"hello","number":42,"decimal":42.5,"flag":true,"generic":[1,2,3]}`,
		},
		{
			name: "todos los campos null",
			input: TestStructure{
				Text:    NullString{Valid: false},
				Number:  NullInt{Valid: false},
				Decimal: NullFloat64{Valid: false},
				Flag:    NullBool{Valid: false},
				Generic: Null[[]int]{Valid: false},
			},
			expected: `{"text":null,"number":null,"decimal":null,"flag":null,"generic":null}`,
		},
		{
			name: "mezcla de valores y null",
			input: TestStructure{
				Text:    NullString{Value: "hello", Valid: true},
				Number:  NullInt{Valid: false},
				Decimal: NullFloat64{Value: 42.5, Valid: true},
				Flag:    NullBool{Valid: false},
				Generic: Null[[]int]{Value: []int{1, 2, 3}, Valid: true},
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

			// Verificar que se puede deserializar de vuelta
			var roundTrip TestStructure
			err = json.Unmarshal(got, &roundTrip)
			if err != nil {
				t.Errorf("Unmarshal() error en roundtrip = %v", err)
			}
		})
	}
}
