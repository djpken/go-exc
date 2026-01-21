package okex

import "github.com/djpken/go-exc/exchanges/okex/constants"

// Re-export types from types package for backward compatibility
type (
	BaseURL              = constants.BaseURL
	InstrumentType       = constants.InstrumentType
	MarginMode           = constants.MarginMode
	ContractType         = constants.ContractType
	PosModeType          = constants.PosModeType
	PositionSide         = constants.PositionSide
	ActualSide           = constants.ActualSide
	TradeMode            = constants.TradeMode
	CountAction          = constants.CountAction
	OrderSide            = constants.OrderSide
	GreekType            = constants.GreekType
	BarSize              = constants.BarSize
	TradeSide            = constants.TradeSide
	ChannelName          = constants.ChannelName
	Operation            = constants.Operation
	EventType            = constants.EventType
	OrderType            = constants.OrderType
	AlgoOrderType        = constants.AlgoOrderType
	QuantityType         = constants.QuantityType
	OrderFlowType        = constants.OrderFlowType
	OrderState           = constants.OrderState
	ActionType           = constants.ActionType
	APIKeyAccess         = constants.APIKeyAccess
	OptionType           = constants.OptionType
	AliasType            = constants.AliasType
	InstrumentState      = constants.InstrumentState
	DeliveryExerciseType = constants.DeliveryExerciseType
	CandleStickWsBarSize = constants.CandleStickWsBarSize

	Destination           = constants.Destination
	BillType              = constants.BillType
	BillSubType           = constants.BillSubType
	FeeCategory           = constants.FeeCategory
	TransferType          = constants.TransferType
	AccountType           = constants.AccountType
	DepositState          = constants.DepositState
	WithdrawalDestination = constants.WithdrawalDestination
	WithdrawalState       = constants.WithdrawalState

	JSONFloat64 = constants.JSONFloat64
	JSONInt64   = constants.JSONInt64
	JSONTime    = constants.JSONTime

	ClientError = constants.ClientError
)

// Re-export constants
const (
	NormalServer = constants.NormalServer
	AwsServer    = constants.AwsServer
	DemoServer   = constants.DemoServer

	RestURL      = constants.RestURL
	PublicWsURL  = constants.PublicWsURL
	PrivateWsURL = constants.PrivateWsURL

	AwsRestURL      = constants.AwsRestURL
	AwsPublicWsURL  = constants.AwsPublicWsURL
	AwsPrivateWsURL = constants.AwsPrivateWsURL

	DemoRestURL      = constants.DemoRestURL
	DemoPublicWsURL  = constants.DemoPublicWsURL
	DemoPrivateWsURL = constants.DemoPrivateWsURL

	AnyInstrument     = constants.AnyInstrument
	SpotInstrument    = constants.SpotInstrument
	MarginInstrument  = constants.MarginInstrument
	SwapInstrument    = constants.SwapInstrument
	FuturesInstrument = constants.FuturesInstrument
	OptionsInstrument = constants.OptionsInstrument

	MarginCrossMode    = constants.MarginCrossMode
	MarginIsolatedMode = constants.MarginIsolatedMode

	ContractLinearType  = constants.ContractLinearType
	ContractInverseType = constants.ContractInverseType
)
