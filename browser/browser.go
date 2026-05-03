package browser

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

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

// findDriverBinary locates the WebDriver binary in the system PATH.
// It returns the full path and an error if not found.
func findDriverBinary(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("%s not found in PATH. Please install it:\n"+
			"  sudo apt-get update && sudo apt-get install -y chromium-chromedriver\n"+
			"  (or for Firefox: sudo apt-get install -y firefox-geckodriver)", name)
	}
	return path, nil
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
		chromeDriverPath, err := findDriverBinary("chromedriver")
		if err != nil {
			return nil, nil, err
		}
		if !binaryExists(chromeDriverPath) {
			return nil, nil, fmt.Errorf("ChromeDriver binary not found at %s", chromeDriverPath)
		}
		service, err = selenium.NewChromeDriverService(
			chromeDriverPath,
			port, // default ChromeDriver port
		)
	case "firefox":
		port = 4444
		geckoDriverPath, err := findDriverBinary("geckodriver")
		if err != nil {
			return nil, nil, err
		}
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

var (
	mu      sync.Mutex
	counter int
)

// TakeScreenshot captures a screenshot of the current browser state,
// saves it to screenshots/<runID>/<seq>_<name>.png, and returns the
// absolute file path. The sequence number is auto-incremented per run.
func TakeScreenshot(wd selenium.WebDriver, name string, runID string) (string, error) {
	mu.Lock()
	counter++
	seq := counter
	mu.Unlock()

	seqStr := fmt.Sprintf("%03d", seq)
	filename := seqStr + "_" + name + ".png"
	dir := filepath.Join("screenshots", runID)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating screenshot directory %s: %w", dir, err)
	}

	pngBytes, err := wd.Screenshot()
	if err != nil {
		return "", fmt.Errorf("capturing screenshot: %w", err)
	}

	filePath := filepath.Join(dir, filename)
	if err := os.WriteFile(filePath, pngBytes, 0644); err != nil {
		return "", fmt.Errorf("writing screenshot file %s: %w", filePath, err)
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return "", fmt.Errorf("getting absolute path for %s: %w", filePath, err)
	}

	return absPath, nil
}

// GenerateRunID creates a unique run identifier based on current time
// and a random suffix.
func GenerateRunID() string {
	now := time.Now().Format("20060102_150405")
	// 4 random alphanumeric characters
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	suffix := make([]byte, 4)
	for i := range suffix {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			// fallback to a simple counter
			suffix[i] = '0'
			continue
		}
		suffix[i] = letters[n.Int64()]
	}
	return fmt.Sprintf("run_%s_%s", now, string(suffix))
}
