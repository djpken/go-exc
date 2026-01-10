package private

import "github.com/djpken/go-exc/exchanges/okex/types"

type (
	Account struct {
		Ccy string `json:"ccy,omitempty"`
	}
	Position struct {
		Uly      string               `json:"uly,omitempty"`
		InstID   string               `json:"instId,omitempty"`
		InstType types.InstrumentType `json:"instType"`
	}
	Order struct {
		Uly      string               `json:"uly,omitempty"`
		InstID   string               `json:"instId,omitempty"`
		InstType types.InstrumentType `json:"instType"`
	}
	AlgoOrder struct {
		Uly      string               `json:"uly,omitempty"`
		InstID   string               `json:"instId,omitempty"`
		InstType types.InstrumentType `json:"instType"`
	}
)
