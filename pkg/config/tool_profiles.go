package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// ToolProfile describes a named MCP tool set. Tool names and groups are
// expanded by internal/mcp, where the registered tool catalog is known.
type ToolProfile struct {
	Description string   `toml:"description,omitempty"`
	Groups      []string `toml:"groups,omitempty"`
	Tools       []string `toml:"tools,omitempty"`
	Exclude     []string `toml:"exclude,omitempty"`
}

type toolProfilesFile struct {
	DefaultProfile string                 `toml:"default_profile,omitempty"`
	Profiles       map[string]ToolProfile `toml:"profiles,omitempty"`
}

// ToolProfiles is the merged project/global profile registry. Project
// profiles replace global profiles with the same name; otherwise global
// profiles remain available.
type ToolProfiles struct {
	Profiles       map[string]ToolProfile
	ProjectDefault string
	GlobalDefault  string
}

// ToolProfilesPaths returns the project and global profile files. They are
// discovered independently of .vsp.json because system configuration uses
// first-found semantics rather than merging project and home files.
func ToolProfilesPaths() (projectPath, globalPath string) {
	projectPath = filepath.Join(".vsp", "profiles.toml")
	if home, err := os.UserHomeDir(); err == nil {
		globalPath = filepath.Join(home, ".vsp", "profiles.toml")
	}
	return projectPath, globalPath
}

// LoadToolProfiles loads and merges the project and global profile files.
// Missing files are fine; malformed or unreadable files are reported.
func LoadToolProfiles() (*ToolProfiles, error) {
	projectPath, globalPath := ToolProfilesPaths()
	return LoadToolProfilesFromPaths(projectPath, globalPath)
}

// LoadToolProfilesFromPaths loads profile files at explicit paths. The global
// file is read first, then project definitions and defaults take precedence.
func LoadToolProfilesFromPaths(projectPath, globalPath string) (*ToolProfiles, error) {
	out := &ToolProfiles{Profiles: make(map[string]ToolProfile)}
	paths := []struct {
		path    string
		project bool
	}{
		{path: globalPath},
		{path: projectPath, project: true},
	}
	seen := make(map[string]bool, len(paths))
	for _, candidate := range paths {
		if candidate.path == "" {
			continue
		}
		abs, err := filepath.Abs(candidate.path)
		if err == nil && seen[abs] {
			continue
		}
		if err == nil {
			seen[abs] = true
		}

		file, err := readToolProfiles(candidate.path)
		if err != nil {
			return nil, err
		}
		if file == nil {
			continue
		}
		for name, profile := range file.Profiles {
			out.Profiles[name] = profile
		}
		if candidate.project {
			out.ProjectDefault = file.DefaultProfile
		} else {
			out.GlobalDefault = file.DefaultProfile
		}
	}
	return out, nil
}

func readToolProfiles(path string) (*toolProfilesFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read tool profiles %s: %w", path, err)
	}
	var file toolProfilesFile
	decoder := toml.NewDecoder(bytes.NewReader(data)).DisallowUnknownFields()
	if err := decoder.Decode(&file); err != nil {
		return nil, fmt.Errorf("parse tool profiles %s: %w", path, err)
	}
	return &file, nil
}

// Resolve selects a profile by explicit name, then project default, then
// global default. A false found result means no profile is configured.
func (p *ToolProfiles) Resolve(explicit string) (name string, profile ToolProfile, found bool, err error) {
	if explicit != "" {
		name = explicit
	} else if p != nil && p.ProjectDefault != "" {
		name = p.ProjectDefault
	} else if p != nil {
		name = p.GlobalDefault
	}
	if name == "" {
		return "", ToolProfile{}, false, nil
	}
	if p == nil {
		return "", ToolProfile{}, false, fmt.Errorf("MCP tool profile %q was selected but no profile files were loaded", name)
	}
	profile, found = p.Profiles[name]
	if !found {
		return "", ToolProfile{}, false, fmt.Errorf("MCP tool profile %q is not defined", name)
	}
	return name, profile, true, nil
}
