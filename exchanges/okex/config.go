package okex

// OKExConfig contains OKEx-specific configuration
type OKExConfig struct {
	// APIKey is the API key
	APIKey string

	// SecretKey is the secret key
	SecretKey string

	// Passphrase is the passphrase
	Passphrase string

	// TestMode indicates whether to use demo environment
	TestMode bool

	// UseAWS indicates whether to use AWS servers
	UseAWS bool
}

// GetDestination returns the appropriate destination based on config
func (c *OKExConfig) GetDestination() Destination {
	if c.TestMode {
		return DemoServer
	}
	if c.UseAWS {
		return AwsServer
	}
	return NormalServer
}
