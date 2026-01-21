package trade

import (
	"github.com/djpken/go-exc/exchanges/okex/constants"
)

type (
	PlaceOrder struct {
		ClOrdID string            `json:"clOrdId"`
		Tag     string            `json:"tag"`
		SMsg    string            `json:"sMsg"`
		SCode   constants.JSONInt64   `json:"sCode"`
		OrdID   constants.JSONFloat64 `json:"ordId"`
	}
	CancelOrder struct {
		OrdID   string            `json:"ordId"`
		ClOrdID string            `json:"clOrdId"`
		SMsg    string            `json:"sMsg"`
		SCode   constants.JSONFloat64 `json:"sCode"`
	}
	AmendOrder struct {
		OrdID   string            `json:"ordId"`
		ClOrdID string            `json:"clOrdId"`
		ReqID   string            `json:"reqId"`
		SMsg    string            `json:"sMsg"`
		SCode   constants.JSONFloat64 `json:"sCode"`
	}
	ClosePosition struct {
		InstID  string             `json:"instId"`
		PosSide constants.PositionSide `json:"posSide"`
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
		Px          constants.JSONFloat64    `json:"px"`
		Sz          constants.JSONFloat64    `json:"sz"`
		Pnl         constants.JSONFloat64    `json:"pnl"`
		AccFillSz   constants.JSONFloat64    `json:"accFillSz"`
		FillPx      constants.JSONFloat64    `json:"fillPx"`
		FillSz      constants.JSONFloat64    `json:"fillSz"`
		FillTime    constants.JSONFloat64    `json:"fillTime"`
		AvgPx       constants.JSONFloat64    `json:"avgPx"`
		Lever       constants.JSONFloat64    `json:"lever"`
		TpTriggerPx constants.JSONFloat64    `json:"tpTriggerPx"`
		TpOrdPx     constants.JSONFloat64    `json:"tpOrdPx"`
		SlTriggerPx constants.JSONFloat64    `json:"slTriggerPx"`
		SlOrdPx     constants.JSONFloat64    `json:"slOrdPx"`
		Fee         constants.JSONFloat64    `json:"fee"`
		Rebate      constants.JSONFloat64    `json:"rebate"`
		State       constants.OrderState     `json:"state"`
		TdMode      constants.TradeMode      `json:"tdMode"`
		PosSide     constants.PositionSide   `json:"posSide"`
		Side        constants.OrderSide      `json:"side"`
		OrdType     constants.OrderType      `json:"ordType"`
		InstType    constants.InstrumentType `json:"instType"`
		TgtCcy      constants.QuantityType   `json:"tgtCcy"`
		UTime       constants.JSONTime       `json:"uTime"`
		CTime       constants.JSONTime       `json:"cTime"`
	}
	TransactionDetail struct {
		InstID   string               `json:"instId"`
		OrdID    string               `json:"ordId"`
		TradeID  string               `json:"tradeId"`
		ClOrdID  string               `json:"clOrdId"`
		BillID   string               `json:"billId"`
		Tag      string               `json:"tag"`
		FillPx   constants.JSONFloat64    `json:"fillPx"`
		FillSz   constants.JSONFloat64    `json:"fillSz"`
		FeeCcy   string               `json:"feeCcy"`
		Fee      constants.JSONFloat64    `json:"fee"`
		InstType constants.InstrumentType `json:"instType"`
		Side     constants.OrderSide      `json:"side"`
		PosSide  constants.PositionSide   `json:"posSide"`
		ExecType constants.OrderFlowType  `json:"execType"`
		TS       constants.JSONTime       `json:"ts"`
	}
	PlaceAlgoOrder struct {
		AlgoID string          `json:"algoId"`
		SMsg   string          `json:"sMsg"`
		SCode  constants.JSONInt64 `json:"sCode"`
	}
	CancelAlgoOrder struct {
		AlgoID string          `json:"algoId"`
		SMsg   string          `json:"sMsg"`
		SCode  constants.JSONInt64 `json:"sCode"`
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
		Px           constants.JSONFloat64    `json:"px"`
		PxVar        constants.JSONFloat64    `json:"pxVar"`
		PxSpread     constants.JSONFloat64    `json:"pxSpread"`
		PxLimit      constants.JSONFloat64    `json:"pxLimit"`
		Sz           constants.JSONFloat64    `json:"sz"`
		SzLimit      constants.JSONFloat64    `json:"szLimit"`
		ActualSz     constants.JSONFloat64    `json:"actualSz"`
		ActualPx     constants.JSONFloat64    `json:"actualPx"`
		Pnl          constants.JSONFloat64    `json:"pnl"`
		AccFillSz    constants.JSONFloat64    `json:"accFillSz"`
		FillPx       constants.JSONFloat64    `json:"fillPx"`
		FillSz       constants.JSONFloat64    `json:"fillSz"`
		FillTime     constants.JSONFloat64    `json:"fillTime"`
		AvgPx        constants.JSONFloat64    `json:"avgPx"`
		Lever        constants.JSONFloat64    `json:"lever"`
		TpTriggerPx  constants.JSONFloat64    `json:"tpTriggerPx"`
		TpOrdPx      constants.JSONFloat64    `json:"tpOrdPx"`
		SlTriggerPx  constants.JSONFloat64    `json:"slTriggerPx"`
		SlOrdPx      constants.JSONFloat64    `json:"slOrdPx"`
		OrdPx        constants.JSONFloat64    `json:"ordPx"`
		Fee          constants.JSONFloat64    `json:"fee"`
		Rebate       constants.JSONFloat64    `json:"rebate"`
		State        constants.OrderState     `json:"state"`
		TdMode       constants.TradeMode      `json:"tdMode"`
		ActualSide   constants.PositionSide   `json:"actualSide"`
		PosSide      constants.PositionSide   `json:"posSide"`
		Side         constants.OrderSide      `json:"side"`
		OrdType      constants.AlgoOrderType  `json:"ordType"`
		InstType     constants.InstrumentType `json:"instType"`
		TgtCcy       constants.QuantityType   `json:"tgtCcy"`
		CTime        constants.JSONTime       `json:"cTime"`
		TriggerTime  constants.JSONTime       `json:"triggerTime"`
	}
	OrderPreCheck struct {
		//TODO
	}
)
