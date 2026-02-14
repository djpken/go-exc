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

// GetPositionV2 retrieves position details V2
//
// API: GET /contract/private/position-v2
// Documentation: https://developer-pro.bitmart.com/en/futures/#get-position-details-v2-signed
//
// Notes:
//   - If symbol is not provided, only positions with holdings are returned
//   - If symbol is provided, returns data even if no position exists (position fields show 0)
//
// Example:
//
//	// Get all positions
//	resp, err := client.Contract.GetPositionV2(contract.GetPositionV2Request{})
//
//	// Get specific symbol position
//	resp, err := client.Contract.GetPositionV2(contract.GetPositionV2Request{
//	    Symbol: "BTCUSDT",
//	})
func (c *Contract) GetPositionV2(req contract.GetPositionV2Request) (*responses.GetPositionV2Response, error) {
	endpoint := "/contract/private/position-v2"

	// Build query parameters
	params := ""
	if req.Symbol != "" {
		params += fmt.Sprintf("?symbol=%s", req.Symbol)
		if req.Account != "" {
			params += fmt.Sprintf("&account=%s", req.Account)
		}
	} else if req.Account != "" {
		params += fmt.Sprintf("?account=%s", req.Account)
	}

	endpoint += params

	var result responses.GetPositionV2Response
	if err := c.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SubmitLeverage sets leverage for a contract trading pair
//
// API: POST /contract/private/submit-leverage
// Documentation: https://developer-pro.bitmart.com/en/futures/#submit-leverage-signed
//
// Example:
//
//	// Set leverage for ETHUSDT to 5x with isolated margin
//	resp, err := client.Contract.SubmitLeverage(contract.SubmitLeverageRequest{
//	    Symbol: "ETHUSDT",
//	    Leverage: "5",
//	    OpenType: "isolated",
//	})
//
//	// Set leverage for BTCUSDT to 10x with cross margin
//	resp, err := client.Contract.SubmitLeverage(contract.SubmitLeverageRequest{
//	    Symbol: "BTCUSDT",
//	    Leverage: "10",
//	    OpenType: "cross",
//	})
func (c *Contract) SubmitLeverage(req contract.SubmitLeverageRequest) (*responses.SubmitLeverageResponse, error) {
	endpoint := "/contract/private/submit-leverage"

	var result responses.SubmitLeverageResponse
	if err := c.client.POST(endpoint, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetContractAssets queries contract assets detail (futures account balance)
//
// API: GET /contract/private/assets-detail
// Documentation: https://developer-pro.bitmart.com/en/futures/#get-contract-assets-keyed
//
// Returns detailed information about contract account assets including:
// - Available balance
// - Frozen balance
// - Position deposit (margin)
// - Total equity
// - Unrealized PnL
//
// Example:
//
//	// Get all contract assets
//	resp, err := client.Contract.GetContractAssets()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, asset := range resp.Data {
//	    fmt.Printf("%s: Available=%s, Equity=%s, Unrealized=%s\n",
//	        asset.Currency, asset.AvailableBalance, asset.Equity, asset.Unrealized)
//	}
func (c *Contract) GetContractAssets() (*responses.GetContractAssetsResponse, error) {
	endpoint := "/contract/private/assets-detail"

	var result responses.GetContractAssetsResponse
	if err := c.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetContractTrades queries contract trade details (order executions)
//
// API: GET /contract/private/trades
// Documentation: https://developer-pro.bitmart.com/en/futures/#get-order-trade-keyed
//
// Returns detailed trade execution information for contract orders including:
// - Execution price and volume
// - Maker/Taker type
// - Realized PnL
// - Trading fees
//
// Notes:
// - If no time range specified, queries last 7 days
// - Time range: end_time must be greater than start_time, max 90 days interval
// - Returns max 200 records per request
// - Supports order types: limit, market, liquidate, bankruptcy, adl, trailing
//
// Example:
//
//	// Get trades for a specific order
//	resp, err := client.Contract.GetContractTrades(contract.GetContractTradesRequest{
//	    OrderID: "220921197409432",
//	})
//
//	// Get trades for a symbol in time range
//	resp, err := client.Contract.GetContractTrades(contract.GetContractTradesRequest{
//	    Symbol:    "BTCUSDT",
//	    StartTime: 1662368173,
//	    EndTime:   1662368179,
//	})
func (c *Contract) GetContractTrades(req contract.GetContractTradesRequest) (*responses.GetContractTradesResponse, error) {
	endpoint := "/contract/private/trades"

	// Build query parameters
	params := make(map[string]string)
	if req.Symbol != "" {
		params["symbol"] = req.Symbol
	}
	if req.Account != "" {
		params["account"] = req.Account
	}
	if req.StartTime > 0 {
		params["start_time"] = fmt.Sprintf("%d", req.StartTime)
	}
	if req.EndTime > 0 {
		params["end_time"] = fmt.Sprintf("%d", req.EndTime)
	}
	if req.OrderID != "" {
		params["order_id"] = req.OrderID
	}
	if req.ClientOrderID != "" {
		params["client_order_id"] = req.ClientOrderID
	}

	// Add query parameters to endpoint
	if len(params) > 0 {
		endpoint += "?"
		first := true
		for key, value := range params {
			if !first {
				endpoint += "&"
			}
			endpoint += fmt.Sprintf("%s=%s", key, value)
			first = false
		}
	}

	var result responses.GetContractTradesResponse
	if err := c.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetKline retrieves kline/candlestick data for contract trading pairs
//
// API: GET /contract/public/kline
// Documentation: https://developer-pro.bitmart.com/en/futures/#get-k-line
//
// Notes:
// - Single request limited to 500 records
// - Symbol should be in format like BTCUSDT (no separator)
// - Step values: 1, 3, 5, 15, 30, 60, 120, 240, 360, 720, 1440, 4320, 10080 (in minutes)
// - start_time and end_time are Unix timestamps in seconds
//
// Example:
//
//	resp, err := client.Contract.GetKline(contract.GetContractKlineRequest{
//	    Symbol:    "BTCUSDT",
//	    Step:      5,
//	    StartTime: 1662518172,
//	    EndTime:   1662518172,
//	})
func (c *Contract) GetKline(req contract.GetContractKlineRequest) (*responses.GetContractKlineResponse, error) {
	endpoint := fmt.Sprintf("/contract/public/kline?symbol=%s&start_time=%d&end_time=%d",
		req.Symbol, req.StartTime, req.EndTime)

	if req.Step > 0 {
		endpoint += fmt.Sprintf("&step=%d", req.Step)
	}

	var result responses.GetContractKlineResponse
	if err := c.client.GET(endpoint, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
