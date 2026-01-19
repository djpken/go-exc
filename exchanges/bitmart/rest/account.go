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

// GetBalance retrieves account balance
//
// API: GET /spot/v1/wallet
func (a *Account) GetBalance(req account.GetBalanceRequest) (*responses.BalanceResponse, error) {
	endpoint := "/spot/v1/wallet"

	if req.Currency != "" {
		endpoint += fmt.Sprintf("?currency=%s", req.Currency)
	}

	var result responses.BalanceResponse
	if err := a.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetWalletBalance retrieves wallet balance by wallet type
//
// API: GET /account/v1/wallet
func (a *Account) GetWalletBalance(req account.GetWalletBalanceRequest) (*responses.WalletBalanceResponse, error) {
	endpoint := "/account/v1/wallet"

	if req.WalletType != "" {
		endpoint += fmt.Sprintf("?wallet_type=%s", req.WalletType)
	}

	var result responses.WalletBalanceResponse
	if err := a.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
