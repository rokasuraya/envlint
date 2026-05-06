package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envlint/internal/profile"
)

const defaultProfilesFile = ".envlint-profiles.yaml"

// profileFlags holds the parsed profile-related CLI flags.
type profileFlags struct {
	ProfileName  string
	ProfilesFile string
}

// registerProfileFlags adds profile-related flags to the given FlagSet.
func registerProfileFlags(fs *flag.FlagSet, pf *profileFlags) {
	fs.StringVar(&pf.ProfileName, "profile", "", "named profile to use from profiles config file")
	fs.StringVar(&pf.ProfilesFile, "profiles-file", defaultProfilesFile, "path to profiles YAML config")
}

// resolveProfile loads the profiles config and returns the env/schema paths
// for the requested profile name. It returns empty strings when no profile
// flag was provided so callers can fall back to --env / --schema flags.
func resolveProfile(pf profileFlags) (envFile, schemaFile string, err error) {
	if pf.ProfileName == "" {
		return "", "", nil
	}
	cfg, err := profile.LoadConfig(pf.ProfilesFile)
	if err != nil {
		return "", "", fmt.Errorf("loading profiles: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return "", "", fmt.Errorf("invalid profiles config: %w", err)
	}
	p, err := cfg.Get(pf.ProfileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "available profiles: %v\n", cfg.Names())
		return "", "", err
	}
	return p.EnvFile, p.SchemaFile, nil
}
