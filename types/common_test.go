package types

import "testing"

func TestDecimal_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		decimal  Decimal
		expected bool
	}{
		{"zero string", MustDecimal("0"), true},
		{"zero with decimal", MustDecimal("0.0"), true},
		{"zero with two decimals", MustDecimal("0.00"), true},
		{"zero with many decimals", MustDecimal("0.000000"), true},
		{"zero decimal constant", ZeroDecimal, true},
		{"positive integer", MustDecimal("1"), false},
		{"negative integer", MustDecimal("-1"), false},
		{"positive decimal", MustDecimal("0.1"), false},
		{"negative decimal", MustDecimal("-0.1"), false},
		{"large number", MustDecimal("1000.50"), false},
		{"scientific notation zero", MustDecimal("0e10"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.decimal.IsZero()
			if result != tt.expected {
				t.Errorf("Decimal(%q).IsZero() = %v, expected %v", tt.decimal, result, tt.expected)
			}
		})
	}
}

func TestDecimal_Float64(t *testing.T) {
	tests := []struct {
		name     string
		decimal  Decimal
		expected float64
	}{
		{"zero", MustDecimal("0"), 0},
		{"positive integer", MustDecimal("123"), 123},
		{"negative integer", MustDecimal("-456"), -456},
		{"positive decimal", MustDecimal("123.456"), 123.456},
		{"negative decimal", MustDecimal("-789.012"), -789.012},
		{"scientific notation", MustDecimal("1.23e10"), 1.23e10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.decimal.Float64()
			if err != nil {
				t.Errorf("Decimal(%q).Float64() unexpected error: %v", tt.decimal, err)
				return
			}
			if result != tt.expected {
				t.Errorf("Decimal(%q).Float64() = %v, expected %v", tt.decimal, result, tt.expected)
			}
		})
	}
}

func TestDecimal_String(t *testing.T) {
	tests := []struct {
		name     string
		decimal  Decimal
		expected string
	}{
		{"zero", MustDecimal("0"), "0"},
		{"integer", MustDecimal("123"), "123"},
		{"decimal", MustDecimal("123.456"), "123.456"},
		{"zero decimal", ZeroDecimal, "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.decimal.String()
			if result != tt.expected {
				t.Errorf("Decimal(%q).String() = %v, expected %v", tt.decimal, result, tt.expected)
			}
		})
	}
}

func TestNewDecimal_InvalidInput(t *testing.T) {
	_, err := NewDecimal("invalid")
	if err == nil {
		t.Error("NewDecimal(\"invalid\") expected error, got nil")
	}

	_, err = NewDecimal("")
	if err == nil {
		t.Error("NewDecimal(\"\") expected error, got nil")
	}
}
