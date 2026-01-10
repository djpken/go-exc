package exc

// Config contains configuration for connecting to an exchange
type Config struct {
	// APIKey is the API key for authentication
	APIKey string

	// SecretKey is the secret key for authentication
	SecretKey string

	// Passphrase is the passphrase for authentication (if required)
	Passphrase string

	// TestMode indicates whether to use test/demo environment
	TestMode bool

	// Extra contains exchange-specific configuration
	Extra map[string]interface{}
}

// ExchangeType represents the type of exchange
type ExchangeType string

// Supported exchanges
const (
	OKEx    ExchangeType = "okex"
	Binance ExchangeType = "binance"
	BitMart ExchangeType = "bitmart"
	BingX   ExchangeType = "bingx"
)
