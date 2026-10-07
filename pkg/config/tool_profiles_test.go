package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadToolProfilesMergesProjectAndGlobal(t *testing.T) {
	dir := t.TempDir()
	globalPath := filepath.Join(dir, "global", "profiles.toml")
	projectPath := filepath.Join(dir, "project", "profiles.toml")
	writeProfilesFile(t, globalPath, `default_profile = "global-default"

[profiles.global-default]
groups = ["source_read"]

[profiles.shared]
tools = ["GetSource"]
`)
	writeProfilesFile(t, projectPath, `default_profile = "project-default"

[profiles.project-default]
groups = ["transport_read"]

[profiles.shared]
tools = ["ListTransports"]
`)

	got, err := LoadToolProfilesFromPaths(projectPath, globalPath)
	if err != nil {
		t.Fatalf("LoadToolProfilesFromPaths: %v", err)
	}
	if got.GlobalDefault != "global-default" || got.ProjectDefault != "project-default" {
		t.Fatalf("defaults = global %q, project %q", got.GlobalDefault, got.ProjectDefault)
	}
	if _, ok := got.Profiles["global-default"]; !ok {
		t.Fatal("global-only profile was not retained")
	}
	if profile := got.Profiles["shared"]; len(profile.Tools) != 1 || profile.Tools[0] != "ListTransports" {
		t.Fatalf("project profile did not replace global profile: %+v", profile)
	}
}

func TestLoadToolProfilesUsesProjectAndHomePaths(t *testing.T) {
	repo := t.TempDir()
	home := t.TempDir()
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(repo); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Error(err)
		}
	})
	t.Setenv("HOME", home)
	writeProfilesFile(t, filepath.Join(home, ".vsp", "profiles.toml"), `default_profile = "global"

[profiles.global]
groups = ["source_read"]
`)
	writeProfilesFile(t, filepath.Join(repo, ".vsp", "profiles.toml"), `[profiles.project]
groups = ["transport_read"]
`)

	got, err := LoadToolProfiles()
	if err != nil {
		t.Fatalf("LoadToolProfiles: %v", err)
	}
	if got.GlobalDefault != "global" || got.Profiles["global"].Groups[0] != "source_read" || got.Profiles["project"].Groups[0] != "transport_read" {
		t.Fatalf("LoadToolProfiles() = %+v; want global and project profiles", got)
	}
}

func TestToolProfilesResolvePrecedence(t *testing.T) {
	profiles := &ToolProfiles{
		Profiles: map[string]ToolProfile{
			"global":  {Tools: []string{"GetSource"}},
			"project": {Tools: []string{"ListTransports"}},
			"flag":    {Tools: []string{"GetTable"}},
		},
		GlobalDefault:  "global",
		ProjectDefault: "project",
	}

	for _, tc := range []struct {
		name     string
		explicit string
		want     string
	}{
		{name: "explicit overrides project", explicit: "flag", want: "flag"},
		{name: "project overrides global", want: "project"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name, profile, found, err := profiles.Resolve(tc.explicit)
			if err != nil || !found {
				t.Fatalf("Resolve: found=%v err=%v", found, err)
			}
			if name != tc.want {
				t.Fatalf("selected %q, want %q", name, tc.want)
			}
			if len(profile.Tools) != 1 {
				t.Fatalf("selected profile = %+v", profile)
			}
		})
	}
}

func TestToolProfilesResolveGlobalDefaultWhenProjectUnset(t *testing.T) {
	profiles := &ToolProfiles{
		Profiles:      map[string]ToolProfile{"global": {Tools: []string{"GetSource"}}},
		GlobalDefault: "global",
	}
	name, _, found, err := profiles.Resolve("")
	if err != nil || !found || name != "global" {
		t.Fatalf("Resolve() = %q, %v, %v; want global, true, nil", name, found, err)
	}
}

func TestToolProfilesNoDefaultUsesLegacyMode(t *testing.T) {
	profiles, err := LoadToolProfilesFromPaths(filepath.Join(t.TempDir(), "project.toml"), filepath.Join(t.TempDir(), "global.toml"))
	if err != nil {
		t.Fatalf("LoadToolProfilesFromPaths: %v", err)
	}
	name, _, found, err := profiles.Resolve("")
	if err != nil || found || name != "" {
		t.Fatalf("Resolve() = %q, %v, %v; want no profile", name, found, err)
	}
}

func TestToolProfilesUnknownDefaultFails(t *testing.T) {
	profiles := &ToolProfiles{Profiles: map[string]ToolProfile{}, GlobalDefault: "missing"}
	if _, _, _, err := profiles.Resolve(""); err == nil {
		t.Fatal("Resolve() accepted an undefined default profile")
	}
}

func TestLoadToolProfilesRejectsUnknownFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "profiles.toml")
	writeProfilesFile(t, path, `default = "code-scout"
`)
	if _, err := LoadToolProfilesFromPaths("", path); err == nil {
		t.Fatal("LoadToolProfilesFromPaths accepted an unknown field")
	}
}

func writeProfilesFile(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatal(err)
	}
}
