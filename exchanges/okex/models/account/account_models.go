package account

import (
	"github.com/djpken/go-exc/exchanges/okex/constants"
)

type (
	Balance struct {
		TotalEq     constants.JSONFloat64 `json:"totalEq"`
		IsoEq       constants.JSONFloat64 `json:"isoEq"`
		AdjEq       constants.JSONFloat64 `json:"adjEq,omitempty"`
		OrdFroz     constants.JSONFloat64 `json:"ordFroz,omitempty"`
		Imr         constants.JSONFloat64 `json:"imr,omitempty"`
		Mmr         constants.JSONFloat64 `json:"mmr,omitempty"`
		MgnRatio    constants.JSONFloat64 `json:"mgnRatio,omitempty"`
		NotionalUsd constants.JSONFloat64 `json:"notionalUsd,omitempty"`
		Details     []*BalanceDetails `json:"details,omitempty"`
		UTime       constants.JSONTime    `json:"uTime"`
	}
	BalanceDetails struct {
		Ccy           string            `json:"ccy"`
		Eq            constants.JSONFloat64 `json:"eq"`
		CashBal       constants.JSONFloat64 `json:"cashBal"`
		IsoEq         constants.JSONFloat64 `json:"isoEq,omitempty"`
		AvailEq       constants.JSONFloat64 `json:"availEq,omitempty"`
		DisEq         constants.JSONFloat64 `json:"disEq"`
		AvailBal      constants.JSONFloat64 `json:"availBal"`
		FrozenBal     constants.JSONFloat64 `json:"frozenBal"`
		OrdFrozen     constants.JSONFloat64 `json:"ordFrozen"`
		Liab          constants.JSONFloat64 `json:"liab,omitempty"`
		Upl           constants.JSONFloat64 `json:"upl,omitempty"`
		UplLib        constants.JSONFloat64 `json:"uplLib,omitempty"`
		CrossLiab     constants.JSONFloat64 `json:"crossLiab,omitempty"`
		IsoLiab       constants.JSONFloat64 `json:"isoLiab,omitempty"`
		MgnRatio      constants.JSONFloat64 `json:"mgnRatio,omitempty"`
		Interest      constants.JSONFloat64 `json:"interest,omitempty"`
		Twap          constants.JSONFloat64 `json:"twap,omitempty"`
		MaxLoan       constants.JSONFloat64 `json:"maxLoan,omitempty"`
		EqUsd         constants.JSONFloat64 `json:"eqUsd"`
		NotionalLever constants.JSONFloat64 `json:"notionalLever,omitempty"`
		StgyEq        constants.JSONFloat64 `json:"stgyEq"`
		IsoUpl        constants.JSONFloat64 `json:"isoUpl,omitempty"`
		UTime         constants.JSONTime    `json:"uTime"`
	}
	Position struct {
		InstID      string               `json:"instId"`
		PosCcy      string               `json:"posCcy,omitempty"`
		LiabCcy     string               `json:"liabCcy,omitempty"`
		OptVal      string               `json:"optVal,omitempty"`
		Ccy         string               `json:"ccy"`
		PosID       string               `json:"posId"`
		TradeID     string               `json:"tradeId"`
		Pos         constants.JSONFloat64    `json:"pos"`
		AvailPos    constants.JSONFloat64    `json:"availPos,omitempty"`
		AvgPx       constants.JSONFloat64    `json:"avgPx"`
		Upl         constants.JSONFloat64    `json:"upl"`
		UplRatio    constants.JSONFloat64    `json:"uplRatio"`
		Lever       constants.JSONFloat64    `json:"lever"`
		LiqPx       constants.JSONFloat64    `json:"liqPx,omitempty"`
		Imr         constants.JSONFloat64    `json:"imr,omitempty"`
		Margin      constants.JSONFloat64    `json:"margin,omitempty"`
		MgnRatio    constants.JSONFloat64    `json:"mgnRatio"`
		Mmr         constants.JSONFloat64    `json:"mmr"`
		Liab        constants.JSONFloat64    `json:"liab,omitempty"`
		Interest    constants.JSONFloat64    `json:"interest"`
		NotionalUsd constants.JSONFloat64    `json:"notionalUsd"`
		ADL         constants.JSONFloat64    `json:"adl"`
		Last        constants.JSONFloat64    `json:"last"`
		DeltaBS     constants.JSONFloat64    `json:"deltaBS"`
		DeltaPA     constants.JSONFloat64    `json:"deltaPA"`
		GammaBS     constants.JSONFloat64    `json:"gammaBS"`
		GammaPA     constants.JSONFloat64    `json:"gammaPA"`
		ThetaBS     constants.JSONFloat64    `json:"thetaBS"`
		ThetaPA     constants.JSONFloat64    `json:"thetaPA"`
		VegaBS      constants.JSONFloat64    `json:"vegaBS"`
		VegaPA      constants.JSONFloat64    `json:"vegaPA"`
		PosSide     constants.PositionSide   `json:"posSide"`
		MgnMode     constants.MarginMode     `json:"mgnMode"`
		InstType    constants.InstrumentType `json:"instType"`
		CTime       constants.JSONTime       `json:"cTime"`
		UTime       constants.JSONTime       `json:"uTime"`
		RealizedPnl constants.JSONFloat64    `json:"realizedPnl"`
	}
	BalanceAndPosition struct {
		EventType constants.EventType `json:"eventType"`
		PTime     constants.JSONTime  `json:"pTime"`
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
		InstType       constants.InstrumentType `json:"instType"`
		MgnMode        constants.MarginMode     `json:"mgnMode"`
		PosSide        constants.PositionSide   `json:"posSide"`
		Pos            constants.JSONFloat64    `json:"pos"`
		Ccy            string               `json:"ccy"`
		PosCcy         string               `json:"posCcy"`
		AvgPx          constants.JSONFloat64    `json:"avgPx"`
		NonSettleAvgPx string               `json:"nonSettleAvgPx"`
		SettledPnl     constants.JSONFloat64    `json:"settledPnl"`
		UTime          constants.JSONTime       `json:"uTime"`
	}
	BalData struct {
		Ccy     string            `json:"ccy"`
		CashBal constants.JSONFloat64 `json:"cashBal"`
		UTime   constants.JSONTime    `json:"uTime"`
	}

	PositionAndAccountRisk struct {
		AdjEq   constants.JSONFloat64                    `json:"adjEq,omitempty"`
		BalData []*PositionAndAccountRiskBalanceData `json:"balData"`
		PosData []*PositionAndAccountRiskBalanceData `json:"posData"`
		TS      constants.JSONTime                       `json:"ts"`
	}
	PositionAndAccountRiskBalanceData struct {
		Ccy   string            `json:"ccy"`
		Eq    constants.JSONFloat64 `json:"eq"`
		DisEq constants.JSONFloat64 `json:"disEq"`
	}
	PositionAndAccountRiskPositionData struct {
		InstID      string               `json:"instId"`
		PosCcy      string               `json:"posCcy,omitempty"`
		Ccy         string               `json:"ccy"`
		NotionalCcy constants.JSONFloat64    `json:"notionalCcy"`
		Pos         constants.JSONFloat64    `json:"pos"`
		NotionalUsd constants.JSONFloat64    `json:"notionalUsd"`
		PosSide     constants.PositionSide   `json:"posSide"`
		InstType    constants.InstrumentType `json:"instType"`
		MgnMode     constants.MarginMode     `json:"mgnMode"`
	}
	Bill struct {
		Ccy       string               `json:"ccy"`
		InstID    string               `json:"instId"`
		Notes     string               `json:"notes"`
		BillID    string               `json:"billId"`
		OrdID     string               `json:"ordId"`
		BalChg    constants.JSONFloat64    `json:"balChg"`
		PosBalChg constants.JSONFloat64    `json:"posBalChg"`
		Bal       constants.JSONFloat64    `json:"bal"`
		PosBal    constants.JSONFloat64    `json:"posBal"`
		Sz        constants.JSONFloat64    `json:"sz"`
		Pnl       constants.JSONFloat64    `json:"pnl"`
		Fee       constants.JSONFloat64    `json:"fee"`
		From      constants.AccountType    `json:"from,string"`
		To        constants.AccountType    `json:"to,string"`
		InstType  constants.InstrumentType `json:"instType"`
		MgnMode   constants.MarginMode     `json:"MgnMode"`
		Type      constants.BillType       `json:"type,string"`
		SubType   constants.BillSubType    `json:"subType,string"`
		TS        constants.JSONTime       `json:"ts"`
	}
	Config struct {
		Level      string            `json:"level"`
		LevelTmp   string            `json:"levelTmp"`
		AcctLv     string            `json:"acctLv"`
		AutoLoan   bool              `json:"autoLoan"`
		UID        string            `json:"uid"`
		GreeksType constants.GreekType   `json:"greeksType"`
		PosMode    constants.PosModeType `json:"posMode"`
	}
	PositionMode struct {
		PosMode constants.PosModeType `json:"posMode"`
	}
	Leverage struct {
		InstID  string             `json:"instId"`
		Lever   constants.JSONFloat64  `json:"lever"`
		MgnMode constants.MarginMode   `json:"mgnMode"`
		PosSide constants.PositionSide `json:"posSide"`
	}
	MaxBuySellAmount struct {
		InstID  string            `json:"instId"`
		Ccy     string            `json:"ccy"`
		MaxBuy  constants.JSONFloat64 `json:"maxBuy"`
		MaxSell constants.JSONFloat64 `json:"maxSell"`
	}
	MaxAvailableTradeAmount struct {
		InstID    string            `json:"instId"`
		AvailBuy  constants.JSONFloat64 `json:"availBuy"`
		AvailSell constants.JSONFloat64 `json:"availSell"`
	}
	MarginBalanceAmount struct {
		InstID  string             `json:"instId"`
		Amt     constants.JSONFloat64  `json:"amt"`
		PosSide constants.PositionSide `json:"posSide,string"`
		Type    constants.CountAction  `json:"type,string"`
	}
	Loan struct {
		InstID  string            `json:"instId"`
		MgnCcy  string            `json:"mgnCcy"`
		Ccy     string            `json:"ccy"`
		MaxLoan constants.JSONFloat64 `json:"maxLoan"`
		MgnMode constants.MarginMode  `json:"mgnMode"`
		Side    constants.OrderSide   `json:"side,string"`
	}
	Fee struct {
		Level    string               `json:"level"`
		Taker    constants.JSONFloat64    `json:"taker"`
		Maker    constants.JSONFloat64    `json:"maker"`
		Delivery constants.JSONFloat64    `json:"delivery,omitempty"`
		Exercise constants.JSONFloat64    `json:"exercise,omitempty"`
		Category constants.FeeCategory    `json:"category,string"`
		InstType constants.InstrumentType `json:"instType"`
		TS       constants.JSONTime       `json:"ts"`
	}
	InterestAccrued struct {
		InstID       string            `json:"instId"`
		Ccy          string            `json:"ccy"`
		Interest     constants.JSONFloat64 `json:"interest"`
		InterestRate constants.JSONFloat64 `json:"interestRate"`
		Liab         constants.JSONFloat64 `json:"liab"`
		MgnMode      constants.MarginMode  `json:"mgnMode"`
		TS           constants.JSONTime    `json:"ts"`
	}
	InterestRate struct {
		Ccy          string            `json:"ccy"`
		InterestRate constants.JSONFloat64 `json:"interestRate"`
	}
	Greek struct {
		GreeksType string `json:"greeksType"`
	}
	MaxWithdrawal struct {
		Ccy   string            `json:"ccy"`
		MaxWd constants.JSONFloat64 `json:"maxWd"`
	}
	AccountLevel struct {
		AcctLv string `json:"acctLv"`
	}
)
