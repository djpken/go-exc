// Package types contains common data types used across different exchanges
package types

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

// Common errors
var (
	// ErrNotSupported is returned when a feature is not supported by the exchange
	ErrNotSupported = errors.New("not supported by this exchange")
)

// ZeroDecimal represents a zero value for Decimal type
var ZeroDecimal = Decimal{decimal.Zero}

// Decimal represents a high-precision decimal number
// Built on shopspring/decimal for financial accuracy
type Decimal struct {
	decimal.Decimal
}

// NewDecimal creates a Decimal from a string
// Returns error if the string cannot be parsed
func NewDecimal(value string) (Decimal, error) {
	d, err := decimal.NewFromString(value)
	if err != nil {
		return ZeroDecimal, err
	}
	return Decimal{d}, nil
}

// NewDecimalFromFloat creates a Decimal from a float64
// Note: May lose precision due to float64 representation
func NewDecimalFromFloat(value float64) Decimal {
	return Decimal{decimal.NewFromFloat(value)}
}

// NewDecimalFromInt creates a Decimal from an int64
func NewDecimalFromInt(value int64) Decimal {
	return Decimal{decimal.NewFromInt(value)}
}

// MustDecimal creates a Decimal from a string, panics on error
// Use this only when you're certain the string is valid
func MustDecimal(value string) Decimal {
	d, err := decimal.NewFromString(value)
	if err != nil {
		panic(err)
	}
	return Decimal{d}
}

// ========== Arithmetic Operations ==========

// Add performs addition: d + other
func (d Decimal) Add(other Decimal) (Decimal, error) {
	return Decimal{d.Decimal.Add(other.Decimal)}, nil
}

// Sub performs subtraction: d - other
func (d Decimal) Sub(other Decimal) (Decimal, error) {
	return Decimal{d.Decimal.Sub(other.Decimal)}, nil
}

// Mul performs multiplication: d * other
func (d Decimal) Mul(other Decimal) (Decimal, error) {
	return Decimal{d.Decimal.Mul(other.Decimal)}, nil
}

// Div performs division: d / other
// Returns error if dividing by zero
func (d Decimal) Div(other Decimal) (Decimal, error) {
	if other.IsZero() {
		return ZeroDecimal, errors.New("division by zero")
	}
	return Decimal{d.Decimal.Div(other.Decimal)}, nil
}

// Mod returns the remainder: d % other
func (d Decimal) Mod(other Decimal) (Decimal, error) {
	if other.IsZero() {
		return ZeroDecimal, errors.New("modulo by zero")
	}
	return Decimal{d.Decimal.Mod(other.Decimal)}, nil
}

// Pow returns d raised to the power of exp: d^exp
func (d Decimal) Pow(exp Decimal) Decimal {
	return Decimal{d.Decimal.Pow(exp.Decimal)}
}

// ========== Comparison Operations ==========

// Cmp compares two Decimal values
// Returns:
//
//	-1 if d < other
//	 0 if d == other
//	 1 if d > other
func (d Decimal) Cmp(other Decimal) (int, error) {
	return d.Decimal.Cmp(other.Decimal), nil
}

// LessThan returns true if d < other
func (d Decimal) LessThan(other Decimal) (bool, error) {
	return d.Decimal.LessThan(other.Decimal), nil
}

// LessThanOrEqual returns true if d <= other
func (d Decimal) LessThanOrEqual(other Decimal) (bool, error) {
	return d.Decimal.LessThanOrEqual(other.Decimal), nil
}

// GreaterThan returns true if d > other
func (d Decimal) GreaterThan(other Decimal) (bool, error) {
	return d.Decimal.GreaterThan(other.Decimal), nil
}

// GreaterThanOrEqual returns true if d >= other
func (d Decimal) GreaterThanOrEqual(other Decimal) (bool, error) {
	return d.Decimal.GreaterThanOrEqual(other.Decimal), nil
}

// Equal returns true if d == other
func (d Decimal) Equal(other Decimal) (bool, error) {
	return d.Decimal.Equal(other.Decimal), nil
}

// ========== Sign and Value Checks ==========

// IsZero checks if the decimal value is zero
func (d Decimal) IsZero() bool {
	return d.Decimal.IsZero()
}

// IsPositive returns true if the decimal value is greater than zero
func (d Decimal) IsPositive() bool {
	return d.Decimal.IsPositive()
}

// IsNegative returns true if the decimal value is less than zero
func (d Decimal) IsNegative() bool {
	return d.Decimal.IsNegative()
}

// Sign returns:
//
//	-1 if d < 0
//	 0 if d == 0
//	 1 if d > 0
func (d Decimal) Sign() int {
	return d.Decimal.Sign()
}

// ========== Absolute Value and Negation ==========

// Abs returns the absolute value of the decimal
func (d Decimal) Abs() (Decimal, error) {
	return Decimal{d.Decimal.Abs()}, nil
}

// Neg returns the negation of the decimal
func (d Decimal) Neg() (Decimal, error) {
	return Decimal{d.Decimal.Neg()}, nil
}

// ========== Min/Max Operations ==========

// Max returns the maximum of two Decimal values
func (d Decimal) Max(other Decimal) (Decimal, error) {
	if d.Decimal.GreaterThanOrEqual(other.Decimal) {
		return d, nil
	}
	return other, nil
}

// Min returns the minimum of two Decimal values
func (d Decimal) Min(other Decimal) (Decimal, error) {
	if d.Decimal.LessThanOrEqual(other.Decimal) {
		return d, nil
	}
	return other, nil
}

// ========== Rounding Operations ==========

// Round rounds to the nearest integer at the specified number of decimal places
// Uses "round half up" strategy (banker's rounding)
// Example: Decimal("2.5").Round(0) => "3"
func (d Decimal) Round(places int32) (Decimal, error) {
	return Decimal{d.Decimal.Round(places)}, nil
}

// RoundBank rounds to the nearest integer at the specified number of decimal places
// Uses "round half to even" strategy (banker's rounding)
// This is the recommended method for financial calculations
// Example: Decimal("2.5").RoundBank(0) => "2", Decimal("3.5").RoundBank(0) => "4"
func (d Decimal) RoundBank(places int32) (Decimal, error) {
	return Decimal{d.Decimal.RoundBank(places)}, nil
}

// RoundCash rounds to the nearest cent (or smallest currency unit)
// interval: the smallest currency unit (e.g., 0.05 for nickels)
func (d Decimal) RoundCash(interval uint8) (Decimal, error) {
	return Decimal{d.Decimal.RoundCash(interval)}, nil
}

// RoundUp rounds up (away from zero) to the specified number of decimal places
// Example: Decimal("2.4").RoundUp(0) => "3", Decimal("-2.4").RoundUp(0) => "-3"
func (d Decimal) RoundUp(places int32) (Decimal, error) {
	return Decimal{d.Decimal.RoundUp(places)}, nil
}

// RoundDown rounds down (towards zero) to the specified number of decimal places
// Example: Decimal("2.6").RoundDown(0) => "2", Decimal("-2.6").RoundDown(0) => "-2"
func (d Decimal) RoundDown(places int32) (Decimal, error) {
	return Decimal{d.Decimal.RoundDown(places)}, nil
}

// Truncate truncates to the specified number of decimal places (no rounding)
// Example: Decimal("2.6822000000000004").Truncate(4) => "2.6822"
func (d Decimal) Truncate(places int32) (Decimal, error) {
	return Decimal{d.Decimal.Truncate(places)}, nil
}

// Floor returns the largest integer value less than or equal to d
func (d Decimal) Floor() Decimal {
	return Decimal{d.Decimal.Floor()}
}

// Ceil returns the smallest integer value greater than or equal to d
func (d Decimal) Ceil() Decimal {
	return Decimal{d.Decimal.Ceil()}
}

// ========== Conversion Methods ==========

// String returns the string representation of the decimal
// This will return the exact decimal value without trailing zeros
func (d Decimal) String() string {
	return d.Decimal.String()
}

// GoString implements fmt.GoStringer for better debugging experience
// This method is called when using %#v format or in debuggers
func (d Decimal) GoString() string {
	return d.Decimal.String()
}

// StringFixed returns a string representation with fixed decimal places
// Example: Decimal("2.5").StringFixed(2) => "2.50"
func (d Decimal) StringFixed(places int32) string {
	return d.Decimal.StringFixed(places)
}

// Float64 converts Decimal to float64
// Note: May lose precision for very large or very small numbers
func (d Decimal) Float64() (float64, error) {
	f, _ := d.Decimal.Float64()
	return f, nil
}

// F64 is a convenience method that returns float64 without error
// Useful for quick debugging: just check d.F64() in watch window
// Note: May lose precision for very large or very small numbers
func (d Decimal) F64() float64 {
	f, _ := d.Decimal.Float64()
	return f
}

// DebugString returns a detailed string representation for debugging
// Format: "value (coefficient: X, exponent: Y)"
func (d Decimal) DebugString() string {
	return d.Decimal.String()
}

// Int64 converts Decimal to int64
// Returns error if the number has fractional parts or is out of range
func (d Decimal) Int64() (int64, error) {
	if !d.Decimal.IsInteger() {
		return 0, errors.New("decimal has fractional part")
	}
	return d.Decimal.IntPart(), nil
}

// IntPart returns the integer part of the decimal
// Truncates any fractional part
func (d Decimal) IntPart() int64 {
	return d.Decimal.IntPart()
}

// Coefficient returns the coefficient of the decimal
// For internal use and advanced operations
func (d Decimal) Coefficient() *int {
	coef := d.Decimal.Coefficient().Int64()
	intCoef := int(coef)
	return &intCoef
}

// Exponent returns the exponent of the decimal
// For internal use and advanced operations
func (d Decimal) Exponent() int32 {
	return d.Decimal.Exponent()
}

// ========== Helper Functions ==========

// IsInteger checks if the decimal is an integer (no fractional part)
func (d Decimal) IsInteger() bool {
	return d.Decimal.IsInteger()
}

// ========== Timestamp Type ==========

// Timestamp represents a timestamp
type Timestamp time.Time

// Time returns the time.Time representation
func (t Timestamp) Time() time.Time {
	return time.Time(t)
}

// Unix returns the Unix timestamp in seconds
func (t Timestamp) Unix() int64 {
	return time.Time(t).Unix()
}

// UnixMilli returns the Unix timestamp in milliseconds
func (t Timestamp) UnixMilli() int64 {
	return time.Time(t).UnixNano() / 1e6
}

// ========== Request Types ==========
// 以下类型用于统一的 API 请求

// PlaceOrderRequest contains parameters for placing an order
type PlaceOrderRequest struct {
	Symbol        string       // Trading symbol (e.g., "BTC-USDT")
	Side          OrderSide    // Order side: "buy" or "sell" (trading direction)
	PosSide       PositionSide // Position side: "long" or "short" (for futures/derivatives, empty for spot)
	TdMode        MarginMode
	Type          string  // Order type: "limit", "market", etc.
	Quantity      float64 // Order quantity
	Price         float64 // Order price (for limit orders)
	ClientOrderID string
	Extra         map[string]interface{} // Exchange-specific parameters
}

// CancelOrderRequest contains parameters for canceling an order
type CancelOrderRequest struct {
	Symbol  string
	OrderID string
	Extra   map[string]interface{}
}

// GetOrderRequest contains parameters for getting order details
type GetOrderRequest struct {
	Symbol        string
	OrderID       string
	ClientOrderID string
	Extra         map[string]interface{}
}

// WithdrawRequest contains parameters for withdrawal
type WithdrawRequest struct {
	Currency string
	Amount   float64
	Address  string
	Tag      string
	Extra    map[string]interface{}
}

// AccountConfig represents account configuration settings
type AccountConfig struct {
	// UID is the user ID
	UID string

	// Level is the account level
	Level string

	// AutoLoan indicates if auto-loan is enabled (for margin trading)
	AutoLoan bool

	// PositionMode is the position mode (long/short mode or net mode)
	// Values: "long_short_mode" or "net_mode"
	PositionMode string

	// Extra contains exchange-specific fields
	Extra map[string]interface{}
}

// GetInstrumentsRequest contains parameters for querying trading instruments
type GetInstrumentsRequest struct {
	// InstrumentType filters by instrument type (spot, futures, swap, etc.)
	// Optional: empty string returns all types
	InstrumentType InstrumentType

	// Extra contains exchange-specific parameters
	Extra map[string]interface{}
}

// GetTickersRequest contains parameters for querying multiple tickers
type GetTickersRequest struct {
	// InstrumentType filters by instrument type (spot, futures, swap, etc.)
	// Optional: empty string returns all types
	InstrumentType InstrumentType

	// Extra contains exchange-specific parameters
	Extra map[string]interface{}
}

// GetCandlesRequest contains parameters for querying candlestick/kline data
type GetCandlesRequest struct {
	// Symbol is the trading symbol
	Symbol string

	// Interval is the candlestick interval (e.g., "1m", "5m", "1H", "1D")
	Interval string

	// Limit is the maximum number of candles to return
	// Optional: defaults to exchange-specific default (usually 100-500)
	Limit int

	// StartTime is the start time for the query
	// Optional: if not specified, returns the most recent candles
	StartTime *time.Time

	// EndTime is the end time for the query
	// Optional: if not specified, returns up to current time
	EndTime *time.Time

	// Extra contains exchange-specific parameters
	Extra map[string]interface{}
}
