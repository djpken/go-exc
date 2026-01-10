package trade

import (
	"github.com/djpken/go-exc/exchanges/okex/types"
)

type (
	PlaceOrder struct {
		ClOrdID string            `json:"clOrdId"`
		Tag     string            `json:"tag"`
		SMsg    string            `json:"sMsg"`
		SCode   types.JSONInt64   `json:"sCode"`
		OrdID   types.JSONFloat64 `json:"ordId"`
	}
	CancelOrder struct {
		OrdID   string            `json:"ordId"`
		ClOrdID string            `json:"clOrdId"`
		SMsg    string            `json:"sMsg"`
		SCode   types.JSONFloat64 `json:"sCode"`
	}
	AmendOrder struct {
		OrdID   string            `json:"ordId"`
		ClOrdID string            `json:"clOrdId"`
		ReqID   string            `json:"reqId"`
		SMsg    string            `json:"sMsg"`
		SCode   types.JSONFloat64 `json:"sCode"`
	}
	ClosePosition struct {
		InstID  string             `json:"instId"`
		PosSide types.PositionSide `json:"posSide"`
		Tag     string             `json:"tag"`
	}
	Order struct {
		InstID      string               `json:"instId"`
		Ccy         string               `json:"ccy"`
		OrdID       string               `json:"ordId"`
		ClOrdID     string               `json:"clOrdId"`
		TradeID     string               `json:"tradeId"`
		Tag         string               `json:"tag"`
		Category    string               `json:"category"`
		FeeCcy      string               `json:"feeCcy"`
		RebateCcy   string               `json:"rebateCcy"`
		Px          types.JSONFloat64    `json:"px"`
		Sz          types.JSONFloat64    `json:"sz"`
		Pnl         types.JSONFloat64    `json:"pnl"`
		AccFillSz   types.JSONFloat64    `json:"accFillSz"`
		FillPx      types.JSONFloat64    `json:"fillPx"`
		FillSz      types.JSONFloat64    `json:"fillSz"`
		FillTime    types.JSONFloat64    `json:"fillTime"`
		AvgPx       types.JSONFloat64    `json:"avgPx"`
		Lever       types.JSONFloat64    `json:"lever"`
		TpTriggerPx types.JSONFloat64    `json:"tpTriggerPx"`
		TpOrdPx     types.JSONFloat64    `json:"tpOrdPx"`
		SlTriggerPx types.JSONFloat64    `json:"slTriggerPx"`
		SlOrdPx     types.JSONFloat64    `json:"slOrdPx"`
		Fee         types.JSONFloat64    `json:"fee"`
		Rebate      types.JSONFloat64    `json:"rebate"`
		State       types.OrderState     `json:"state"`
		TdMode      types.TradeMode      `json:"tdMode"`
		PosSide     types.PositionSide   `json:"posSide"`
		Side        types.OrderSide      `json:"side"`
		OrdType     types.OrderType      `json:"ordType"`
		InstType    types.InstrumentType `json:"instType"`
		TgtCcy      types.QuantityType   `json:"tgtCcy"`
		UTime       types.JSONTime       `json:"uTime"`
		CTime       types.JSONTime       `json:"cTime"`
	}
	TransactionDetail struct {
		InstID   string               `json:"instId"`
		OrdID    string               `json:"ordId"`
		TradeID  string               `json:"tradeId"`
		ClOrdID  string               `json:"clOrdId"`
		BillID   string               `json:"billId"`
		Tag      string               `json:"tag"`
		FillPx   types.JSONFloat64    `json:"fillPx"`
		FillSz   types.JSONFloat64    `json:"fillSz"`
		FeeCcy   string               `json:"feeCcy"`
		Fee      types.JSONFloat64    `json:"fee"`
		InstType types.InstrumentType `json:"instType"`
		Side     types.OrderSide      `json:"side"`
		PosSide  types.PositionSide   `json:"posSide"`
		ExecType types.OrderFlowType  `json:"execType"`
		TS       types.JSONTime       `json:"ts"`
	}
	PlaceAlgoOrder struct {
		AlgoID string          `json:"algoId"`
		SMsg   string          `json:"sMsg"`
		SCode  types.JSONInt64 `json:"sCode"`
	}
	CancelAlgoOrder struct {
		AlgoID string          `json:"algoId"`
		SMsg   string          `json:"sMsg"`
		SCode  types.JSONInt64 `json:"sCode"`
	}
	AlgoOrder struct {
		InstID       string               `json:"instId"`
		Ccy          string               `json:"ccy"`
		OrdID        string               `json:"ordId"`
		AlgoID       string               `json:"algoId"`
		ClOrdID      string               `json:"clOrdId"`
		TradeID      string               `json:"tradeId"`
		Tag          string               `json:"tag"`
		Category     string               `json:"category"`
		FeeCcy       string               `json:"feeCcy"`
		RebateCcy    string               `json:"rebateCcy"`
		TimeInterval string               `json:"timeInterval"`
		Px           types.JSONFloat64    `json:"px"`
		PxVar        types.JSONFloat64    `json:"pxVar"`
		PxSpread     types.JSONFloat64    `json:"pxSpread"`
		PxLimit      types.JSONFloat64    `json:"pxLimit"`
		Sz           types.JSONFloat64    `json:"sz"`
		SzLimit      types.JSONFloat64    `json:"szLimit"`
		ActualSz     types.JSONFloat64    `json:"actualSz"`
		ActualPx     types.JSONFloat64    `json:"actualPx"`
		Pnl          types.JSONFloat64    `json:"pnl"`
		AccFillSz    types.JSONFloat64    `json:"accFillSz"`
		FillPx       types.JSONFloat64    `json:"fillPx"`
		FillSz       types.JSONFloat64    `json:"fillSz"`
		FillTime     types.JSONFloat64    `json:"fillTime"`
		AvgPx        types.JSONFloat64    `json:"avgPx"`
		Lever        types.JSONFloat64    `json:"lever"`
		TpTriggerPx  types.JSONFloat64    `json:"tpTriggerPx"`
		TpOrdPx      types.JSONFloat64    `json:"tpOrdPx"`
		SlTriggerPx  types.JSONFloat64    `json:"slTriggerPx"`
		SlOrdPx      types.JSONFloat64    `json:"slOrdPx"`
		OrdPx        types.JSONFloat64    `json:"ordPx"`
		Fee          types.JSONFloat64    `json:"fee"`
		Rebate       types.JSONFloat64    `json:"rebate"`
		State        types.OrderState     `json:"state"`
		TdMode       types.TradeMode      `json:"tdMode"`
		ActualSide   types.PositionSide   `json:"actualSide"`
		PosSide      types.PositionSide   `json:"posSide"`
		Side         types.OrderSide      `json:"side"`
		OrdType      types.AlgoOrderType  `json:"ordType"`
		InstType     types.InstrumentType `json:"instType"`
		TgtCcy       types.QuantityType   `json:"tgtCcy"`
		CTime        types.JSONTime       `json:"cTime"`
		TriggerTime  types.JSONTime       `json:"triggerTime"`
	}
	OrderPreCheck struct {
		//TODO
	}
)
