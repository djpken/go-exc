package types

import "testing"

func TestDecimal_IsZero(t *testing.T) {
	tests := []struct {
		name     string
		decimal  Decimal
		expected bool
	}{
		{"zero string", "0", true},
		{"zero with decimal", "0.0", true},
		{"zero with two decimals", "0.00", true},
		{"zero with many decimals", "0.000000", true},
		{"empty string", "", true},
		{"zero decimal constant", ZeroDecimal, true},
		{"positive integer", "1", false},
		{"negative integer", "-1", false},
		{"positive decimal", "0.1", false},
		{"negative decimal", "-0.1", false},
		{"large number", "1000.50", false},
		{"scientific notation zero", "0e10", true},
		{"invalid string", "abc", false},
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
		hasError bool
	}{
		{"zero", "0", 0, false},
		{"positive integer", "123", 123, false},
		{"negative integer", "-456", -456, false},
		{"positive decimal", "123.456", 123.456, false},
		{"negative decimal", "-789.012", -789.012, false},
		{"scientific notation", "1.23e10", 1.23e10, false},
		{"invalid", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.decimal.Float64()
			if tt.hasError {
				if err == nil {
					t.Errorf("Decimal(%q).Float64() expected error, got nil", tt.decimal)
				}
			} else {
				if err != nil {
					t.Errorf("Decimal(%q).Float64() unexpected error: %v", tt.decimal, err)
				}
				if result != tt.expected {
					t.Errorf("Decimal(%q).Float64() = %v, expected %v", tt.decimal, result, tt.expected)
				}
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
		{"zero", "0", "0"},
		{"integer", "123", "123"},
		{"decimal", "123.456", "123.456"},
		{"empty", "", ""},
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
