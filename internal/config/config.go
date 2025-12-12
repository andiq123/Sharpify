package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/andiq123/sharpify/internal/rules"
)

type Config struct {
	TargetVersion string `json:"targetVersion"`
	SafeOnly      bool   `json:"safeOnly"`
	BackupEnabled bool   `json:"backupEnabled"`
}

func DefaultConfig() *Config {
	return &Config{
		TargetVersion: "12",
		SafeOnly:      true,
		BackupEnabled: false,
	}
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".sharpify.json")
}

func Load() *Config {
	cfg := DefaultConfig()

	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg
	}

	_ = json.Unmarshal(data, cfg)
	return cfg
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0644)
}

func (c *Config) GetVersion() rules.CSharpVersion {
	switch c.TargetVersion {
	case "6":
		return rules.CSharp6
	case "7":
		return rules.CSharp7
	case "8":
		return rules.CSharp8
	case "9":
		return rules.CSharp9
	case "10":
		return rules.CSharp10
	case "11":
		return rules.CSharp11
	case "13":
		return rules.CSharp13
	default:
		return rules.CSharp12
	}
}

func (c *Config) SetVersion(v rules.CSharpVersion) {
	switch v {
	case rules.CSharp6:
		c.TargetVersion = "6"
	case rules.CSharp7:
		c.TargetVersion = "7"
	case rules.CSharp8:
		c.TargetVersion = "8"
	case rules.CSharp9:
		c.TargetVersion = "9"
	case rules.CSharp10:
		c.TargetVersion = "10"
	case rules.CSharp11:
		c.TargetVersion = "11"
	case rules.CSharp13:
		c.TargetVersion = "13"
	default:
		c.TargetVersion = "12"
	}
}
