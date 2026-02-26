package rest

// ListenKeyData holds the listen key returned by the BingX API.
type ListenKeyData struct {
	ListenKey string `json:"listenKey"`
}

// CreateListenKey creates a new listen key for private WebSocket connections.
// The key is valid for 1 hour; call ExtendListenKey every 30 minutes to keep it alive.
func (c *ClientRest) CreateListenKey() (string, error) {
	var resp Response[ListenKeyData]
	if err := c.POST("/openApi/user/auth/userDataStream", nil, &resp); err != nil {
		return "", err
	}
	return resp.Data.ListenKey, nil
}

// ExtendListenKey resets the 60-minute expiry timer for an existing listen key.
func (c *ClientRest) ExtendListenKey(listenKey string) error {
	return c.PUT("/openApi/user/auth/userDataStream", map[string]string{
		"listenKey": listenKey,
	}, nil)
}

// DeleteListenKey invalidates a listen key and closes the associated private stream.
func (c *ClientRest) DeleteListenKey(listenKey string) error {
	return c.DELETE("/openApi/user/auth/userDataStream", map[string]string{
		"listenKey": listenKey,
	}, nil)
}
