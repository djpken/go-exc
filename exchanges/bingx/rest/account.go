package rest

import "fmt"

// Account provides BingX account endpoints
type Account struct {
	client *ClientRest
}

func NewAccount(c *ClientRest) *Account { return &Account{client: c} }

// BalanceAsset represents a single asset in the account balance
type BalanceAsset struct {
	Asset            string `json:"asset"`
	Balance          string `json:"balance"`
	Equity           string `json:"equity"`
	UnrealizedProfit string `json:"unrealizedProfit"`
	RealisedProfit   string `json:"realisedProfit"`
	AvailableMargin  string `json:"availableMargin"`
	UsedMargin       string `json:"usedMargin"`
	FreezedMargin    string `json:"freezedMargin"`
}

// BalanceResponse is the full API response for account balance
type BalanceResponse struct {
	Code int          `json:"code"`
	Data BalanceAsset `json:"data"`
}

// GetBalance retrieves the perpetual swap account balance
// GET /openApi/swap/v2/user/balance
func (a *Account) GetBalance(currency string) (*BalanceResponse, error) {
	var result BalanceResponse
	params := map[string]string{}
	if currency != "" {
		params["currency"] = currency
	}
	if err := a.client.GET("/openApi/swap/v2/user/balance", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PositionData holds details for a single position
type PositionData struct {
	Symbol           string `json:"symbol"`
	PositionID       string `json:"positionId"`
	PositionSide     string `json:"positionSide"`
	Isolated         bool   `json:"isolated"`
	PositionAmt      string `json:"positionAmt"`
	AvailableAmt     string `json:"availableAmt"`
	UnrealizedProfit string `json:"unrealizedProfit"`
	RealisedProfit   string `json:"realisedProfit"`
	InitialMargin    string `json:"initialMargin"`
	AvgPrice         string `json:"avgPrice"`
	Leverage         int    `json:"leverage"`
}

// PositionsResponse is the full API response for positions
type PositionsResponse struct {
	Code int            `json:"code"`
	Data []PositionData `json:"data"`
}

// GetPositions retrieves all open positions (or a specific symbol if provided)
// GET /openApi/swap/v2/user/positions
func (a *Account) GetPositions(symbol string) (*PositionsResponse, error) {
	var result PositionsResponse
	params := map[string]string{}
	if symbol != "" {
		params["symbol"] = symbol
	}
	if err := a.client.GET("/openApi/swap/v2/user/positions", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// LeverageData holds leverage information for a symbol
type LeverageData struct {
	LongLeverage   int64 `json:"longLeverage"`
	ShortLeverage  int64 `json:"shortLeverage"`
	MaxLongLeverage  int64 `json:"maxLongLeverage"`
	MaxShortLeverage int64 `json:"maxShortLeverage"`
}

// LeverageResponse is the full API response for leverage
type LeverageResponse struct {
	Code int          `json:"code"`
	Data LeverageData `json:"data"`
}

// GetLeverage queries the current leverage for a symbol
// GET /openApi/swap/v2/trade/leverage
func (a *Account) GetLeverage(symbol string) (*LeverageResponse, error) {
	var result LeverageResponse
	params := map[string]string{"symbol": symbol}
	if err := a.client.GET("/openApi/swap/v2/trade/leverage", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SetLeverageData holds the result after setting leverage
type SetLeverageData struct {
	Symbol        string `json:"symbol"`
	Side          string `json:"side"`
	Leverage      int    `json:"leverage"`
	MaxNotional   string `json:"maxNotionalValue"`
}

// SetLeverageResponse is the full API response for setting leverage
type SetLeverageResponse struct {
	Code int             `json:"code"`
	Data SetLeverageData `json:"data"`
}

// SetLeverage sets the leverage for a symbol and side
// POST /openApi/swap/v2/trade/leverage
func (a *Account) SetLeverage(symbol, side string, leverage int) (*SetLeverageResponse, error) {
	var result SetLeverageResponse
	params := map[string]string{
		"symbol":   symbol,
		"side":     side,
		"leverage": fmt.Sprintf("%d", leverage),
	}
	if err := a.client.POST("/openApi/swap/v2/trade/leverage", params, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
