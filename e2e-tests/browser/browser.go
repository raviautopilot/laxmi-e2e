package browser

import (
	"fmt"
	"os"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"

	"github.com/laxmi/e2e-tests/config"
)

// binaryExists checks whether the given file path exists and is executable.
func binaryExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	// Check that it's a regular file (not a directory) and has execute permission
	return !info.IsDir() && (info.Mode()&0111 != 0)
}

// NewWebDriver creates a new Selenium WebDriver based on the provided config.
// It uses a local WebDriver binary (ChromeDriver or GeckoDriver) and does not
// require Docker or Selenium Grid.
// It returns the WebDriver and the underlying service so the caller can stop it.
func NewWebDriver(cfg *config.Config) (selenium.WebDriver, *selenium.Service, error) {
	// Determine which browser binary to use
	var service *selenium.Service
	var port int
	var err error

	switch cfg.Browser {
	case "chrome":
		port = 9515
		chromeDriverPath := "/usr/local/bin/chromedriver" // adjust path as needed
		if !binaryExists(chromeDriverPath) {
			return nil, nil, fmt.Errorf("ChromeDriver binary not found at %s", chromeDriverPath)
		}
		service, err = selenium.NewChromeDriverService(
			chromeDriverPath,
			port, // default ChromeDriver port
		)
	case "firefox":
		port = 4444
		geckoDriverPath := "/usr/local/bin/geckodriver" // adjust path as needed
		if !binaryExists(geckoDriverPath) {
			return nil, nil, fmt.Errorf("GeckoDriver binary not found at %s", geckoDriverPath)
		}
		service, err = selenium.NewGeckoDriverService(
			geckoDriverPath,
			port, // default GeckoDriver port
		)
	default:
		return nil, nil, fmt.Errorf("unsupported browser: %s", cfg.Browser)
	}
	if err != nil {
		return nil, nil, fmt.Errorf("starting %s driver service: %w", cfg.Browser, err)
	}

	// Build capabilities
	caps := selenium.Capabilities{"browserName": cfg.Browser}
	if cfg.Browser == "chrome" {
		chromeCaps := chrome.Capabilities{
			Args: []string{
				"--no-sandbox",
				"--disable-dev-shm-usage",
			},
		}
		if cfg.Headless {
			chromeCaps.Args = append(chromeCaps.Args, "--headless")
		}
		caps.AddChrome(chromeCaps)
	}

	// Connect to the WebDriver instance
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		service.Stop()
		return nil, nil, fmt.Errorf("connecting to WebDriver: %w", err)
	}

	// Set implicit wait timeout
	if cfg.Timeout > 0 {
		if err := wd.SetImplicitWaitTimeout(cfg.Timeout); err != nil {
			service.Stop()
			wd.Quit()
			return nil, nil, fmt.Errorf("setting implicit wait: %w", err)
		}
	}

	return wd, service, nil
}

// Cleanup stops the WebDriver service and quits the driver.
func Cleanup(wd selenium.WebDriver, svc *selenium.Service) {
	if wd != nil {
		wd.Quit()
	}
	if svc != nil {
		svc.Stop()
	}
}
