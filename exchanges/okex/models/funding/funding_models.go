package funding

import "github.com/djpken/go-exc/exchanges/okex/types"

type (
	Currency struct {
		Ccy         string `json:"ccy"`
		Name        string `json:"name"`
		Chain       string `json:"chain"`
		MinWd       string `json:"minWd"`
		MinFee      string `json:"minFee"`
		MaxFee      string `json:"maxFee"`
		CanDep      bool   `json:"canDep"`
		CanWd       bool   `json:"canWd"`
		CanInternal bool   `json:"canInternal"`
	}
	Balance struct {
		Ccy       string `json:"ccy"`
		Bal       string `json:"bal"`
		FrozenBal string `json:"frozenBal"`
		AvailBal  string `json:"availBal"`
	}
	Transfer struct {
		TransID string            `json:"transId"`
		Ccy     string            `json:"ccy"`
		Amt     types.JSONFloat64 `json:"amt"`
		From    types.AccountType `json:"from,string"`
		To      types.AccountType `json:"to,string"`
	}
	Bill struct {
		BillID string            `json:"billId"`
		Ccy    string            `json:"ccy"`
		Bal    types.JSONFloat64 `json:"bal"`
		BalChg types.JSONFloat64 `json:"balChg"`
		Type   types.BillType    `json:"type,string"`
		TS     types.JSONTime    `json:"ts"`
	}
	DepositAddress struct {
		Addr     string            `json:"addr"`
		Tag      string            `json:"tag,omitempty"`
		Memo     string            `json:"memo,omitempty"`
		PmtID    string            `json:"pmtId,omitempty"`
		Ccy      string            `json:"ccy"`
		Chain    string            `json:"chain"`
		CtAddr   string            `json:"ctAddr"`
		Selected bool              `json:"selected"`
		To       types.AccountType `json:"to,string"`
		TS       types.JSONTime    `json:"ts"`
	}
	DepositHistory struct {
		Ccy   string             `json:"ccy"`
		Chain string             `json:"chain"`
		TxID  string             `json:"txId"`
		From  string             `json:"from"`
		To    string             `json:"to"`
		DepId string             `json:"depId"`
		Amt   types.JSONFloat64  `json:"amt"`
		State types.DepositState `json:"state,string"`
		TS    types.JSONTime     `json:"ts"`
	}
	Withdrawal struct {
		Ccy   string            `json:"ccy"`
		Chain string            `json:"chain"`
		WdID  types.JSONInt64   `json:"wdId"`
		Amt   types.JSONFloat64 `json:"amt"`
	}
	WithdrawalHistory struct {
		Ccy   string                `json:"ccy"`
		Chain string                `json:"chain"`
		TxID  string                `json:"txId"`
		From  string                `json:"from"`
		To    string                `json:"to"`
		Tag   string                `json:"tag,omitempty"`
		PmtID string                `json:"pmtId,omitempty"`
		Memo  string                `json:"memo,omitempty"`
		Amt   types.JSONFloat64     `json:"amt"`
		Fee   types.JSONFloat64     `json:"fee"`
		WdID  types.JSONInt64       `json:"wdId"`
		State types.WithdrawalState `json:"state,string"`
		TS    types.JSONTime        `json:"ts"`
	}
	PiggyBank struct {
		Ccy  string            `json:"ccy"`
		Amt  types.JSONFloat64 `json:"amt"`
		Side types.ActionType  `json:"side,string"`
	}
	PiggyBankBalance struct {
		Ccy      string            `json:"ccy"`
		Amt      types.JSONFloat64 `json:"amt"`
		Earnings types.JSONFloat64 `json:"earnings"`
	}
)
