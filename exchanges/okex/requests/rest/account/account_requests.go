package account

import "github.com/djpken/go-exc/exchanges/okex/types"

type (
	GetBalance struct {
		Ccy []string `json:"ccy,omitempty"`
	}
	GetPositions struct {
		InstID   []string             `json:"instId,omitempty"`
		PosID    []string             `json:"posId,omitempty"`
		InstType types.InstrumentType `json:"instType,omitempty"`
	}
	GetAccountAndPositionRisk struct {
		InstType types.InstrumentType `json:"instType,omitempty"`
	}
	GetBills struct {
		Ccy      string               `json:"ccy,omitempty"`
		After    int64                `json:"after,omitempty,string"`
		Before   int64                `json:"before,omitempty,string"`
		Limit    int64                `json:"limit,omitempty,string"`
		InstType types.InstrumentType `json:"instType,omitempty"`
		MgnMode  types.MarginMode     `json:"mgnMode,omitempty"`
		CtType   types.ContractType   `json:"ctType,omitempty"`
		Type     types.BillType       `json:"type,omitempty,string"`
		SubType  types.BillSubType    `json:"subType,omitempty,string"`
	}
	SetPositionMode struct {
		PosMode types.PosModeType `json:"posMode"`
	}
	SetLeverage struct {
		Lever   int64              `json:"lever,string"`
		InstID  string             `json:"instId,omitempty"`
		Ccy     string             `json:"ccy,omitempty"`
		MgnMode types.MarginMode   `json:"mgnMode"`
		PosSide types.PositionSide `json:"posSide,omitempty"`
	}
	GetMaxBuySellAmount struct {
		Ccy    string          `json:"ccy,omitempty"`
		Px     float64         `json:"px,string,omitempty"`
		InstID []string        `json:"instId"`
		TdMode types.TradeMode `json:"tdMode"`
	}
	GetMaxAvailableTradeAmount struct {
		Ccy        string          `json:"ccy,omitempty"`
		InstID     string          `json:"instId"`
		ReduceOnly bool            `json:"reduceOnly,omitempty"`
		TdMode     types.TradeMode `json:"tdMode"`
	}
	IncreaseDecreaseMargin struct {
		InstID     string             `json:"instId"`
		Amt        float64            `json:"amt,string"`
		PosSide    types.PositionSide `json:"posSide"`
		ActionType types.CountAction  `json:"actionType"`
	}
	GetLeverage struct {
		InstID  []string         `json:"instId"`
		MgnMode types.MarginMode `json:"mgnMode"`
	}
	GetMaxLoan struct {
		InstID  string           `json:"instId"`
		MgnCcy  string           `json:"mgnCcy,omitempty"`
		MgnMode types.MarginMode `json:"mgnMode"`
	}
	GetFeeRates struct {
		InstID   string               `json:"instId,omitempty"`
		Uly      string               `json:"uly,omitempty"`
		Category types.FeeCategory    `json:"category,omitempty,string"`
		InstType types.InstrumentType `json:"instType"`
	}
	GetInterestAccrued struct {
		InstID  string           `json:"instId,omitempty"`
		Ccy     string           `json:"ccy,omitempty"`
		After   int64            `json:"after,omitempty,string"`
		Before  int64            `json:"before,omitempty,string"`
		Limit   int64            `json:"limit,omitempty,string"`
		MgnMode types.MarginMode `json:"mgnMode,omitempty"`
	}
	SetGreeks struct {
		GreeksType types.GreekType `json:"greeksType"`
	}
	SetAccountLevel struct {
		AcctLv string `json:"acctLv"`
	}
)
