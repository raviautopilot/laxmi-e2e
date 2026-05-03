//go:build web

package web_test

import (
	"testing"

	"github.com/tebeka/selenium"

	"github.com/laxmi/e2e-tests/config"
)

// ExampleWebTest demonstrates a basic Selenium test.
func ExampleWebTest(t *testing.T, wd selenium.WebDriver, cfg *config.Config) {
	t.Helper()

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
