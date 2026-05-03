package browser

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/tebeka/selenium"
)

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
