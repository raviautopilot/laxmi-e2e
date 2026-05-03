//go:build web

package e2e_tests

import (
	"testing"
)

// TestExampleWeb demonstrates a basic Selenium test.
func TestExampleWeb(t *testing.T) {
	// Use global wd and cfg from main_test.go
	if wd == nil {
		t.Skip("WebDriver not initialized")
	}

	// Navigate to the base URL
	if err := wd.Get(cfg.BaseURL); err != nil {
		t.Fatalf("failed to navigate to %s: %v", cfg.BaseURL, err)
	}

	// Verify page title
	title, err := wd.Title()
	if err != nil {
		t.Fatalf("failed to get page title: %v", err)
	}
	if title == "" {
		t.Error("page title is empty")
	}
}
