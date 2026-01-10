package account

import (
	"github.com/djpken/go-exc/exchanges/okex/types"
)

type (
	Balance struct {
		TotalEq     types.JSONFloat64 `json:"totalEq"`
		IsoEq       types.JSONFloat64 `json:"isoEq"`
		AdjEq       types.JSONFloat64 `json:"adjEq,omitempty"`
		OrdFroz     types.JSONFloat64 `json:"ordFroz,omitempty"`
		Imr         types.JSONFloat64 `json:"imr,omitempty"`
		Mmr         types.JSONFloat64 `json:"mmr,omitempty"`
		MgnRatio    types.JSONFloat64 `json:"mgnRatio,omitempty"`
		NotionalUsd types.JSONFloat64 `json:"notionalUsd,omitempty"`
		Details     []*BalanceDetails `json:"details,omitempty"`
		UTime       types.JSONTime    `json:"uTime"`
	}
	BalanceDetails struct {
		Ccy           string            `json:"ccy"`
		Eq            types.JSONFloat64 `json:"eq"`
		CashBal       types.JSONFloat64 `json:"cashBal"`
		IsoEq         types.JSONFloat64 `json:"isoEq,omitempty"`
		AvailEq       types.JSONFloat64 `json:"availEq,omitempty"`
		DisEq         types.JSONFloat64 `json:"disEq"`
		AvailBal      types.JSONFloat64 `json:"availBal"`
		FrozenBal     types.JSONFloat64 `json:"frozenBal"`
		OrdFrozen     types.JSONFloat64 `json:"ordFrozen"`
		Liab          types.JSONFloat64 `json:"liab,omitempty"`
		Upl           types.JSONFloat64 `json:"upl,omitempty"`
		UplLib        types.JSONFloat64 `json:"uplLib,omitempty"`
		CrossLiab     types.JSONFloat64 `json:"crossLiab,omitempty"`
		IsoLiab       types.JSONFloat64 `json:"isoLiab,omitempty"`
		MgnRatio      types.JSONFloat64 `json:"mgnRatio,omitempty"`
		Interest      types.JSONFloat64 `json:"interest,omitempty"`
		Twap          types.JSONFloat64 `json:"twap,omitempty"`
		MaxLoan       types.JSONFloat64 `json:"maxLoan,omitempty"`
		EqUsd         types.JSONFloat64 `json:"eqUsd"`
		NotionalLever types.JSONFloat64 `json:"notionalLever,omitempty"`
		StgyEq        types.JSONFloat64 `json:"stgyEq"`
		IsoUpl        types.JSONFloat64 `json:"isoUpl,omitempty"`
		UTime         types.JSONTime    `json:"uTime"`
	}
	Position struct {
		InstID      string               `json:"instId"`
		PosCcy      string               `json:"posCcy,omitempty"`
		LiabCcy     string               `json:"liabCcy,omitempty"`
		OptVal      string               `json:"optVal,omitempty"`
		Ccy         string               `json:"ccy"`
		PosID       string               `json:"posId"`
		TradeID     string               `json:"tradeId"`
		Pos         types.JSONFloat64    `json:"pos"`
		AvailPos    types.JSONFloat64    `json:"availPos,omitempty"`
		AvgPx       types.JSONFloat64    `json:"avgPx"`
		Upl         types.JSONFloat64    `json:"upl"`
		UplRatio    types.JSONFloat64    `json:"uplRatio"`
		Lever       types.JSONFloat64    `json:"lever"`
		LiqPx       types.JSONFloat64    `json:"liqPx,omitempty"`
		Imr         types.JSONFloat64    `json:"imr,omitempty"`
		Margin      types.JSONFloat64    `json:"margin,omitempty"`
		MgnRatio    types.JSONFloat64    `json:"mgnRatio"`
		Mmr         types.JSONFloat64    `json:"mmr"`
		Liab        types.JSONFloat64    `json:"liab,omitempty"`
		Interest    types.JSONFloat64    `json:"interest"`
		NotionalUsd types.JSONFloat64    `json:"notionalUsd"`
		ADL         types.JSONFloat64    `json:"adl"`
		Last        types.JSONFloat64    `json:"last"`
		DeltaBS     types.JSONFloat64    `json:"deltaBS"`
		DeltaPA     types.JSONFloat64    `json:"deltaPA"`
		GammaBS     types.JSONFloat64    `json:"gammaBS"`
		GammaPA     types.JSONFloat64    `json:"gammaPA"`
		ThetaBS     types.JSONFloat64    `json:"thetaBS"`
		ThetaPA     types.JSONFloat64    `json:"thetaPA"`
		VegaBS      types.JSONFloat64    `json:"vegaBS"`
		VegaPA      types.JSONFloat64    `json:"vegaPA"`
		PosSide     types.PositionSide   `json:"posSide"`
		MgnMode     types.MarginMode     `json:"mgnMode"`
		InstType    types.InstrumentType `json:"instType"`
		CTime       types.JSONTime       `json:"cTime"`
		UTime       types.JSONTime       `json:"uTime"`
		RealizedPnl types.JSONFloat64    `json:"realizedPnl"`
	}
	BalanceAndPosition struct {
		EventType types.EventType `json:"eventType"`
		PTime     types.JSONTime  `json:"pTime"`
		PosData   []*PosData      `json:"posData"`
		BalData   []*BalData      `json:"balData"`
		Trades    []*Trades       `json:"trades"`
	}
	Trades struct {
		InstId  string `json:"instId"`
		TradeId string `json:"tradeId"`
	}

	PosData struct {
		PosId          string               `json:"posId"`
		TradeId        string               `json:"tradeId"`
		InstId         string               `json:"instId"`
		InstType       types.InstrumentType `json:"instType"`
		MgnMode        types.MarginMode     `json:"mgnMode"`
		PosSide        types.PositionSide   `json:"posSide"`
		Pos            types.JSONFloat64    `json:"pos"`
		Ccy            string               `json:"ccy"`
		PosCcy         string               `json:"posCcy"`
		AvgPx          types.JSONFloat64    `json:"avgPx"`
		NonSettleAvgPx string               `json:"nonSettleAvgPx"`
		SettledPnl     types.JSONFloat64    `json:"settledPnl"`
		UTime          types.JSONTime       `json:"uTime"`
	}
	BalData struct {
		Ccy     string            `json:"ccy"`
		CashBal types.JSONFloat64 `json:"cashBal"`
		UTime   types.JSONTime    `json:"uTime"`
	}

	PositionAndAccountRisk struct {
		AdjEq   types.JSONFloat64                    `json:"adjEq,omitempty"`
		BalData []*PositionAndAccountRiskBalanceData `json:"balData"`
		PosData []*PositionAndAccountRiskBalanceData `json:"posData"`
		TS      types.JSONTime                       `json:"ts"`
	}
	PositionAndAccountRiskBalanceData struct {
		Ccy   string            `json:"ccy"`
		Eq    types.JSONFloat64 `json:"eq"`
		DisEq types.JSONFloat64 `json:"disEq"`
	}
	PositionAndAccountRiskPositionData struct {
		InstID      string               `json:"instId"`
		PosCcy      string               `json:"posCcy,omitempty"`
		Ccy         string               `json:"ccy"`
		NotionalCcy types.JSONFloat64    `json:"notionalCcy"`
		Pos         types.JSONFloat64    `json:"pos"`
		NotionalUsd types.JSONFloat64    `json:"notionalUsd"`
		PosSide     types.PositionSide   `json:"posSide"`
		InstType    types.InstrumentType `json:"instType"`
		MgnMode     types.MarginMode     `json:"mgnMode"`
	}
	Bill struct {
		Ccy       string               `json:"ccy"`
		InstID    string               `json:"instId"`
		Notes     string               `json:"notes"`
		BillID    string               `json:"billId"`
		OrdID     string               `json:"ordId"`
		BalChg    types.JSONFloat64    `json:"balChg"`
		PosBalChg types.JSONFloat64    `json:"posBalChg"`
		Bal       types.JSONFloat64    `json:"bal"`
		PosBal    types.JSONFloat64    `json:"posBal"`
		Sz        types.JSONFloat64    `json:"sz"`
		Pnl       types.JSONFloat64    `json:"pnl"`
		Fee       types.JSONFloat64    `json:"fee"`
		From      types.AccountType    `json:"from,string"`
		To        types.AccountType    `json:"to,string"`
		InstType  types.InstrumentType `json:"instType"`
		MgnMode   types.MarginMode     `json:"MgnMode"`
		Type      types.BillType       `json:"type,string"`
		SubType   types.BillSubType    `json:"subType,string"`
		TS        types.JSONTime       `json:"ts"`
	}
	Config struct {
		Level      string            `json:"level"`
		LevelTmp   string            `json:"levelTmp"`
		AcctLv     string            `json:"acctLv"`
		AutoLoan   bool              `json:"autoLoan"`
		UID        string            `json:"uid"`
		GreeksType types.GreekType   `json:"greeksType"`
		PosMode    types.PosModeType `json:"posMode"`
	}
	PositionMode struct {
		PosMode types.PosModeType `json:"posMode"`
	}
	Leverage struct {
		InstID  string             `json:"instId"`
		Lever   types.JSONFloat64  `json:"lever"`
		MgnMode types.MarginMode   `json:"mgnMode"`
		PosSide types.PositionSide `json:"posSide"`
	}
	MaxBuySellAmount struct {
		InstID  string            `json:"instId"`
		Ccy     string            `json:"ccy"`
		MaxBuy  types.JSONFloat64 `json:"maxBuy"`
		MaxSell types.JSONFloat64 `json:"maxSell"`
	}
	MaxAvailableTradeAmount struct {
		InstID    string            `json:"instId"`
		AvailBuy  types.JSONFloat64 `json:"availBuy"`
		AvailSell types.JSONFloat64 `json:"availSell"`
	}
	MarginBalanceAmount struct {
		InstID  string             `json:"instId"`
		Amt     types.JSONFloat64  `json:"amt"`
		PosSide types.PositionSide `json:"posSide,string"`
		Type    types.CountAction  `json:"type,string"`
	}
	Loan struct {
		InstID  string            `json:"instId"`
		MgnCcy  string            `json:"mgnCcy"`
		Ccy     string            `json:"ccy"`
		MaxLoan types.JSONFloat64 `json:"maxLoan"`
		MgnMode types.MarginMode  `json:"mgnMode"`
		Side    types.OrderSide   `json:"side,string"`
	}
	Fee struct {
		Level    string               `json:"level"`
		Taker    types.JSONFloat64    `json:"taker"`
		Maker    types.JSONFloat64    `json:"maker"`
		Delivery types.JSONFloat64    `json:"delivery,omitempty"`
		Exercise types.JSONFloat64    `json:"exercise,omitempty"`
		Category types.FeeCategory    `json:"category,string"`
		InstType types.InstrumentType `json:"instType"`
		TS       types.JSONTime       `json:"ts"`
	}
	InterestAccrued struct {
		InstID       string            `json:"instId"`
		Ccy          string            `json:"ccy"`
		Interest     types.JSONFloat64 `json:"interest"`
		InterestRate types.JSONFloat64 `json:"interestRate"`
		Liab         types.JSONFloat64 `json:"liab"`
		MgnMode      types.MarginMode  `json:"mgnMode"`
		TS           types.JSONTime    `json:"ts"`
	}
	InterestRate struct {
		Ccy          string            `json:"ccy"`
		InterestRate types.JSONFloat64 `json:"interestRate"`
	}
	Greek struct {
		GreeksType string `json:"greeksType"`
	}
	MaxWithdrawal struct {
		Ccy   string            `json:"ccy"`
		MaxWd types.JSONFloat64 `json:"maxWd"`
	}
	AccountLevel struct {
		AcctLv string `json:"acctLv"`
	}
)
