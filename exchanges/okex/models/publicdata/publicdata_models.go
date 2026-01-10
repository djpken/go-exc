package publicdata

import (
	"github.com/djpken/go-exc/exchanges/okex/types"
)

type (
	Instrument struct {
		InstID    string                `json:"instId"`
		Uly       string                `json:"uly,omitempty"`
		BaseCcy   string                `json:"baseCcy,omitempty"`
		QuoteCcy  string                `json:"quoteCcy,omitempty"`
		SettleCcy string                `json:"settleCcy,omitempty"`
		CtValCcy  string                `json:"ctValCcy,omitempty"`
		CtVal     types.JSONFloat64     `json:"ctVal,omitempty"`
		CtMult    types.JSONFloat64     `json:"ctMult,omitempty"`
		Stk       types.JSONFloat64     `json:"stk,omitempty"`
		TickSz    types.JSONFloat64     `json:"tickSz,omitempty"`
		LotSz     types.JSONFloat64     `json:"lotSz,omitempty"`
		MinSz     types.JSONFloat64     `json:"minSz,omitempty"`
		Lever     types.JSONFloat64     `json:"lever"`
		InstType  types.InstrumentType  `json:"instType"`
		Category  types.FeeCategory     `json:"category,string"`
		OptType   types.OptionType      `json:"optType,omitempty"`
		ListTime  types.JSONTime        `json:"listTime"`
		ExpTime   types.JSONTime        `json:"expTime,omitempty"`
		CtType    types.ContractType    `json:"ctType,omitempty"`
		Alias     types.AliasType       `json:"alias,omitempty"`
		State     types.InstrumentState `json:"state"`
	}
	DeliveryExerciseHistory struct {
		Details []*DeliveryExerciseHistoryDetails `json:"details"`
		TS      types.JSONTime                    `json:"ts"`
	}
	DeliveryExerciseHistoryDetails struct {
		InstID string                     `json:"instId"`
		Px     types.JSONFloat64          `json:"px"`
		Type   types.DeliveryExerciseType `json:"type"`
	}
	OpenInterest struct {
		InstID   string               `json:"instId"`
		Oi       types.JSONFloat64    `json:"oi"`
		OiCcy    types.JSONFloat64    `json:"oiCcy"`
		InstType types.InstrumentType `json:"instType"`
		TS       types.JSONTime       `json:"ts"`
	}
	FundingRate struct {
		InstID          string               `json:"instId"`
		InstType        types.InstrumentType `json:"instType"`
		FundingRate     types.JSONFloat64    `json:"fundingRate"`
		NextFundingRate types.JSONFloat64    `json:"NextFundingRate"`
		FundingTime     types.JSONTime       `json:"fundingTime"`
		NextFundingTime types.JSONTime       `json:"nextFundingTime"`
	}
	LimitPrice struct {
		InstID   string               `json:"instId"`
		InstType types.InstrumentType `json:"instType"`
		BuyLmt   types.JSONFloat64    `json:"buyLmt"`
		SellLmt  types.JSONFloat64    `json:"sellLmt"`
		TS       types.JSONTime       `json:"ts"`
	}
	EstimatedDeliveryExercisePrice struct {
		InstID   string               `json:"instId"`
		InstType types.InstrumentType `json:"instType"`
		SettlePx types.JSONFloat64    `json:"settlePx"`
		TS       types.JSONTime       `json:"ts"`
	}
	OptionMarketData struct {
		InstID   string               `json:"instId"`
		Uly      string               `json:"uly"`
		InstType types.InstrumentType `json:"instType"`
		Delta    types.JSONFloat64    `json:"delta"`
		Gamma    types.JSONFloat64    `json:"gamma"`
		Vega     types.JSONFloat64    `json:"vega"`
		Theta    types.JSONFloat64    `json:"theta"`
		DeltaBS  types.JSONFloat64    `json:"deltaBS"`
		GammaBS  types.JSONFloat64    `json:"gammaBS"`
		VegaBS   types.JSONFloat64    `json:"vegaBS"`
		ThetaBS  types.JSONFloat64    `json:"thetaBS"`
		Lever    types.JSONFloat64    `json:"lever"`
		MarkVol  types.JSONFloat64    `json:"markVol"`
		BidVol   types.JSONFloat64    `json:"bidVol"`
		AskVol   types.JSONFloat64    `json:"askVol"`
		RealVol  types.JSONFloat64    `json:"realVol"`
		TS       types.JSONTime       `json:"ts"`
	}
	GetDiscountRateAndInterestFreeQuota struct {
		Ccy          string            `json:"ccy"`
		Amt          types.JSONFloat64 `json:"amt"`
		DiscountLv   types.JSONInt64   `json:"discountLv"`
		DiscountInfo []*DiscountInfo   `json:"discountInfo"`
	}
	DiscountInfo struct {
		DiscountRate types.JSONInt64 `json:"discountRate"`
		MaxAmt       types.JSONInt64 `json:"maxAmt"`
		MinAmt       types.JSONInt64 `json:"minAmt"`
	}
	SystemTime struct {
		TS types.JSONTime `json:"ts"`
	}
	LiquidationOrder struct {
		InstID    string                    `json:"instId"`
		Uly       string                    `json:"uly,omitempty"`
		InstType  types.InstrumentType      `json:"instType"`
		TotalLoss types.JSONFloat64         `json:"totalLoss"`
		Details   []*LiquidationOrderDetail `json:"details"`
	}
	LiquidationOrderDetail struct {
		Ccy     string             `json:"ccy,omitempty"`
		Side    types.OrderSide    `json:"side"`
		OosSide types.PositionSide `json:"posSide"`
		BkPx    types.JSONFloat64  `json:"bkPx"`
		Sz      types.JSONFloat64  `json:"sz"`
		BkLoss  types.JSONFloat64  `json:"bkLoss"`
		TS      types.JSONTime     `json:"ts"`
	}
	MarkPrice struct {
		InstID   string               `json:"instId"`
		InstType types.InstrumentType `json:"instType"`
		MarkPx   types.JSONFloat64    `json:"markPx"`
		TS       types.JSONTime       `json:"ts"`
	}
	PositionTier struct {
		InstID       string               `json:"instId"`
		Uly          string               `json:"uly,omitempty"`
		InstType     types.InstrumentType `json:"instType"`
		Tier         types.JSONInt64      `json:"tier"`
		MinSz        types.JSONFloat64    `json:"minSz"`
		MaxSz        types.JSONFloat64    `json:"maxSz"`
		Mmr          types.JSONFloat64    `json:"mmr"`
		Imr          types.JSONFloat64    `json:"imr"`
		OptMgnFactor types.JSONFloat64    `json:"optMgnFactor,omitempty"`
		QuoteMaxLoan types.JSONFloat64    `json:"quoteMaxLoan,omitempty"`
		BaseMaxLoan  types.JSONFloat64    `json:"baseMaxLoan,omitempty"`
		MaxLever     types.JSONFloat64    `json:"maxLever"`
		TS           types.JSONTime       `json:"ts"`
	}
	InterestRateAndLoanQuota struct {
		Basic   []*InterestRateAndLoanBasic `json:"basic"`
		Vip     []*InterestRateAndLoanUser  `json:"vip"`
		Regular []*InterestRateAndLoanUser  `json:"regular"`
	}
	InterestRateAndLoanBasic struct {
		Ccy   string            `json:"ccy"`
		Rate  types.JSONFloat64 `json:"rate"`
		Quota types.JSONFloat64 `json:"quota"`
	}
	InterestRateAndLoanUser struct {
		Level         string            `json:"level"`
		IrDiscount    types.JSONFloat64 `json:"irDiscount"`
		LoanQuotaCoef int               `json:"loanQuotaCoef,string"`
	}
	State struct {
		Title       string         `json:"title"`
		State       string         `json:"state"`
		Href        string         `json:"href"`
		ServiceType string         `json:"serviceType"`
		System      string         `json:"system"`
		ScheDesc    string         `json:"scheDesc"`
		Begin       types.JSONTime `json:"begin"`
		End         types.JSONTime `json:"end"`
	}
)
