//go:build web

package laxmi_e2e

import (
	"testing"

	"github.com/laxmi/e2e-tests/browser"
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

	// Take a screenshot for documentation
	if _, err := browser.TakeScreenshot(wd, "example-navigation", runID); err != nil {
		t.Logf("screenshot failed: %v", err)
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
