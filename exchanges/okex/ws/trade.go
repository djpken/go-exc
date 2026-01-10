package ws

import (
	requests "github.com/djpken/go-exc/exchanges/okex/requests/ws/trade"
	"github.com/djpken/go-exc/exchanges/okex/types"
	"github.com/djpken/go-exc/exchanges/okex/utils"
)

// Trade
//
// https://www.okex.com/docs-v5/en/#websocket-api-trade
type Trade struct {
	*ClientWs
}

// NewTrade returns a pointer to a fresh Trade
func NewTrade(c *ClientWs) *Trade {
	return &Trade{ClientWs: c}
}

// PlaceOrder
// You can place an order only if you have sufficient funds.
//
// https://www.okex.com/docs-v5/en/#websocket-api-trade-place-order
//
// Place orders in a batch. Maximum 20 orders can be placed at a time
//
// https://www.okex.com/docs-v5/en/#websocket-api-trade-place-multiple-orders
func (c *Trade) PlaceOrder(req ...requests.PlaceOrder) error {
	tmpArgs := make([]map[string]string, len(req))
	op := types.OrderOperation
	if len(req) > 1 {
		op = types.BatchOrderOperation
	}
	for i, order := range req {
		tmpArgs[i] = utils.S2M(order)
	}
	return c.Send(true, op, tmpArgs, map[string]string{"id": req[0].ID})
}

// CancelOrder
// Cancel an incomplete order
//
// https://www.okex.com/docs-v5/en/#websocket-api-trade-place-order
//
// Cancel incomplete orders in batches. Maximum 20 orders can be canceled at a time.
//
// https://www.okex.com/docs-v5/en/#websocket-api-trade-cancel-multiple-orders
func (c *Trade) CancelOrder(req ...requests.CancelOrder) error {
	tmpArgs := make([]map[string]string, len(req))
	op := types.CancelOrderOperation
	if len(req) > 1 {
		op = types.BatchCancelOrderOperation
	}
	for i, order := range req {
		tmpArgs[i] = utils.S2M(order)
	}
	return c.Send(true, op, tmpArgs, map[string]string{"id": req[0].ID})
}

// AmendOrder
// Amend an incomplete order.
//
// https://www.okex.com/docs-v5/en/#websocket-api-trade-place-order
//
// Amend incomplete orders in batches. Maximum 20 orders can be amended at a time.
//
// https://www.okex.com/docs-v5/en/#websocket-api-trade-amend-multiple-orders
func (c *Trade) AmendOrder(req ...requests.AmendOrder) error {
	tmpArgs := make([]map[string]string, len(req))
	op := types.AmendOrderOperation
	if len(req) > 1 {
		op = types.BatchAmendOrderOperation
	}
	for i, order := range req {
		tmpArgs[i] = utils.S2M(order)
	}
	return c.Send(true, op, tmpArgs, map[string]string{"id": req[0].ID})
}
