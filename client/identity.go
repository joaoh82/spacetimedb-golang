package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// IdentityResponse represents the response from identity-related endpoints
type IdentityResponse struct {
	Identity string `json:"identity"`
	Token    string `json:"token"`
}

// CreateIdentity creates a new identity and returns the identity and token
func (c *Client) CreateIdentity() (*IdentityResponse, error) {
	url := fmt.Sprintf("%s/v1/identity", c.baseURL)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var identityResp IdentityResponse
	if err := json.NewDecoder(resp.Body).Decode(&identityResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	fmt.Println("identityResp", identityResp)

	return &identityResp, nil
}

// VerifyIdentity verifies an identity and token pair
func (c *Client) VerifyIdentity(identity string) error {
	if c.token == "" {
		return fmt.Errorf("token is required for identity verification")
	}

	url := fmt.Sprintf("%s/v1/identity/%s/verify", c.baseURL, identity)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", c.getAuthHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// GetDatabases returns a list of databases owned by an identity
func (c *Client) GetDatabases(identity string) ([]string, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token is required for getting databases")
	}

	url := fmt.Sprintf("%s/v1/identity/%s/databases", c.baseURL, identity)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", c.getAuthHeader())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Addresses []string `json:"addresses"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response.Addresses, nil
}
