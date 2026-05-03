package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/laxmi/e2e-tests/config"
)

// Client is a lightweight HTTP client for API testing.
type Client struct {
	baseURL    string
	httpClient *http.Client
	authToken  string
}

// NewClient creates a new API client from the given config.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL: cfg.BaseURL,
		httpClient: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// SetAuthToken sets a bearer token for subsequent requests.
func (c *Client) SetAuthToken(token string) {
	c.authToken = token
}

// Get performs an HTTP GET request to the given path.
func (c *Client) Get(path string) (*http.Response, error) {
	return c.doRequest(http.MethodGet, path, nil)
}

// Post performs an HTTP POST request with an optional JSON body.
func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
	return c.doRequest(http.MethodPost, path, body)
}

// Put performs an HTTP PUT request with an optional JSON body.
func (c *Client) Put(path string, body interface{}) (*http.Response, error) {
	return c.doRequest(http.MethodPut, path, body)
}

// Delete performs an HTTP DELETE request.
func (c *Client) Delete(path string) (*http.Response, error) {
	return c.doRequest(http.MethodDelete, path, nil)
}

func (c *Client) doRequest(method, path string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	if c.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.authToken)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	return resp, nil
}
