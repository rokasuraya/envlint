package schema

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// rawSchema is the intermediate representation used for YAML unmarshalling.
type rawSchema struct {
	Vars map[string]struct {
		Required    bool    `yaml:"required"`
		Type        string  `yaml:"type"`
		Pattern     string  `yaml:"pattern"`
		AllowEmpty  bool    `yaml:"allow_empty"`
		Description string  `yaml:"description"`
	} `yaml:"vars"`
}

// LoadFile reads and parses a YAML schema file from the given path.
func LoadFile(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: reading file %q: %w", path, err)
	}
	return Load(data)
}

// Load parses a YAML schema from raw bytes.
func Load(data []byte) (*Schema, error) {
	var raw rawSchema
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("schema: parsing YAML: %w", err)
	}

	s := &Schema{Vars: make(map[string]VarSchema, len(raw.Vars))}
	for key, v := range raw.Vars {
		vt := VarType(v.Type)
		if vt == "" {
			vt = TypeString
		} else if !vt.IsValid() {
			return nil, fmt.Errorf("schema: var %q has unknown type %q", key, v.Type)
		}
		s.Vars[key] = VarSchema{
			Required:    v.Required,
			Type:        vt,
			Pattern:     v.Pattern,
			AllowEmpty:  v.AllowEmpty,
			Description: v.Description,
		}
	}
	return s, nil
}
