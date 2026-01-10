package okex

import "github.com/djpken/go-exc/exchanges/okex/types"

// Re-export types from types package for backward compatibility
type (
	BaseURL              = types.BaseURL
	InstrumentType       = types.InstrumentType
	MarginMode           = types.MarginMode
	ContractType         = types.ContractType
	PosModeType          = types.PosModeType
	PositionSide         = types.PositionSide
	ActualSide           = types.ActualSide
	TradeMode            = types.TradeMode
	CountAction          = types.CountAction
	OrderSide            = types.OrderSide
	GreekType            = types.GreekType
	BarSize              = types.BarSize
	TradeSide            = types.TradeSide
	ChannelName          = types.ChannelName
	Operation            = types.Operation
	EventType            = types.EventType
	OrderType            = types.OrderType
	AlgoOrderType        = types.AlgoOrderType
	QuantityType         = types.QuantityType
	OrderFlowType        = types.OrderFlowType
	OrderState           = types.OrderState
	ActionType           = types.ActionType
	APIKeyAccess         = types.APIKeyAccess
	OptionType           = types.OptionType
	AliasType            = types.AliasType
	InstrumentState      = types.InstrumentState
	DeliveryExerciseType = types.DeliveryExerciseType
	CandleStickWsBarSize = types.CandleStickWsBarSize

	Destination           = types.Destination
	BillType              = types.BillType
	BillSubType           = types.BillSubType
	FeeCategory           = types.FeeCategory
	TransferType          = types.TransferType
	AccountType           = types.AccountType
	DepositState          = types.DepositState
	WithdrawalDestination = types.WithdrawalDestination
	WithdrawalState       = types.WithdrawalState

	JSONFloat64 = types.JSONFloat64
	JSONInt64   = types.JSONInt64
	JSONTime    = types.JSONTime

	ClientError = types.ClientError
)

// Re-export constants
const (
	NormalServer = types.NormalServer
	AwsServer    = types.AwsServer
	DemoServer   = types.DemoServer

	RestURL      = types.RestURL
	PublicWsURL  = types.PublicWsURL
	PrivateWsURL = types.PrivateWsURL

	AwsRestURL      = types.AwsRestURL
	AwsPublicWsURL  = types.AwsPublicWsURL
	AwsPrivateWsURL = types.AwsPrivateWsURL

	DemoRestURL      = types.DemoRestURL
	DemoPublicWsURL  = types.DemoPublicWsURL
	DemoPrivateWsURL = types.DemoPrivateWsURL

	AnyInstrument     = types.AnyInstrument
	SpotInstrument    = types.SpotInstrument
	MarginInstrument  = types.MarginInstrument
	SwapInstrument    = types.SwapInstrument
	FuturesInstrument = types.FuturesInstrument
	OptionsInstrument = types.OptionsInstrument

	MarginCrossMode    = types.MarginCrossMode
	MarginIsolatedMode = types.MarginIsolatedMode

	ContractLinearType  = types.ContractLinearType
	ContractInverseType = types.ContractInverseType
)
