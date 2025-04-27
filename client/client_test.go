package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		wantErr bool
	}{
		{
			name:    "valid URL",
			baseURL: "https://example.com",
			wantErr: false,
		},
		{
			name:    "invalid URL",
			baseURL: "not-a-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(tt.baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateIdentity(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/v1/identity" {
			t.Errorf("Expected path /v1/identity, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"identity":"test-identity","token":"test-token"}`))
	}))
	defer server.Close()

	// Create a client with the test server URL
	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test CreateIdentity
	identityResp, err := client.CreateIdentity()
	if err != nil {
		t.Fatalf("CreateIdentity() error = %v", err)
	}

	if identityResp.Identity != "test-identity" {
		t.Errorf("Expected identity test-identity, got %s", identityResp.Identity)
	}
	if identityResp.Token != "test-token" {
		t.Errorf("Expected token test-token, got %s", identityResp.Token)
	}
}

func TestGetAuthHeader(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		expected string
	}{
		{
			name:     "with token",
			token:    "test-token",
			expected: "Bearer test-token",
		},
		{
			name:     "without token",
			token:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{token: tt.token}
			got := client.getAuthHeader()
			if got != tt.expected {
				t.Errorf("getAuthHeader() = %v, want %v", got, tt.expected)
			}
		})
	}
}
