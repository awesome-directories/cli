package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/caarlos0/env/v11"
	"gopkg.in/yaml.v3"
)

var (
	BuildSupabaseURL     string
	BuildSupabaseAnonKey string
)

// Config holds all configuration for the CLI
type Config struct {
	// Supabase configuration
	SupabaseURL     string `env:"SUPABASE_URL" yaml:"supabase_url"`
	SupabaseAnonKey string `env:"SUPABASE_ANON_KEY" yaml:"supabase_anon_key"`

	// Auth configuration
	AuthToken string `env:"AUTH_TOKEN" yaml:"auth_token"`

	// Cache configuration
	CacheDir string        `env:"CACHE_DIR" yaml:"cache_dir"`
	CacheTTL time.Duration `env:"CACHE_TTL" yaml:"cache_ttl"`

	// General settings
	Debug   bool `env:"DEBUG" yaml:"debug"`
	NoColor bool `env:"NO_COLOR" yaml:"no_color"`
}

// Default values
const (
	DefaultCacheTTL = 24 * time.Hour
)

// Load loads configuration from environment and config file
func Load() (*Config, error) {
	cfg := &Config{
		SupabaseURL:     BuildSupabaseURL,
		SupabaseAnonKey: BuildSupabaseAnonKey,
		CacheTTL:        DefaultCacheTTL,
	}

	// Get config directory
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	// Set cache directory
	cfg.CacheDir = filepath.Join(configDir, "cache")

	// Load from config file if it exists
	configFile := filepath.Join(configDir, "config.yaml")
	if _, err := os.Stat(configFile); err == nil {
		if err := loadFromFile(configFile, cfg); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	}

	// Override with environment variables
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	if cfg.SupabaseURL == "" || cfg.SupabaseAnonKey == "" {
		return nil, fmt.Errorf("supabase URL and anon key are missing. provide them with env var SUPABASE_URL & SUPABASE_ANON_KEY")
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(cfg.CacheDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return cfg, nil
}

// Save saves configuration to file
func (c *Config) Save() error {
	configDir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.yaml")

	// Marshal to YAML
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// getConfigDir returns the configuration directory path
func getConfigDir() (string, error) {
	// Try XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "awesome-directories"), nil
	}

	// Fall back to ~/.config
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".config", "awesome-directories"), nil
}

// GetConfigDir returns the configuration directory (public helper)
func GetConfigDir() (string, error) {
	return getConfigDir()
}

// loadFromFile loads configuration from YAML file
func loadFromFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg)
}
