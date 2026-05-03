package e2e_tests

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/tebeka/selenium"

	"github.com/laxmi/e2e-tests/api"
	"github.com/laxmi/e2e-tests/browser"
	"github.com/laxmi/e2e-tests/config"
)

// Global variables accessible by all test files.
var (
	cfg      *config.Config
	wd       selenium.WebDriver
	apiClient *api.Client
)

// TestMain is the entry point for all E2E tests.
func TestMain(m *testing.M) {
	// Parse command-line flags (e.g., -test.v)
	flag.Parse()

	// Load configuration
	var err error
	cfg, err = config.Load("config/config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize browser (only if web tests are enabled)
	if hasWebTests() {
		wd, err = browser.NewWebDriver(cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create WebDriver: %v\n", err)
			os.Exit(1)
		}
		defer browser.Cleanup(wd)
	}

	// Initialize API client
	apiClient = api.NewClient(cfg)

	// Run all tests
	exitCode := m.Run()

	os.Exit(exitCode)
}

// hasWebTests checks whether any web tests are registered.
// This is a simple heuristic: if the build tag "web" is present, we assume web tests exist.
func hasWebTests() bool {
	// We rely on build tags; if the binary was compiled without "web" tag,
	// the web test files won't be included, so we can skip browser init.
	// For simplicity, always attempt to initialize the browser.
	return true
}
