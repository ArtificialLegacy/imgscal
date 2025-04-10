package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/ArtificialLegacy/imgscal/pkg/config"
)

func main() {
	cfgDir, err := os.UserConfigDir()
	if err != nil {
		panic(fmt.Sprintf("cannot access user config directory! (%s)", err))
	}

	cfgPath := path.Join(cfgDir, "imgscal", "config.json")
	if envCfg, exists := os.LookupEnv("IMGSCAL_CONFIG"); exists {
		cfgPath = envCfg
	}

	_, err = os.Stat(cfgPath)
	if err != nil {
		panic(fmt.Sprintf("no config file, cannot access log directory setting! (%s)", err))
	}

	b, err := os.ReadFile(cfgPath)
	if err != nil {
		panic(fmt.Sprintf("failed to read config file! (%s)", err))
	}

	cfg := config.NewConfig()
	err = json.Unmarshal(b, cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal json config file! (%s)", err))
	}

	_, err = os.Stat(cfg.LogDirectory)
	if err != nil {
		panic(fmt.Sprintf("failed to find log directory! (%s)", err))
	}

	latestPath := path.Join(cfg.LogDirectory, "@latest.txt")
	_, err = os.Stat(latestPath)
	if err != nil {
		panic(fmt.Sprintf("no @latest.txt log available! (%s)", err))
	}

	lb, err := os.ReadFile(latestPath)
	if err != nil {
		panic(fmt.Sprintf("failed to read @latest.txt! (%s)", err))
	}

	fmt.Print(string(lb[:]))
}
