package private

import (
	"github.com/djpken/go-exc/exchanges/okex/events"
	"github.com/djpken/go-exc/exchanges/okex/models/account"
	"github.com/djpken/go-exc/exchanges/okex/models/trade"
	"github.com/djpken/go-exc/exchanges/okex/constants"
)

type (
	Account struct {
		Arg       *events.Argument   `json:"arg"`
		EventType constants.EventType    `json:"eventType"`
		CurPage   int                `json:"curPage"`
		LastPage  bool               `json:"lastPage"`
		Balances  []*account.Balance `json:"data"`
	}
	Position struct {
		Arg       *events.Argument    `json:"arg"`
		EventType constants.EventType     `json:"eventType"`
		CurPage   int                 `json:"curPage"`
		LastPage  bool                `json:"lastPage"`
		Positions []*account.Position `json:"data"`
	}
	BalanceAndPosition struct {
		Arg                 *events.Argument              `json:"arg"`
		BalanceAndPositions []*account.BalanceAndPosition `json:"data"`
	}
	Order struct {
		Arg    *events.Argument `json:"arg"`
		Orders []*trade.Order   `json:"data"`
	}
)
