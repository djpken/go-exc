package rest

import (
	"fmt"

	"github.com/djpken/go-exc/exchanges/bitmart/requests/rest/account"
	responses "github.com/djpken/go-exc/exchanges/bitmart/responses/account"
)

// Account provides access to BitMart account API
type Account struct {
	client *ClientRest
}

// NewAccount creates a new Account API instance
func NewAccount(c *ClientRest) *Account {
	return &Account{client: c}
}

// GetWalletBalance retrieves wallet balance
//
// API: GET /account/v1/wallet
// Parameters:
//   - currency: Optional currency code (e.g., "BTC")
//   - needUsdValuation: Optional flag to return USD valuation (default: false)
func (a *Account) GetWalletBalance(req account.GetWalletBalanceRequest) (*responses.WalletBalanceResponse, error) {
	endpoint := "/account/v1/wallet"

	// Build query parameters
	params := []string{}
	if req.Currency != "" {
		params = append(params, fmt.Sprintf("currency=%s", req.Currency))
	}
	if req.NeedUsdValuation {
		params = append(params, "needUsdValuation=true")
	}

	// Add query parameters to endpoint
	if len(params) > 0 {
		endpoint += "?"
		for i, param := range params {
			if i > 0 {
				endpoint += "&"
			}
			endpoint += param
		}
	}

	var result responses.WalletBalanceResponse
	if err := a.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
