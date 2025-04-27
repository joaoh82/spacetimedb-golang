package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// DatabaseInfo represents information about a database
type DatabaseInfo struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

// ConnectWebSocket establishes a WebSocket connection to a database
func (c *Client) ConnectWebSocket(dbAddress string) error {
	if c.token == "" {
		return fmt.Errorf("token is required for WebSocket connection")
	}

	wsURL := url.URL{
		Scheme:   "ws",
		Host:     c.baseURL,
		Path:     fmt.Sprintf("/v1/database/%s/subscribe", dbAddress),
		RawQuery: "token=" + c.token,
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
	}

	conn, _, err := dialer.Dial(wsURL.String(), nil)
	if err != nil {
		return fmt.Errorf("error connecting to WebSocket: %w", err)
	}

	c.wsClient = conn
	return nil
}

// SendMessage sends a message through the WebSocket connection
func (c *Client) SendMessage(message interface{}) error {
	if c.wsClient == nil {
		return fmt.Errorf("WebSocket connection not established")
	}

	return c.wsClient.WriteJSON(message)
}

// ReceiveMessage receives a message from the WebSocket connection
func (c *Client) ReceiveMessage() (interface{}, error) {
	if c.wsClient == nil {
		return nil, fmt.Errorf("WebSocket connection not established")
	}

	var message interface{}
	err := c.wsClient.ReadJSON(&message)
	if err != nil {
		return nil, fmt.Errorf("error reading message: %w", err)
	}

	return message, nil
}

// GetDatabaseInfo retrieves information about a database
func (c *Client) GetDatabaseInfo(dbAddress string) (*DatabaseInfo, error) {
	url := fmt.Sprintf("%s/v1/database/%s", c.baseURL, dbAddress)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", c.getAuthHeader())
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var dbInfo DatabaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&dbInfo); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &dbInfo, nil
}

// ExecuteQuery executes a SQL query against the database
func (c *Client) ExecuteQuery(dbAddress string, query string) (interface{}, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token is required for executing queries")
	}

	url := fmt.Sprintf("%s/v1/database/%s/query", c.baseURL, dbAddress)

	reqBody := struct {
		Query string `json:"query"`
	}{
		Query: query,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", c.getAuthHeader())
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return result, nil
}
