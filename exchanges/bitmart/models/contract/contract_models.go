package contract

// ContractDetail represents detailed information about a contract
type ContractDetail struct {
	// Symbol is the trading pair name
	Symbol string `json:"symbol"`

	// ProductType is the contract type: 1=Perpetual, 2=Futures (currently only perpetual)
	ProductType int `json:"product_type"`

	// OpenTimestamp is the first opening time
	OpenTimestamp int64 `json:"open_timestamp"`

	// ExpireTimestamp is the expiration date, 0 if never expires
	ExpireTimestamp int64 `json:"expire_timestamp"`

	// SettleTimestamp is the settlement date, 0 if no automatic settlement
	SettleTimestamp int64 `json:"settle_timestamp"`

	// BaseCurrency is the contract base currency
	BaseCurrency string `json:"base_currency"`

	// QuoteCurrency is the contract quote currency
	QuoteCurrency string `json:"quote_currency"`

	// LastPrice is the latest trade price
	LastPrice string `json:"last_price"`

	// Volume24h is the 24h trading volume
	Volume24h string `json:"volume_24h"`

	// Turnover24h is the 24h turnover
	Turnover24h string `json:"turnover_24h"`

	// IndexPrice is the index price
	IndexPrice string `json:"index_price"`

	// IndexName is the index name
	IndexName string `json:"index_name"`

	// ContractSize is the contract size
	ContractSize string `json:"contract_size"`

	// MinLeverage is the minimum leverage multiplier
	MinLeverage string `json:"min_leverage"`

	// MaxLeverage is the maximum leverage multiplier
	MaxLeverage string `json:"max_leverage"`

	// PricePrecision is the price precision
	PricePrecision string `json:"price_precision"`

	// VolPrecision is the volume precision
	VolPrecision string `json:"vol_precision"`

	// MaxVolume is the maximum volume for limit orders
	MaxVolume string `json:"max_volume"`

	// MarketMaxVolume is the maximum volume for market orders
	MarketMaxVolume string `json:"market_max_volume"`

	// MinVolume is the minimum order volume
	MinVolume string `json:"min_volume"`

	// FundingRate is the previous settlement funding rate
	FundingRate string `json:"funding_rate"`

	// ExpectedFundingRate is the expected settlement funding rate
	ExpectedFundingRate string `json:"expected_funding_rate"`

	// OpenInterest is the open interest (position volume)
	OpenInterest string `json:"open_interest"`

	// OpenInterestValue is the open interest value
	OpenInterestValue string `json:"open_interest_value"`

	// High24h is the 24h highest price
	High24h string `json:"high_24h"`

	// Low24h is the 24h lowest price
	Low24h string `json:"low_24h"`

	// Change24h is the 24h price change
	Change24h string `json:"change_24h"`

	// FundingIntervalHours is the funding rate collection interval in hours
	FundingIntervalHours int `json:"funding_interval_hours"`

	// Status is the trading status: Trading=Trading, Delisted=Delisted
	Status string `json:"status"`

	// DelistTime is the delist time (if status=Trading, it's the estimated delist time)
	DelistTime int64 `json:"delist_time"`
}
