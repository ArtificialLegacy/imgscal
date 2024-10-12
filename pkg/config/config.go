package config

import (
	"path"
)

const CONFIG_VERSION = "0.2.0"

type Config struct {
	ConfigVersion     string `json:"config_version"`
	ConfigDirectory   string `json:"config_directory"`
	WorkflowDirectory string `json:"workflow_directory"`
	OutputDirectory   string `json:"output_directory"`
	LogDirectory      string `json:"log_directory"`
	PluginDirectory   string `json:"plugin_directory"`
	DisableLogs       bool   `json:"disable_logs"`
	AlwaysConfirm     bool   `json:"always_confirm"`
}

func NewConfig() *Config {
	return &Config{}
}

func NewConfigWithDefaults(homeDir string) *Config {
	return &Config{
		ConfigVersion:     CONFIG_VERSION,
		ConfigDirectory:   path.Join(homeDir, "imgscal", "config"),
		WorkflowDirectory: path.Join(homeDir, "imgscal", "workflow"),
		LogDirectory:      path.Join(homeDir, "imgscal", "log"),
		OutputDirectory:   path.Join(homeDir, "imgscal", "output"),
		PluginDirectory:   path.Join(homeDir, "imgscal", "plugin"),
		DisableLogs:       false,
		AlwaysConfirm:     false,
	}
}
