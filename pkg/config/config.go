package config

import (
	"path"
)

const (
	CONFIG_VERSION = "0.5.0"
	SCHEMA         = "https://gist.githubusercontent.com/ArtificialLegacy/bf37b79d4fc943006f333cc35467266c/raw/933fdffd6d871d3bf5a281a7815b7d408fcd51b2/imgscal_config.json"
)

type Config struct {
	Schema            string `json:"$schema"`
	ConfigVersion     string `json:"config_version"`
	ConfigDirectory   string `json:"config_directory"`
	WorkflowDirectory string `json:"workflow_directory"`
	OutputDirectory   string `json:"output_directory"`
	InputDirectory    string `json:"input_directory"`
	LogDirectory      string `json:"log_directory"`
	PluginDirectory   string `json:"plugin_directory"`
	DefaultAuthor     string `json:"default_author"`
	DisableLogs       bool   `json:"disable_logs"`
	AlwaysConfirm     bool   `json:"always_confirm"`
	DisableBell       bool   `json:"disable_bell"`
}

func NewConfig() *Config {
	return &Config{
		Schema: SCHEMA,
	}
}

func NewConfigWithDefaults(homeDir string) *Config {
	return &Config{
		Schema:            SCHEMA,
		ConfigVersion:     CONFIG_VERSION,
		ConfigDirectory:   path.Join(homeDir, "imgscal", "config"),
		WorkflowDirectory: path.Join(homeDir, "imgscal", "workflow"),
		LogDirectory:      path.Join(homeDir, "imgscal", "log"),
		OutputDirectory:   path.Join(homeDir, "imgscal", "output"),
		InputDirectory:    path.Join(homeDir, "imgscal", "input"),
		PluginDirectory:   path.Join(homeDir, "imgscal", "plugin"),
		DefaultAuthor:     "",
		DisableLogs:       false,
		AlwaysConfirm:     false,
		DisableBell:       false,
	}
}
