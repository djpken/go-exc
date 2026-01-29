package rest

import (
	"fmt"

	"github.com/djpken/go-exc/exchanges/bitmart/requests/rest/contract"
	responses "github.com/djpken/go-exc/exchanges/bitmart/responses/contract"
)

// Contract provides access to BitMart contract API
type Contract struct {
	client *ClientRest
}

// NewContract creates a new Contract API instance
func NewContract(c *ClientRest) *Contract {
	return &Contract{client: c}
}

// GetContractDetails retrieves contract details for perpetual contracts
//
// API: GET /contract/public/details
// Documentation: https://developer-pro.bitmart.com/en/futures/#get-contract-details
func (c *Contract) GetContractDetails(req contract.GetContractDetailsRequest) (*responses.ContractDetailsResponse, error) {
	endpoint := "/contract/public/details"

	// Add symbol parameter if specified
	if req.Symbol != "" {
		endpoint += fmt.Sprintf("?symbol=%s", req.Symbol)
	}

	var result responses.ContractDetailsResponse
	if err := c.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SubmitOrder places a new contract order
//
// API: POST /contract/private/submit-order
// Documentation: https://developer-pro.bitmart.com/en/futures/#place-order-signed
//
// Example:
//
//	resp, err := client.Contract.SubmitOrder(contract.SubmitOrderRequest{
//	    Symbol: "BTCUSDT",
//	    Side: 1,           // 1=Open long, 4=Open short
//	    Type: "limit",
//	    Size: 10,          // Number of contracts
//	    Price: "40000",
//	    Leverage: "10",
//	    OpenType: "isolated",
//	})
func (c *Contract) SubmitOrder(req contract.SubmitOrderRequest) (*responses.SubmitOrderResponse, error) {
	endpoint := "/contract/private/submit-order"

	var result responses.SubmitOrderResponse
	if err := c.client.POST(endpoint, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
