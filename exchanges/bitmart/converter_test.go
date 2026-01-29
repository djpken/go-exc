package bitmart

import (
	"testing"

	accountmodels "github.com/djpken/go-exc/exchanges/bitmart/models/account"
)

func TestConverter_ConvertAccountBalance(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name              string
		balances          []accountmodels.Balance
		expectedEquity    string
		expectedBalances  int
	}{
		{
			name: "with USD valuation",
			balances: []accountmodels.Balance{
				{
					Currency:              "BTC",
					Name:                  "Bitcoin",
					Available:             "1.5",
					AvailableUsdValuation: "60000.00",
					Frozen:                "0.0",
					UnAvailable:           "0.0",
				},
				{
					Currency:              "ETH",
					Name:                  "Ethereum",
					Available:             "10.0",
					AvailableUsdValuation: "30000.00",
					Frozen:                "0.0",
					UnAvailable:           "0.0",
				},
				{
					Currency:              "USDT",
					Name:                  "Tether USD",
					Available:             "5000.0",
					AvailableUsdValuation: "5000.00",
					Frozen:                "100.0",
					UnAvailable:           "100.0",
				},
			},
			expectedEquity:   "95000", // 60000 + 30000 + 5000
			expectedBalances: 3,
		},
		{
			name: "without USD valuation",
			balances: []accountmodels.Balance{
				{
					Currency:    "BTC",
					Name:        "Bitcoin",
					Available:   "1.0",
					Frozen:      "0.0",
					UnAvailable: "0.0",
				},
			},
			expectedEquity:   "0",
			expectedBalances: 1,
		},
		{
			name:             "empty balances",
			balances:         []accountmodels.Balance{},
			expectedEquity:   "0",
			expectedBalances: 0,
		},
		{
			name:             "nil balances",
			balances:         nil,
			expectedEquity:   "",
			expectedBalances: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertAccountBalance(tt.balances)

			// Check nil case
			if tt.balances == nil {
				if result != nil {
					t.Errorf("Expected nil result for nil balances")
				}
				return
			}

			// Check balances count
			if len(result.Balances) != tt.expectedBalances {
				t.Errorf("Expected %d balances, got %d", tt.expectedBalances, len(result.Balances))
			}

			// Check total equity
			equity := result.TotalEquity.String()
			if equity != tt.expectedEquity {
				t.Errorf("Expected total equity %s, got %s", tt.expectedEquity, equity)
			}

			// Verify balance details for non-empty cases
			if len(tt.balances) > 0 && len(result.Balances) > 0 {
				firstBal := result.Balances[0]
				if firstBal.Currency != tt.balances[0].Currency {
					t.Errorf("Expected currency %s, got %s", tt.balances[0].Currency, firstBal.Currency)
				}

				// Check Extra fields
				if firstBal.Extra != nil {
					if name, ok := firstBal.Extra["name"].(string); ok {
						if name != tt.balances[0].Name {
							t.Errorf("Expected name %s in Extra, got %s", tt.balances[0].Name, name)
						}
					}

					if usdVal, ok := firstBal.Extra["available_usd_valuation"].(string); ok {
						if usdVal != tt.balances[0].AvailableUsdValuation {
							t.Errorf("Expected USD valuation %s in Extra, got %s", tt.balances[0].AvailableUsdValuation, usdVal)
						}
					}
				}
			}
		})
	}
}

func TestConverter_ConvertAccountBalance_TotalCalculation(t *testing.T) {
	converter := NewConverter()

	balances := []accountmodels.Balance{
		{
			Currency:    "BTC",
			Name:        "Bitcoin",
			Available:   "1.5",
			Frozen:      "0.2",
			UnAvailable: "0.5", // Available + UnAvailable = Total (should be calculated)
		},
	}

	result := converter.ConvertAccountBalance(balances)

	if len(result.Balances) != 1 {
		t.Fatalf("Expected 1 balance, got %d", len(result.Balances))
	}

	bal := result.Balances[0]

	// Total should be Available (1.5) + UnAvailable (0.5) = 2.0
	totalFloat, err := bal.Total.Float64()
	if err != nil {
		t.Fatalf("Failed to parse total: %v", err)
	}

	expectedTotal := 2.0
	if totalFloat != expectedTotal {
		t.Errorf("Expected total %f, got %f", expectedTotal, totalFloat)
	}
}

func TestConverter_ConvertBalance(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name          string
		balance       *accountmodels.Balance
		expectedTotal string
	}{
		{
			name: "with all fields",
			balance: &accountmodels.Balance{
				Currency:              "BTC",
				Name:                  "Bitcoin",
				Available:             "1.0",
				AvailableUsdValuation: "50000.00",
				Frozen:                "0.1",
				UnAvailable:           "0.2", // Should calculate total as 1.0 + 0.2 = 1.2
			},
			expectedTotal: "1.2",
		},
		{
			name: "calculate total from available and unavailable",
			balance: &accountmodels.Balance{
				Currency:    "ETH",
				Name:        "Ethereum",
				Available:   "10.0",
				Frozen:      "1.0",
				UnAvailable: "2.0", // Should be calculated as 10.0 + 2.0 = 12.0
			},
			expectedTotal: "12",
		},
		{
			name:          "nil balance",
			balance:       nil,
			expectedTotal: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertBalance(tt.balance)

			if tt.balance == nil {
				if result != nil {
					t.Errorf("Expected nil result for nil balance")
				}
				return
			}

			if result.Total.String() != tt.expectedTotal {
				t.Errorf("Expected total %s, got %s", tt.expectedTotal, result.Total.String())
			}

			// Verify Extra fields
			if result.Extra != nil {
				if name, ok := result.Extra["name"].(string); ok {
					if name != tt.balance.Name {
						t.Errorf("Expected name %s in Extra, got %s", tt.balance.Name, name)
					}
				}
			}
		})
	}
}
