package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the E2E test framework.
type Config struct {
	BaseURL     string        `yaml:"base_url"`
	TestURLs    []string      `yaml:"test_urls"`
	Credentials Credentials   `yaml:"credentials"`
	Browser     string        `yaml:"browser"`
	Headless    bool          `yaml:"headless"`
	Timeout     time.Duration `yaml:"timeout"`
}

// Credentials holds authentication details.
type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Load reads configuration from the given YAML file path.
// Environment variables override YAML values when present.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", path, err)
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %s: %w", path, err)
	}

	// Apply environment variable overrides
	if v := os.Getenv("BASE_URL"); v != "" {
		cfg.BaseURL = v
	}
	if v := os.Getenv("BROWSER"); v != "" {
		cfg.Browser = v
	}
	if v := os.Getenv("HEADLESS"); v == "true" || v == "1" {
		cfg.Headless = true
	}
	if v := os.Getenv("TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err == nil {
			cfg.Timeout = d
		}
	}
	if v := os.Getenv("USERNAME"); v != "" {
		cfg.Credentials.Username = v
	}
	if v := os.Getenv("PASSWORD"); v != "" {
		cfg.Credentials.Password = v
	}

	return cfg, nil
}
