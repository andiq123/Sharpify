package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/andiq123/sharpify/internal/rules"
)

type Config struct {
	TargetVersion string   `json:"targetVersion"`
	SafeOnly      bool     `json:"safeOnly"`
	BackupEnabled bool     `json:"backupEnabled"`
	DisabledRules []string `json:"disabledRules,omitempty"`
	WorkingPath   string   `json:"workingPath,omitempty"`
}

func DefaultConfig() *Config {
	return &Config{
		TargetVersion: "12",
		SafeOnly:      true,
		BackupEnabled: false,
		DisabledRules: []string{},
		WorkingPath:   "",
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

// IsRuleDisabled checks if a rule is in the disabled list
func (c *Config) IsRuleDisabled(name string) bool {
	for _, r := range c.DisabledRules {
		if r == name {
			return true
		}
	}
	return false
}

// SetRuleDisabled adds or removes a rule from the disabled list
func (c *Config) SetRuleDisabled(name string, disabled bool) {
	if disabled {
		if !c.IsRuleDisabled(name) {
			c.DisabledRules = append(c.DisabledRules, name)
		}
	} else {
		newList := make([]string, 0, len(c.DisabledRules))
		for _, r := range c.DisabledRules {
			if r != name {
				newList = append(newList, r)
			}
		}
		c.DisabledRules = newList
	}
}

// GetEnabledRules returns the list of rules filtering out disabled ones
func (c *Config) GetEnabledRules(allRules []rules.Rule) []rules.Rule {
	result := make([]rules.Rule, 0, len(allRules))
	for _, r := range allRules {
		if !c.IsRuleDisabled(r.Name()) {
			result = append(result, r)
		}
	}
	return result
}

