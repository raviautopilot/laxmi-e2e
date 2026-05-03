//go:build api

package api_test

import (
	"io"
	"net/http"
	"testing"

	"github.com/laxmi/e2e-tests/api"
)

// ExampleAPITest demonstrates a basic API test.
func ExampleAPITest(t *testing.T, client *api.Client) {
	t.Helper()

	// Perform a GET request to the base URL
	resp, err := client.Get("/")
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
