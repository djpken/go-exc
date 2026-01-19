package rest

import (
	"fmt"

	"github.com/djpken/go-exc/exchanges/bitmart/requests/rest/funding"
	responses "github.com/djpken/go-exc/exchanges/bitmart/responses/funding"
)

// Funding provides access to BitMart funding API
type Funding struct {
	client *ClientRest
}

// NewFunding creates a new Funding API instance
func NewFunding(c *ClientRest) *Funding {
	return &Funding{client: c}
}

// GetDepositAddress retrieves deposit address for a currency
//
// API: GET /account/v1/deposit/address
func (f *Funding) GetDepositAddress(req funding.GetDepositAddressRequest) (*responses.DepositAddressResponse, error) {
	endpoint := fmt.Sprintf("/account/v1/deposit/address?currency=%s", req.Currency)

	if req.Chain != "" {
		endpoint += fmt.Sprintf("&chain=%s", req.Chain)
	}

	var result responses.DepositAddressResponse
	if err := f.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Withdraw initiates a withdrawal
//
// API: POST /account/v1/withdraw/apply
func (f *Funding) Withdraw(req funding.WithdrawRequest) (*responses.WithdrawResponse, error) {
	endpoint := "/account/v1/withdraw/apply"

	var result responses.WithdrawResponse
	if err := f.client.POST(endpoint, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetDepositHistory retrieves deposit history
//
// API: GET /account/v2/deposit-withdraw/history
func (f *Funding) GetDepositHistory(req funding.GetDepositHistoryRequest) (*responses.DepositHistoryResponse, error) {
	endpoint := "/account/v2/deposit-withdraw/history?operation_type=deposit"

	if req.Currency != "" {
		endpoint += fmt.Sprintf("&currency=%s", req.Currency)
	}

	var result responses.DepositHistoryResponse
	if err := f.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetWithdrawHistory retrieves withdrawal history
//
// API: GET /account/v2/deposit-withdraw/history
func (f *Funding) GetWithdrawHistory(req funding.GetWithdrawHistoryRequest) (*responses.WithdrawHistoryResponse, error) {
	endpoint := "/account/v2/deposit-withdraw/history?operation_type=withdraw"

	if req.Currency != "" {
		endpoint += fmt.Sprintf("&currency=%s", req.Currency)
	}

	var result responses.WithdrawHistoryResponse
	if err := f.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
