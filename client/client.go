package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a SpacetimeDB client
type Client struct {
	baseURL    string
	httpClient *http.Client
	wsClient   *websocket.Conn
	token      string
	identity   string
	ctx        context.Context
	cancelFunc context.CancelFunc
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// NewClient creates a new SpacetimeDB client
func NewClient(baseURL string, opts ...ClientOption) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{
		baseURL: parsedURL.String(),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		ctx:        ctx,
		cancelFunc: cancel,
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// WithToken sets the authentication token for the client
func WithToken(token string) ClientOption {
	return func(c *Client) {
		c.token = token
	}
}

// WithIdentity sets the identity for the client
func WithIdentity(identity string) ClientOption {
	return func(c *Client) {
		c.identity = identity
	}
}

// Close closes the client and its connections
func (c *Client) Close() error {
	if c.wsClient != nil {
		if err := c.wsClient.Close(); err != nil {
			return fmt.Errorf("error closing websocket connection: %w", err)
		}
	}
	c.cancelFunc()
	return nil
}

// getAuthHeader returns the Authorization header value
func (c *Client) getAuthHeader() string {
	if c.token == "" {
		return ""
	}
	return fmt.Sprintf("Bearer %s", c.token)
}
