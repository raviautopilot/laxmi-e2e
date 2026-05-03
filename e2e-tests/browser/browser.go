package browser

import (
	"fmt"

	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"

	"github.com/laxmi/e2e-tests/config"
)

// NewWebDriver creates a new Selenium WebDriver based on the provided config.
// It uses a local WebDriver binary (ChromeDriver or GeckoDriver) and does not
// require Docker or Selenium Grid.
// It returns the WebDriver and the underlying service so the caller can stop it.
func NewWebDriver(cfg *config.Config) (selenium.WebDriver, *selenium.Service, error) {
	// Determine which browser binary to use
	var service *selenium.Service
	var err error

	switch cfg.Browser {
	case "chrome":
		service, err = selenium.NewChromeDriverService(
			"/usr/local/bin/chromedriver", // adjust path as needed
			9515,                          // default ChromeDriver port
			nil,                           // no output
			selenium.Output(nil),
		)
	case "firefox":
		service, err = selenium.NewGeckoDriverService(
			"/usr/local/bin/geckodriver", // adjust path as needed
			4444,                         // default GeckoDriver port
			nil,
			selenium.Output(nil),
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
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", service.Port()))
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
