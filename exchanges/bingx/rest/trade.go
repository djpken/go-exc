package rest

import "fmt"

// Trade provides BingX trading endpoints
type Trade struct {
	client *ClientRest
}

func NewTrade(c *ClientRest) *Trade { return &Trade{client: c} }

// OrderData holds details for a placed or queried order
type OrderData struct {
	Symbol        string `json:"symbol"`
	OrderID       int64  `json:"orderId"`
	Side          string `json:"side"`
	PositionSide  string `json:"positionSide"`
	Type          string `json:"type"`
	ClientOrderID string `json:"clientOrderID"`
	Price         string `json:"price"`
	OrigQty       string `json:"origQty"`
	ExecutedQty   string `json:"executedQty"`
	AvgPrice      string `json:"avgPrice"`
	CumQuote      string `json:"cumQuote"`
	StopPrice     string `json:"stopPrice"`
	Status        string `json:"status"`
	Profit        string `json:"profit"`
	Commission    string `json:"commission"`
	Time          int64  `json:"time"`
	UpdateTime    int64  `json:"updateTime"`
	WorkingType   string `json:"workingType"`
}

// PlaceOrderResponseData is the data returned when placing an order
type PlaceOrderResponseData struct {
	Order OrderData `json:"order"`
}

// PlaceOrderResponse is the full API response for placing an order
type PlaceOrderResponse struct {
	Code int                    `json:"code"`
	Data PlaceOrderResponseData `json:"data"`
}

// PlaceOrder places a new perpetual swap order
// POST /openApi/swap/v2/trade/order
func (t *Trade) PlaceOrder(
	symbol, side, positionSide, orderType string,
	price, quantity float64,
	clientOrderID string,
	extra map[string]string,
) (*PlaceOrderResponse, error) {
	params := map[string]string{
		"symbol": symbol,
		"side":   side,
		"type":   orderType,
	}
	if positionSide != "" {
		params["positionSide"] = positionSide
	}
	if price > 0 {
		params["price"] = fmt.Sprintf("%f", price)
	}
	if quantity > 0 {
		params["quantity"] = fmt.Sprintf("%f", quantity)
	}
	if clientOrderID != "" {
		params["clientOrderID"] = clientOrderID
	}
	for k, v := range extra {
		params[k] = v
	}

	var result PlaceOrderResponse
	if err := t.client.POST("/openApi/swap/v2/trade/order", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// QueryOrderResponse is the full API response for querying an order
type QueryOrderResponse struct {
	Code int       `json:"code"`
	Data OrderData `json:"data"`
}

// GetOrder queries a single order by orderID or clientOrderID
// GET /openApi/swap/v2/trade/order
func (t *Trade) GetOrder(symbol string, orderID int64, clientOrderID string) (*QueryOrderResponse, error) {
	params := map[string]string{"symbol": symbol}
	if orderID > 0 {
		params["orderId"] = fmt.Sprintf("%d", orderID)
	}
	if clientOrderID != "" {
		params["clientOrderID"] = clientOrderID
	}

	var result QueryOrderResponse
	if err := t.client.GET("/openApi/swap/v2/trade/order", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CancelOrderResponse is the full API response for canceling an order
type CancelOrderResponse struct {
	Code int       `json:"code"`
	Data OrderData `json:"data"`
}

// CancelOrder cancels an open order
// DELETE /openApi/swap/v2/trade/order
func (t *Trade) CancelOrder(symbol string, orderID int64, clientOrderID string) (*CancelOrderResponse, error) {
	params := map[string]string{"symbol": symbol}
	if orderID > 0 {
		params["orderId"] = fmt.Sprintf("%d", orderID)
	}
	if clientOrderID != "" {
		params["clientOrderID"] = clientOrderID
	}

	var result CancelOrderResponse
	if err := t.client.DELETE("/openApi/swap/v2/trade/order", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
