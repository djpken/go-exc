package exc

// Config contains configuration for connecting to an exchange
type Config struct {
	// APIKey is the API key for authentication
	APIKey string

	// SecretKey is the secret key for authentication
	SecretKey string

	// Passphrase is the passphrase for authentication (if required)
	Passphrase string

	// Extra contains exchange-specific configuration
	Extra map[string]interface{}
}

// ExchangeType represents the type of exchange
type ExchangeType string

// Supported exchanges
const (
	OKX         ExchangeType = "OKX"
	OKXTest     ExchangeType = "OKX-TEST"
	Binance     ExchangeType = "BINANCE"
	BitMart     ExchangeType = "BIT-MART"
	BitMartTEST ExchangeType = "BIT-MART-Test"
	BingX       ExchangeType = "BING-X"
)
