package public

import (
	"encoding/json"
	"testing"
)

func TestFuturesTickerEventUnmarshal(t *testing.T) {
	// This is the actual JSON format from BitMart
	jsonData := `{
		"data": {
			"symbol": "BTCUSDT",
			"last_price": "97153.6",
			"volume_24": "25502894",
			"range": "0.0016599204475393",
			"mark_price": "97153.7",
			"index_price": "97185.614",
			"ask_price": "97153.9",
			"ask_vol": "28",
			"bid_price": "97153.4",
			"bid_vol": "428"
		},
		"group": "futures/ticker:BTCUSDT"
	}`

	var event FuturesTickerEvent
	err := json.Unmarshal([]byte(jsonData), &event)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify group field
	if event.Group != "futures/ticker:BTCUSDT" {
		t.Errorf("Expected group 'futures/ticker:BTCUSDT', got '%s'", event.Group)
	}

	// Verify data fields
	if event.Data.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got '%s'", event.Data.Symbol)
	}

	if event.Data.LastPrice != "97153.6" {
		t.Errorf("Expected last_price '97153.6', got '%s'", event.Data.LastPrice)
	}

	if event.Data.Volume24 != "25502894" {
		t.Errorf("Expected volume_24 '25502894', got '%s'", event.Data.Volume24)
	}

	if event.Data.Range != "0.0016599204475393" {
		t.Errorf("Expected range '0.0016599204475393', got '%s'", event.Data.Range)
	}

	if event.Data.MarkPrice != "97153.7" {
		t.Errorf("Expected mark_price '97153.7', got '%s'", event.Data.MarkPrice)
	}

	if event.Data.IndexPrice != "97185.614" {
		t.Errorf("Expected index_price '97185.614', got '%s'", event.Data.IndexPrice)
	}

	if event.Data.AskPrice != "97153.9" {
		t.Errorf("Expected ask_price '97153.9', got '%s'", event.Data.AskPrice)
	}

	if event.Data.AskVol != "28" {
		t.Errorf("Expected ask_vol '28', got '%s'", event.Data.AskVol)
	}

	if event.Data.BidPrice != "97153.4" {
		t.Errorf("Expected bid_price '97153.4', got '%s'", event.Data.BidPrice)
	}

	if event.Data.BidVol != "428" {
		t.Errorf("Expected bid_vol '428', got '%s'", event.Data.BidVol)
	}

	t.Logf("✓ Successfully parsed FuturesTickerEvent")
	t.Logf("  Symbol: %s", event.Data.Symbol)
	t.Logf("  Last Price: %s", event.Data.LastPrice)
	t.Logf("  Mark Price: %s", event.Data.MarkPrice)
	t.Logf("  Index Price: %s", event.Data.IndexPrice)
}
