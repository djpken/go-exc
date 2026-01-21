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

// Decimal represents a decimal number as a string for precision
type Decimal string

// Float64 converts Decimal to float64
func (d Decimal) Float64() (float64, error) {
	return strconv.ParseFloat(string(d), 64)
}

// String returns the string representation
func (d Decimal) String() string {
	return string(d)
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
	Symbol   string
	Side     string
	Type     string
	Quantity float64
	Price    float64
	Extra    map[string]interface{}
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
