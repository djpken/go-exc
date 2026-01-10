package tradedata

import "github.com/djpken/go-exc/exchanges/okex/types"

type (
	GetTakerVolume struct {
		Ccy      string               `json:"ccy"`
		Begin    int64                `json:"before,omitempty,string"`
		End      int64                `json:"limit,omitempty,string"`
		InstType types.InstrumentType `json:"instType"`
		Period   types.BarSize        `json:"period,string,omitempty"`
	}
	GetRatio struct {
		Ccy    string        `json:"ccy"`
		Begin  int64         `json:"before,omitempty,string"`
		End    int64         `json:"limit,omitempty,string"`
		Period types.BarSize `json:"period,string,omitempty"`
	}
	GetOpenInterestAndVolumeStrike struct {
		Ccy     string        `json:"ccy"`
		ExpTime string        `json:"expTime"`
		Period  types.BarSize `json:"period,string,omitempty"`
	}
)
