package publicdata

import (
	"github.com/djpken/go-exc/exchanges/okex/constants"
)

type (
	Instrument struct {
		InstID    string                    `json:"instId"`
		Uly       string                    `json:"uly,omitempty"`
		BaseCcy   string                    `json:"baseCcy,omitempty"`
		QuoteCcy  string                    `json:"quoteCcy,omitempty"`
		SettleCcy string                    `json:"settleCcy,omitempty"`
		CtValCcy  string                    `json:"ctValCcy,omitempty"`
		CtVal     constants.JSONFloat64     `json:"ctVal,omitempty"`
		CtMult    constants.JSONFloat64     `json:"ctMult,omitempty"`
		Stk       constants.JSONFloat64     `json:"stk,omitempty"`
		TickSz    constants.JSONFloat64     `json:"tickSz,omitempty"`
		LotSz     constants.JSONFloat64     `json:"lotSz,omitempty"`
		MinSz     constants.JSONFloat64     `json:"minSz,omitempty"`
		Lever     constants.JSONInt64       `json:"lever"`
		InstType  constants.InstrumentType  `json:"instType"`
		Category  constants.FeeCategory     `json:"category,string"`
		OptType   constants.OptionType      `json:"optType,omitempty"`
		ListTime  constants.JSONTime        `json:"listTime"`
		ExpTime   constants.JSONTime        `json:"expTime,omitempty"`
		CtType    constants.ContractType    `json:"ctType,omitempty"`
		Alias     constants.AliasType       `json:"alias,omitempty"`
		State     constants.InstrumentState `json:"state"`
	}
	DeliveryExerciseHistory struct {
		Details []*DeliveryExerciseHistoryDetails `json:"details"`
		TS      constants.JSONTime                `json:"ts"`
	}
	DeliveryExerciseHistoryDetails struct {
		InstID string                         `json:"instId"`
		Px     constants.JSONFloat64          `json:"px"`
		Type   constants.DeliveryExerciseType `json:"type"`
	}
	OpenInterest struct {
		InstID   string                   `json:"instId"`
		Oi       constants.JSONFloat64    `json:"oi"`
		OiCcy    constants.JSONFloat64    `json:"oiCcy"`
		InstType constants.InstrumentType `json:"instType"`
		TS       constants.JSONTime       `json:"ts"`
	}
	FundingRate struct {
		InstID          string                   `json:"instId"`
		InstType        constants.InstrumentType `json:"instType"`
		FundingRate     constants.JSONFloat64    `json:"fundingRate"`
		NextFundingRate constants.JSONFloat64    `json:"NextFundingRate"`
		FundingTime     constants.JSONTime       `json:"fundingTime"`
		NextFundingTime constants.JSONTime       `json:"nextFundingTime"`
	}
	LimitPrice struct {
		InstID   string                   `json:"instId"`
		InstType constants.InstrumentType `json:"instType"`
		BuyLmt   constants.JSONFloat64    `json:"buyLmt"`
		SellLmt  constants.JSONFloat64    `json:"sellLmt"`
		TS       constants.JSONTime       `json:"ts"`
	}
	EstimatedDeliveryExercisePrice struct {
		InstID   string                   `json:"instId"`
		InstType constants.InstrumentType `json:"instType"`
		SettlePx constants.JSONFloat64    `json:"settlePx"`
		TS       constants.JSONTime       `json:"ts"`
	}
	OptionMarketData struct {
		InstID   string                   `json:"instId"`
		Uly      string                   `json:"uly"`
		InstType constants.InstrumentType `json:"instType"`
		Delta    constants.JSONFloat64    `json:"delta"`
		Gamma    constants.JSONFloat64    `json:"gamma"`
		Vega     constants.JSONFloat64    `json:"vega"`
		Theta    constants.JSONFloat64    `json:"theta"`
		DeltaBS  constants.JSONFloat64    `json:"deltaBS"`
		GammaBS  constants.JSONFloat64    `json:"gammaBS"`
		VegaBS   constants.JSONFloat64    `json:"vegaBS"`
		ThetaBS  constants.JSONFloat64    `json:"thetaBS"`
		Lever    constants.JSONFloat64    `json:"lever"`
		MarkVol  constants.JSONFloat64    `json:"markVol"`
		BidVol   constants.JSONFloat64    `json:"bidVol"`
		AskVol   constants.JSONFloat64    `json:"askVol"`
		RealVol  constants.JSONFloat64    `json:"realVol"`
		TS       constants.JSONTime       `json:"ts"`
	}
	GetDiscountRateAndInterestFreeQuota struct {
		Ccy          string                `json:"ccy"`
		Amt          constants.JSONFloat64 `json:"amt"`
		DiscountLv   constants.JSONInt64   `json:"discountLv"`
		DiscountInfo []*DiscountInfo       `json:"discountInfo"`
	}
	DiscountInfo struct {
		DiscountRate constants.JSONInt64 `json:"discountRate"`
		MaxAmt       constants.JSONInt64 `json:"maxAmt"`
		MinAmt       constants.JSONInt64 `json:"minAmt"`
	}
	SystemTime struct {
		TS constants.JSONTime `json:"ts"`
	}
	LiquidationOrder struct {
		InstID    string                    `json:"instId"`
		Uly       string                    `json:"uly,omitempty"`
		InstType  constants.InstrumentType  `json:"instType"`
		TotalLoss constants.JSONFloat64     `json:"totalLoss"`
		Details   []*LiquidationOrderDetail `json:"details"`
	}
	LiquidationOrderDetail struct {
		Ccy     string                 `json:"ccy,omitempty"`
		Side    constants.OrderSide    `json:"side"`
		OosSide constants.PositionSide `json:"posSide"`
		BkPx    constants.JSONFloat64  `json:"bkPx"`
		Sz      constants.JSONFloat64  `json:"sz"`
		BkLoss  constants.JSONFloat64  `json:"bkLoss"`
		TS      constants.JSONTime     `json:"ts"`
	}
	MarkPrice struct {
		InstID   string                   `json:"instId"`
		InstType constants.InstrumentType `json:"instType"`
		MarkPx   constants.JSONFloat64    `json:"markPx"`
		TS       constants.JSONTime       `json:"ts"`
	}
	PositionTier struct {
		InstID       string                   `json:"instId"`
		Uly          string                   `json:"uly,omitempty"`
		InstType     constants.InstrumentType `json:"instType"`
		Tier         constants.JSONInt64      `json:"tier"`
		MinSz        constants.JSONFloat64    `json:"minSz"`
		MaxSz        constants.JSONFloat64    `json:"maxSz"`
		Mmr          constants.JSONFloat64    `json:"mmr"`
		Imr          constants.JSONFloat64    `json:"imr"`
		OptMgnFactor constants.JSONFloat64    `json:"optMgnFactor,omitempty"`
		QuoteMaxLoan constants.JSONFloat64    `json:"quoteMaxLoan,omitempty"`
		BaseMaxLoan  constants.JSONFloat64    `json:"baseMaxLoan,omitempty"`
		MaxLever     constants.JSONFloat64    `json:"maxLever"`
		TS           constants.JSONTime       `json:"ts"`
	}
	InterestRateAndLoanQuota struct {
		Basic   []*InterestRateAndLoanBasic `json:"basic"`
		Vip     []*InterestRateAndLoanUser  `json:"vip"`
		Regular []*InterestRateAndLoanUser  `json:"regular"`
	}
	InterestRateAndLoanBasic struct {
		Ccy   string                `json:"ccy"`
		Rate  constants.JSONFloat64 `json:"rate"`
		Quota constants.JSONFloat64 `json:"quota"`
	}
	InterestRateAndLoanUser struct {
		Level         string                `json:"level"`
		IrDiscount    constants.JSONFloat64 `json:"irDiscount"`
		LoanQuotaCoef int                   `json:"loanQuotaCoef,string"`
	}
	State struct {
		Title       string             `json:"title"`
		State       string             `json:"state"`
		Href        string             `json:"href"`
		ServiceType string             `json:"serviceType"`
		System      string             `json:"system"`
		ScheDesc    string             `json:"scheDesc"`
		Begin       constants.JSONTime `json:"begin"`
		End         constants.JSONTime `json:"end"`
	}
)
