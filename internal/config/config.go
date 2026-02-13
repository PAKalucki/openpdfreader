// Package config provides application configuration management.
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds application configuration.
type Config struct {
	RecentFiles    []string `json:"recent_files"`
	LastOpenDir    string   `json:"last_open_dir"`
	WindowWidth    int      `json:"window_width"`
	WindowHeight   int      `json:"window_height"`
	Theme          string   `json:"theme"` // "light", "dark", "system"
	DefaultZoom    float64  `json:"default_zoom"`
	ShowThumbnails bool     `json:"show_thumbnails"`
}

// Default returns the default configuration.
func Default() *Config {
	return &Config{
		RecentFiles:    []string{},
		LastOpenDir:    "",
		WindowWidth:    1200,
		WindowHeight:   800,
		Theme:          "system",
		DefaultZoom:    1.0,
		ShowThumbnails: true,
	}
}

// Load loads configuration from the config file.
func Load() *Config {
	configPath := getConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		return Default()
	}

	config := Default()
	if err := json.Unmarshal(data, config); err != nil {
		return Default()
	}

	return config
}

// Save writes the configuration to disk.
func (c *Config) Save() error {
	configPath := getConfigPath()

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// AddRecentFile adds a file to the recent files list.
func (c *Config) AddRecentFile(path string) {
	// Remove if already exists
	for i, f := range c.RecentFiles {
		if f == path {
			c.RecentFiles = append(c.RecentFiles[:i], c.RecentFiles[i+1:]...)
			break
		}
	}

	// Add to front
	c.RecentFiles = append([]string{path}, c.RecentFiles...)

	// Limit to 10 recent files
	if len(c.RecentFiles) > 10 {
		c.RecentFiles = c.RecentFiles[:10]
	}
}

func getConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir = "."
	}
	return filepath.Join(configDir, "openpdfreader", "config.json")
}
