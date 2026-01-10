// Package types contains common data types used across different exchanges
package types

import (
	"strconv"
	"time"
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
