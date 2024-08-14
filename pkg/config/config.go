package config

import (
	"path"
)

const CONFIG_VERSION = "0.1.0"

type Config struct {
	ConfigVersion     string `json:"config_version"`
	WorkflowDirectory string `json:"workflow_directory"`
	OutputDirectory   string `json:"output_directory"`
	LogDirectory      string `json:"log_directory"`
	DisableLogs       bool   `json:"disable_logs"`
	AlwaysConfirm     bool   `json:"always_confirm"`
}

func NewConfig() *Config {
	return &Config{}
}

func NewConfigWithDefaults(homeDir string) *Config {
	return &Config{
		ConfigVersion:     CONFIG_VERSION,
		WorkflowDirectory: path.Join(homeDir, "imgscal", "workflow"),
		LogDirectory:      path.Join(homeDir, "imgscal", "log"),
		OutputDirectory:   path.Join(homeDir, "imgscal", "output"),
		DisableLogs:       false,
		AlwaysConfirm:     false,
	}
}
