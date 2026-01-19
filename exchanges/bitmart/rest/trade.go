package rest

import (
	"github.com/djpken/go-exc/exchanges/bitmart/requests/rest/trade"
	responses "github.com/djpken/go-exc/exchanges/bitmart/responses/trade"
)

// Trade provides access to BitMart trading API
type Trade struct {
	client *ClientRest
}

// NewTrade creates a new Trade API instance
func NewTrade(c *ClientRest) *Trade {
	return &Trade{client: c}
}

// PlaceOrder places a new order
//
// API: POST /spot/v2/submit_order
func (t *Trade) PlaceOrder(req trade.PlaceOrderRequest) (*responses.PlaceOrderResponse, error) {
	endpoint := "/spot/v2/submit_order"

	var result responses.PlaceOrderResponse
	if err := t.client.POST(endpoint, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CancelOrder cancels an existing order
//
// API: POST /spot/v3/cancel_order
func (t *Trade) CancelOrder(req trade.CancelOrderRequest) (*responses.CancelOrderResponse, error) {
	endpoint := "/spot/v3/cancel_order"

	var result responses.CancelOrderResponse
	if err := t.client.POST(endpoint, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CancelAllOrders cancels all orders for a symbol or all symbols
//
// API: POST /spot/v1/cancel_orders
func (t *Trade) CancelAllOrders(req trade.CancelAllOrdersRequest) (*responses.CancelAllOrdersResponse, error) {
	endpoint := "/spot/v1/cancel_orders"

	var result responses.CancelAllOrdersResponse
	if err := t.client.POST(endpoint, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOrder retrieves order details
//
// API: GET /spot/v2/order_detail
func (t *Trade) GetOrder(req trade.GetOrderRequest) (*responses.OrderResponse, error) {
	endpoint := "/spot/v2/order_detail"

	var result responses.OrderResponse
	if err := t.client.POST(endpoint, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOrders retrieves order list
//
// API: POST /spot/v2/orders
func (t *Trade) GetOrders(req trade.GetOrdersRequest) (*responses.OrdersResponse, error) {
	endpoint := "/spot/v2/orders"

	var result responses.OrdersResponse
	if err := t.client.POST(endpoint, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTrades retrieves trade history
//
// API: POST /spot/v2/trades
func (t *Trade) GetTrades(req trade.GetTradesRequest) (*responses.TradesResponse, error) {
	endpoint := "/spot/v2/trades"

	var result responses.TradesResponse
	if err := t.client.POST(endpoint, req, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
