package config

import (
	"path"
)

const CONFIG_VERSION = "0.4.0"

type Config struct {
	ConfigVersion     string `json:"config_version"`
	ConfigDirectory   string `json:"config_directory"`
	WorkflowDirectory string `json:"workflow_directory"`
	OutputDirectory   string `json:"output_directory"`
	InputDirectory    string `json:"input_directory"`
	LogDirectory      string `json:"log_directory"`
	PluginDirectory   string `json:"plugin_directory"`
	DisableLogs       bool   `json:"disable_logs"`
	AlwaysConfirm     bool   `json:"always_confirm"`
	DefaultAuthor     string `json:"default_author"`
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
		InputDirectory:    path.Join(homeDir, "imgscal", "input"),
		PluginDirectory:   path.Join(homeDir, "imgscal", "plugin"),
		DisableLogs:       false,
		AlwaysConfirm:     false,
		DefaultAuthor:     "",
	}
}
