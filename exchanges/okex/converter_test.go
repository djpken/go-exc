package okex

import (
	"testing"

	okexconstants "github.com/djpken/go-exc/exchanges/okex/constants"
	commontypes "github.com/djpken/go-exc/types"
)

func TestConverter_ConvertInstrumentType(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		input    commontypes.InstrumentType
		expected okexconstants.InstrumentType
	}{
		{
			name:     "any",
			input:    commontypes.InstrumentAny,
			expected: okexconstants.AnyInstrument,
		},
		{
			name:     "spot",
			input:    commontypes.InstrumentSpot,
			expected: okexconstants.SpotInstrument,
		},
		{
			name:     "margin",
			input:    commontypes.InstrumentMargin,
			expected: okexconstants.MarginInstrument,
		},
		{
			name:     "swap",
			input:    commontypes.InstrumentSwap,
			expected: okexconstants.SwapInstrument,
		},
		{
			name:     "futures",
			input:    commontypes.InstrumentFutures,
			expected: okexconstants.FuturesInstrument,
		},
		{
			name:     "option",
			input:    commontypes.InstrumentOption,
			expected: okexconstants.OptionsInstrument,
		},
		{
			name:     "default (empty)",
			input:    commontypes.InstrumentType(""),
			expected: okexconstants.AnyInstrument,
		},
		{
			name:     "default (unknown)",
			input:    commontypes.InstrumentType("unknown"),
			expected: okexconstants.AnyInstrument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertInstrumentType(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertInstrumentType(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConverter_ConvertOrderSide(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		input    string
		expected okexconstants.OrderSide
	}{
		{
			name:     "buy",
			input:    "buy",
			expected: okexconstants.OrderBuy,
		},
		{
			name:     "sell",
			input:    "sell",
			expected: okexconstants.OrderSell,
		},
		{
			name:     "default",
			input:    "unknown",
			expected: okexconstants.OrderBuy,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertOrderSide(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertOrderSide(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestConverter_ConvertOrderType(t *testing.T) {
	converter := NewConverter()

	tests := []struct {
		name     string
		input    string
		expected okexconstants.OrderType
	}{
		{
			name:     "market",
			input:    "market",
			expected: okexconstants.OrderMarket,
		},
		{
			name:     "limit",
			input:    "limit",
			expected: okexconstants.OrderLimit,
		},
		{
			name:     "post_only",
			input:    "post_only",
			expected: okexconstants.OrderPostOnly,
		},
		{
			name:     "fok",
			input:    "fok",
			expected: okexconstants.OrderFOK,
		},
		{
			name:     "ioc",
			input:    "ioc",
			expected: okexconstants.OrderIOC,
		},
		{
			name:     "default",
			input:    "unknown",
			expected: okexconstants.OrderLimit,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := converter.ConvertOrderType(tt.input)
			if result != tt.expected {
				t.Errorf("ConvertOrderType(%v) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}
