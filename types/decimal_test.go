package types

import (
	"testing"
)

func TestDecimal_Add(t *testing.T) {
	tests := []struct {
		name     string
		d1       Decimal
		d2       Decimal
		expected string
		wantErr  bool
	}{
		{
			name:     "add positive numbers",
			d1:       Decimal("10.5"),
			d2:       Decimal("20.3"),
			expected: "30.8",
			wantErr:  false,
		},
		{
			name:     "add with zero",
			d1:       Decimal("100"),
			d2:       ZeroDecimal,
			expected: "100",
			wantErr:  false,
		},
		{
			name:     "add negative numbers",
			d1:       Decimal("-10.5"),
			d2:       Decimal("-5.5"),
			expected: "-16",
			wantErr:  false,
		},
		{
			name:     "add positive and negative",
			d1:       Decimal("100"),
			d2:       Decimal("-30"),
			expected: "70",
			wantErr:  false,
		},
		{
			name:     "invalid decimal",
			d1:       Decimal("invalid"),
			d2:       Decimal("10"),
			expected: "0",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.d1.Add(tt.d2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.expected {
				t.Errorf("Add() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecimal_Sub(t *testing.T) {
	tests := []struct {
		name     string
		d1       Decimal
		d2       Decimal
		expected string
		wantErr  bool
	}{
		{
			name:     "subtract positive numbers",
			d1:       Decimal("30.8"),
			d2:       Decimal("10.3"),
			expected: "20.5",
			wantErr:  false,
		},
		{
			name:     "subtract with zero",
			d1:       Decimal("100"),
			d2:       ZeroDecimal,
			expected: "100",
			wantErr:  false,
		},
		{
			name:     "subtract negative numbers",
			d1:       Decimal("-10"),
			d2:       Decimal("-5"),
			expected: "-5",
			wantErr:  false,
		},
		{
			name:     "result is negative",
			d1:       Decimal("10"),
			d2:       Decimal("30"),
			expected: "-20",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.d1.Sub(tt.d2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.expected {
				t.Errorf("Sub() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecimal_Mul(t *testing.T) {
	tests := []struct {
		name     string
		d1       Decimal
		d2       Decimal
		expected string
		wantErr  bool
	}{
		{
			name:     "multiply positive numbers",
			d1:       Decimal("10"),
			d2:       Decimal("5"),
			expected: "50",
			wantErr:  false,
		},
		{
			name:     "multiply with zero",
			d1:       Decimal("100"),
			d2:       ZeroDecimal,
			expected: "0",
			wantErr:  false,
		},
		{
			name:     "multiply decimals",
			d1:       Decimal("2.5"),
			d2:       Decimal("4"),
			expected: "10",
			wantErr:  false,
		},
		{
			name:     "multiply negative numbers",
			d1:       Decimal("-10"),
			d2:       Decimal("-5"),
			expected: "50",
			wantErr:  false,
		},
		{
			name:     "multiply positive and negative",
			d1:       Decimal("10"),
			d2:       Decimal("-5"),
			expected: "-50",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.d1.Mul(tt.d2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Mul() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.expected {
				t.Errorf("Mul() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecimal_Div(t *testing.T) {
	tests := []struct {
		name     string
		d1       Decimal
		d2       Decimal
		expected string
		wantErr  bool
	}{
		{
			name:     "divide positive numbers",
			d1:       Decimal("50"),
			d2:       Decimal("5"),
			expected: "10",
			wantErr:  false,
		},
		{
			name:     "divide by zero",
			d1:       Decimal("100"),
			d2:       ZeroDecimal,
			expected: "0",
			wantErr:  true,
		},
		{
			name:     "divide decimals",
			d1:       Decimal("10"),
			d2:       Decimal("4"),
			expected: "2.5",
			wantErr:  false,
		},
		{
			name:     "divide negative numbers",
			d1:       Decimal("-50"),
			d2:       Decimal("-5"),
			expected: "10",
			wantErr:  false,
		},
		{
			name:     "divide positive and negative",
			d1:       Decimal("50"),
			d2:       Decimal("-5"),
			expected: "-10",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.d1.Div(tt.d2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Div() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.expected {
				t.Errorf("Div() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecimal_Cmp(t *testing.T) {
	tests := []struct {
		name     string
		d1       Decimal
		d2       Decimal
		expected int
		wantErr  bool
	}{
		{
			name:     "less than",
			d1:       Decimal("10"),
			d2:       Decimal("20"),
			expected: -1,
			wantErr:  false,
		},
		{
			name:     "equal",
			d1:       Decimal("10"),
			d2:       Decimal("10"),
			expected: 0,
			wantErr:  false,
		},
		{
			name:     "greater than",
			d1:       Decimal("20"),
			d2:       Decimal("10"),
			expected: 1,
			wantErr:  false,
		},
		{
			name:     "decimal comparison",
			d1:       Decimal("10.5"),
			d2:       Decimal("10.49"),
			expected: 1,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.d1.Cmp(tt.d2)
			if (err != nil) != tt.wantErr {
				t.Errorf("Cmp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("Cmp() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecimal_LessThan(t *testing.T) {
	tests := []struct {
		name     string
		d1       Decimal
		d2       Decimal
		expected bool
	}{
		{
			name:     "true case",
			d1:       Decimal("10"),
			d2:       Decimal("20"),
			expected: true,
		},
		{
			name:     "false case - equal",
			d1:       Decimal("10"),
			d2:       Decimal("10"),
			expected: false,
		},
		{
			name:     "false case - greater",
			d1:       Decimal("20"),
			d2:       Decimal("10"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.d1.LessThan(tt.d2)
			if err != nil {
				t.Errorf("LessThan() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("LessThan() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecimal_GreaterThan(t *testing.T) {
	tests := []struct {
		name     string
		d1       Decimal
		d2       Decimal
		expected bool
	}{
		{
			name:     "true case",
			d1:       Decimal("20"),
			d2:       Decimal("10"),
			expected: true,
		},
		{
			name:     "false case - equal",
			d1:       Decimal("10"),
			d2:       Decimal("10"),
			expected: false,
		},
		{
			name:     "false case - less",
			d1:       Decimal("10"),
			d2:       Decimal("20"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.d1.GreaterThan(tt.d2)
			if err != nil {
				t.Errorf("GreaterThan() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("GreaterThan() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecimal_Equal(t *testing.T) {
	tests := []struct {
		name     string
		d1       Decimal
		d2       Decimal
		expected bool
	}{
		{
			name:     "true case",
			d1:       Decimal("10"),
			d2:       Decimal("10"),
			expected: true,
		},
		{
			name:     "false case - less",
			d1:       Decimal("10"),
			d2:       Decimal("20"),
			expected: false,
		},
		{
			name:     "false case - greater",
			d1:       Decimal("20"),
			d2:       Decimal("10"),
			expected: false,
		},
		{
			name:     "different formats same value",
			d1:       Decimal("10.0"),
			d2:       Decimal("10"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.d1.Equal(tt.d2)
			if err != nil {
				t.Errorf("Equal() error = %v", err)
				return
			}
			if result != tt.expected {
				t.Errorf("Equal() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecimal_Abs(t *testing.T) {
	tests := []struct {
		name     string
		d        Decimal
		expected string
		wantErr  bool
	}{
		{
			name:     "positive number",
			d:        Decimal("10.5"),
			expected: "10.5",
			wantErr:  false,
		},
		{
			name:     "negative number",
			d:        Decimal("-10.5"),
			expected: "10.5",
			wantErr:  false,
		},
		{
			name:     "zero",
			d:        ZeroDecimal,
			expected: "0",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.d.Abs()
			if (err != nil) != tt.wantErr {
				t.Errorf("Abs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.expected {
				t.Errorf("Abs() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecimal_Neg(t *testing.T) {
	tests := []struct {
		name     string
		d        Decimal
		expected string
		wantErr  bool
	}{
		{
			name:     "positive number",
			d:        Decimal("10.5"),
			expected: "-10.5",
			wantErr:  false,
		},
		{
			name:     "negative number",
			d:        Decimal("-10.5"),
			expected: "10.5",
			wantErr:  false,
		},
		{
			name:     "zero",
			d:        ZeroDecimal,
			expected: "0",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.d.Neg()
			if (err != nil) != tt.wantErr {
				t.Errorf("Neg() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result.String() != tt.expected {
				t.Errorf("Neg() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDecimal_ChainOperations(t *testing.T) {
	// Test chaining multiple operations
	d1 := Decimal("10")
	d2 := Decimal("5")
	d3 := Decimal("2")

	// (10 + 5) * 2 = 30
	result, err := d1.Add(d2)
	if err != nil {
		t.Fatalf("Add failed: %v", err)
	}

	result, err = result.Mul(d3)
	if err != nil {
		t.Fatalf("Mul failed: %v", err)
	}

	expected := "30"
	if result.String() != expected {
		t.Errorf("Chain operations result = %v, want %v", result, expected)
	}

	// (10 - 5) / 2 = 2.5
	result2, err := d1.Sub(d2)
	if err != nil {
		t.Fatalf("Sub failed: %v", err)
	}

	result2, err = result2.Div(d3)
	if err != nil {
		t.Fatalf("Div failed: %v", err)
	}

	expected2 := "2.5"
	if result2.String() != expected2 {
		t.Errorf("Chain operations result = %v, want %v", result2, expected2)
	}
}
