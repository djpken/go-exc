// Package types contains common data types used across different exchanges
package types

import (
	"errors"
	"strconv"
	"time"
)

// Common errors
var (
	// ErrNotSupported is returned when a feature is not supported by the exchange
	ErrNotSupported = errors.New("not supported by this exchange")
)

// ZeroDecimal represents a decimal number as a string for precision
const ZeroDecimal Decimal = "0"

type Decimal string

// Float64 converts Decimal to float64
func (d Decimal) Float64() (float64, error) {
	return strconv.ParseFloat(string(d), 64)
}

// String returns the string representation
func (d Decimal) String() string {
	return string(d)
}

// IsZero checks if the decimal value is zero
// Returns true for "0", "0.0", "0.00", empty string, etc.
func (d Decimal) IsZero() bool {
	if d == "" || d == "0" {
		return true
	}

	val, err := d.Float64()
	if err != nil {
		// If parsing fails, treat as non-zero to be safe
		return false
	}

	return val == 0
}

// Add performs addition: d + other
// Returns a new Decimal with the sum
func (d Decimal) Add(other Decimal) (Decimal, error) {
	val1, err := d.Float64()
	if err != nil {
		return ZeroDecimal, err
	}

	val2, err := other.Float64()
	if err != nil {
		return ZeroDecimal, err
	}

	result := val1 + val2
	return Decimal(strconv.FormatFloat(result, 'f', -1, 64)), nil
}

// Sub performs subtraction: d - other
// Returns a new Decimal with the difference
func (d Decimal) Sub(other Decimal) (Decimal, error) {
	val1, err := d.Float64()
	if err != nil {
		return ZeroDecimal, err
	}

	val2, err := other.Float64()
	if err != nil {
		return ZeroDecimal, err
	}

	result := val1 - val2
	return Decimal(strconv.FormatFloat(result, 'f', -1, 64)), nil
}

// Mul performs multiplication: d * other
// Returns a new Decimal with the product
func (d Decimal) Mul(other Decimal) (Decimal, error) {
	val1, err := d.Float64()
	if err != nil {
		return ZeroDecimal, err
	}

	val2, err := other.Float64()
	if err != nil {
		return ZeroDecimal, err
	}

	result := val1 * val2
	return Decimal(strconv.FormatFloat(result, 'f', -1, 64)), nil
}

// Div performs division: d / other
// Returns a new Decimal with the quotient
// Returns error if dividing by zero
func (d Decimal) Div(other Decimal) (Decimal, error) {
	val1, err := d.Float64()
	if err != nil {
		return ZeroDecimal, err
	}

	val2, err := other.Float64()
	if err != nil {
		return ZeroDecimal, err
	}

	if val2 == 0 {
		return ZeroDecimal, errors.New("division by zero")
	}

	result := val1 / val2
	return Decimal(strconv.FormatFloat(result, 'f', -1, 64)), nil
}

// Cmp compares two Decimal values
// Returns:
//   -1 if d < other
//    0 if d == other
//    1 if d > other
func (d Decimal) Cmp(other Decimal) (int, error) {
	val1, err := d.Float64()
	if err != nil {
		return 0, err
	}

	val2, err := other.Float64()
	if err != nil {
		return 0, err
	}

	if val1 < val2 {
		return -1, nil
	} else if val1 > val2 {
		return 1, nil
	}
	return 0, nil
}

// LessThan returns true if d < other
func (d Decimal) LessThan(other Decimal) (bool, error) {
	cmp, err := d.Cmp(other)
	if err != nil {
		return false, err
	}
	return cmp < 0, nil
}

// GreaterThan returns true if d > other
func (d Decimal) GreaterThan(other Decimal) (bool, error) {
	cmp, err := d.Cmp(other)
	if err != nil {
		return false, err
	}
	return cmp > 0, nil
}

// Equal returns true if d == other
func (d Decimal) Equal(other Decimal) (bool, error) {
	cmp, err := d.Cmp(other)
	if err != nil {
		return false, err
	}
	return cmp == 0, nil
}

// Abs returns the absolute value of the decimal
func (d Decimal) Abs() (Decimal, error) {
	val, err := d.Float64()
	if err != nil {
		return ZeroDecimal, err
	}

	if val < 0 {
		val = -val
	}

	return Decimal(strconv.FormatFloat(val, 'f', -1, 64)), nil
}

// Neg returns the negation of the decimal
func (d Decimal) Neg() (Decimal, error) {
	val, err := d.Float64()
	if err != nil {
		return ZeroDecimal, err
	}

	result := -val
	// Handle negative zero case
	if result == 0 {
		result = 0
	}

	return Decimal(strconv.FormatFloat(result, 'f', -1, 64)), nil
}

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
	Symbol   string  // Trading symbol (e.g., "BTC-USDT")
	Side     string  // Order side: "buy" or "sell" (trading direction)
	PosSide  string  // Position side: "long" or "short" (for futures/derivatives, empty for spot)
	Type     string  // Order type: "limit", "market", etc.
	Quantity float64 // Order quantity
	Price    float64 // Order price (for limit orders)
	Extra    map[string]interface{} // Exchange-specific parameters
}

// CancelOrderRequest contains parameters for canceling an order
type CancelOrderRequest struct {
	Symbol  string
	OrderID string
	Extra   map[string]interface{}
}

// GetOrderRequest contains parameters for getting order details
type GetOrderRequest struct {
	Symbol  string
	OrderID string
	Extra   map[string]interface{}
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
