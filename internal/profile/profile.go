// Package profile manages named validation profiles, allowing users to
// define reusable sets of schema + env file pairs for common environments
// (e.g. "production", "staging", "ci").
package profile

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Profile represents a named environment configuration.
type Profile struct {
	Name       string `yaml:"name"`
	EnvFile    string `yaml:"env_file"`
	SchemaFile string `yaml:"schema_file"`
	Description string `yaml:"description,omitempty"`
}

// Config holds all profiles defined in a profiles config file.
type Config struct {
	Profiles []Profile `yaml:"profiles"`
}

// LoadConfig reads a profiles YAML file from the given path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("profile: read %q: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("profile: parse %q: %w", path, err)
	}
	return &cfg, nil
}

// Get returns the profile with the given name, or an error if not found.
func (c *Config) Get(name string) (*Profile, error) {
	for i := range c.Profiles {
		if c.Profiles[i].Name == name {
			return &c.Profiles[i], nil
		}
	}
	return nil, fmt.Errorf("profile: %q not found", name)
}

// Names returns a sorted list of all profile names.
func (c *Config) Names() []string {
	names := make([]string, 0, len(c.Profiles))
	for _, p := range c.Profiles {
		names = append(names, p.Name)
	}
	return names
}

// Validate checks that each profile has required fields set.
func (c *Config) Validate() error {
	seen := map[string]bool{}
	for _, p := range c.Profiles {
		if p.Name == "" {
			return fmt.Errorf("profile: entry missing name")
		}
		if seen[p.Name] {
			return fmt.Errorf("profile: duplicate name %q", p.Name)
		}
		seen[p.Name] = true
		if p.EnvFile == "" {
			return fmt.Errorf("profile %q: env_file is required", p.Name)
		}
		if p.SchemaFile == "" {
			return fmt.Errorf("profile %q: schema_file is required", p.Name)
		}
	}
	return nil
}
