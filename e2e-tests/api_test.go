//go:build api

package e2e_tests

import (
	"io"
	"net/http"
	"testing"
)

// TestExampleAPI demonstrates a basic API test.
func TestExampleAPI(t *testing.T) {
	// Use global apiClient from main_test.go
	if apiClient == nil {
		t.Skip("API client not initialized")
	}

	// Perform a GET request to the base URL
	resp, err := apiClient.Get("/")
	if err != nil {
		t.Fatalf("GET / failed: %v", err)
	}
	defer resp.Body.Close()

	// Verify status code
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	// Read response body (optional)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	if len(body) == 0 {
		t.Error("response body is empty")
	}
}
