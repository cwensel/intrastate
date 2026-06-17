// Package config loads and validates the project config file
// (intrastate.toml). It owns config discovery, parsing, and the
// structured errors that surface when the file is missing or invalid —
// every error is a clierr.CLIError so the CLI's never-silent contract
// holds at the config boundary.
//
// The parser is intentionally not wired yet: ConfigFileName and the
// Config shape lock the on-disk contract, and Load establishes the
// discovery + error envelope. A future cut chooses a TOML library and
// fills in parse(); the call sites and error codes are already in place.
package config

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/newcoinc/intrastate/internal/cli/clierr"
)

// ConfigFileName is the project config file discovered by walking up
// from the working directory.
const ConfigFileName = "intrastate.toml"

// CurrentSchemaVersion is the config schema version this build writes
// and accepts. Bump deliberately when the on-disk shape changes.
const CurrentSchemaVersion = "1"

// Config is the loaded, validated view of intrastate.toml. Grow this as
// the config surface is defined; keep zero values meaningful so an
// absent section reads as "use defaults".
type Config struct {
	SchemaVersion string
	// Raw is the unparsed file bytes, retained for surgical rewrites
	// that must preserve formatting and comments.
	Raw []byte
	// Path is the resolved absolute path the config was loaded from.
	Path string
}

// Warning is a non-fatal load-time advisory (e.g. a deprecated key).
type Warning struct {
	Code    string
	Message string
}

// Discover walks up from dir looking for ConfigFileName and returns the
// first match. Returns a config-not-found CLIError when none is found
// before the filesystem root.
func Discover(dir string) (string, error) {
	cur := dir
	for {
		candidate := filepath.Join(cur, ConfigFileName)
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return "", &clierr.CLIError{
				Code:    "config-not-found",
				Message: "no " + ConfigFileName + " found in " + dir + " or any parent directory",
				Group:   clierr.GroupUserEnv,
				Hint:    "run `intrastate init` to create one",
			}
		}
		cur = parent
	}
}

// Load reads and validates the config at path. Errors carry stable
// structured codes (config-not-found, config-invalid); callers test
// them via clierr.ErrorCode.
func Load(path string) (*Config, []Warning, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil, &clierr.CLIError{
				Code:    "config-not-found",
				Message: "config file not found: " + path,
				Group:   clierr.GroupUserEnv,
			}
		}
		return nil, nil, &clierr.CLIError{
			Code:    "config-read-error",
			Message: "could not read config: " + path,
			Detail:  err.Error(),
			Group:   clierr.GroupInternal,
			Cause:   err,
		}
	}

	cfg := &Config{
		SchemaVersion: CurrentSchemaVersion,
		Raw:           data,
		Path:          path,
	}
	// TODO: parse `data` into cfg once a TOML library is chosen, then
	// validate cfg.SchemaVersion against CurrentSchemaVersion and emit a
	// `config-invalid` CLIError on mismatch.
	return cfg, nil, nil
}
