package funding

import "github.com/djpken/go-exc/exchanges/okex/constants"

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
		Amt     constants.JSONFloat64 `json:"amt"`
		From    constants.AccountType `json:"from,string"`
		To      constants.AccountType `json:"to,string"`
	}
	Bill struct {
		BillID string            `json:"billId"`
		Ccy    string            `json:"ccy"`
		Bal    constants.JSONFloat64 `json:"bal"`
		BalChg constants.JSONFloat64 `json:"balChg"`
		Type   constants.BillType    `json:"type,string"`
		TS     constants.JSONTime    `json:"ts"`
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
		To       constants.AccountType `json:"to,string"`
		TS       constants.JSONTime    `json:"ts"`
	}
	DepositHistory struct {
		Ccy   string             `json:"ccy"`
		Chain string             `json:"chain"`
		TxID  string             `json:"txId"`
		From  string             `json:"from"`
		To    string             `json:"to"`
		DepId string             `json:"depId"`
		Amt   constants.JSONFloat64  `json:"amt"`
		State constants.DepositState `json:"state,string"`
		TS    constants.JSONTime     `json:"ts"`
	}
	Withdrawal struct {
		Ccy   string            `json:"ccy"`
		Chain string            `json:"chain"`
		WdID  constants.JSONInt64   `json:"wdId"`
		Amt   constants.JSONFloat64 `json:"amt"`
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
		Amt   constants.JSONFloat64     `json:"amt"`
		Fee   constants.JSONFloat64     `json:"fee"`
		WdID  constants.JSONInt64       `json:"wdId"`
		State constants.WithdrawalState `json:"state,string"`
		TS    constants.JSONTime        `json:"ts"`
	}
	PiggyBank struct {
		Ccy  string            `json:"ccy"`
		Amt  constants.JSONFloat64 `json:"amt"`
		Side constants.ActionType  `json:"side,string"`
	}
	PiggyBankBalance struct {
		Ccy      string            `json:"ccy"`
		Amt      constants.JSONFloat64 `json:"amt"`
		Earnings constants.JSONFloat64 `json:"earnings"`
	}
)
