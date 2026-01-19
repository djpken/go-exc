package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// GenerateSignature generates HMAC-SHA256 signature for BitMart API
func GenerateSignature(timestamp, body, secretKey string) string {
	message := timestamp + "#" + "bitmart.WebSocket" + "#" + body
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))
	return hex.EncodeToString(mac.Sum(nil))
}

// GetTimestamp returns current timestamp in milliseconds
func GetTimestamp() string {
	return strconv.FormatInt(time.Now().UnixMilli(), 10)
}

// FormatFloat formats float64 to string with given precision
func FormatFloat(f float64, precision int) string {
	return fmt.Sprintf("%.*f", precision, f)
}

// StringToFloat converts string to float64
func StringToFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}
