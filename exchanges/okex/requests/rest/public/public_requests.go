package public

import "github.com/djpken/go-exc/exchanges/okex/constants"

type (
	GetInstruments struct {
		Uly      string               `json:"uly,omitempty"`
		InstID   string               `json:"instId,omitempty"`
		InstType constants.InstrumentType `json:"instType"`
	}
	GetDeliveryExerciseHistory struct {
		Uly      string               `json:"uly"`
		After    int64                `json:"after,omitempty,string"`
		Before   int64                `json:"before,omitempty,string"`
		Limit    int64                `json:"limit,omitempty,string"`
		InstType constants.InstrumentType `json:"instType"`
	}
	GetOpenInterest struct {
		Uly      string               `json:"uly,omitempty"`
		InstID   string               `json:"instId,omitempty"`
		InstType constants.InstrumentType `json:"instType"`
	}
	GetFundingRate struct {
		InstID string `json:"instId"`
	}
	GetLimitPrice struct {
		InstID string `json:"instId"`
	}
	GetOptionMarketData struct {
		Uly     string `json:"uly"`
		ExpTime string `json:"expTime,omitempty"`
	}
	GetEstimatedDeliveryExercisePrice struct {
		Uly     string `json:"uly"`
		ExpTime string `json:"expTime,omitempty"`
	}
	GetDiscountRateAndInterestFreeQuota struct {
		Uly        string  `json:"uly"`
		Ccy        string  `json:"ccy,omitempty"`
		DiscountLv float64 `json:"discountLv,string"`
	}
	GetLiquidationOrders struct {
		InstID   string               `json:"instId,omitempty"`
		Ccy      string               `json:"ccy,omitempty"`
		Uly      string               `json:"uly,omitempty"`
		After    int64                `json:"after,omitempty,string"`
		Before   int64                `json:"before,omitempty,string"`
		Limit    int64                `json:"limit,omitempty,string"`
		InstType constants.InstrumentType `json:"instType"`
		MgnMode  constants.MarginMode     `json:"mgnMode,omitempty"`
		Alias    constants.AliasType      `json:"alias,omitempty"`
		State    constants.OrderState     `json:"state,omitempty"`
	}
	GetMarkPrice struct {
		InstID   string               `json:"instId,omitempty"`
		Uly      string               `json:"uly,omitempty"`
		InstType constants.InstrumentType `json:"instType"`
	}
	GetPositionTiers struct {
		InstID   string               `json:"instId,omitempty"`
		Uly      string               `json:"uly,omitempty"`
		InstType constants.InstrumentType `json:"instType"`
		TdMode   constants.TradeMode      `json:"tdMode"`
		Tier     constants.JSONInt64      `json:"tier,omitempty"`
	}
	GetUnderlying struct {
		InstType constants.InstrumentType `json:"instType"`
	}
	Status struct {
		State string `json:"state,omitempty"`
	}
)
