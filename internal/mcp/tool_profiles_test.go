package mcp

import (
	"path/filepath"
	"testing"

	"github.com/oisee/vibing-steampunk/pkg/config"
)

func TestExpandToolProfileCombinesGroupsAndAppliesExclusions(t *testing.T) {
	got, err := ExpandToolProfile(config.ToolProfile{
		Groups:  []string{"transport_read", "transport_manage"},
		Tools:   []string{"SAP"},
		Exclude: []string{"ReleaseTransport", "GetAbapHelp"},
	})
	if err != nil {
		t.Fatalf("ExpandToolProfile: %v", err)
	}
	for _, name := range []string{"SAP", "GetConnectionInfo", "GetFeatures", "ListTransports", "GetTransport", "CreateTransport", "DeleteTransport"} {
		if !got[name] {
			t.Errorf("tool %s missing from expanded profile", name)
		}
	}
	for _, name := range []string{"ReleaseTransport", "GetAbapHelp"} {
		if got[name] {
			t.Errorf("excluded tool %s remains in expanded profile", name)
		}
	}
}

func TestExpandToolProfileRejectsUnknownReferences(t *testing.T) {
	for _, tc := range []struct {
		name    string
		profile config.ToolProfile
	}{
		{name: "tool", profile: config.ToolProfile{Tools: []string{"NotATool"}}},
		{name: "group", profile: config.ToolProfile{Groups: []string{"not_a_group"}}},
		{name: "exclude", profile: config.ToolProfile{Exclude: []string{"NotATool"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := ExpandToolProfile(tc.profile); err == nil {
				t.Fatal("ExpandToolProfile accepted an unknown reference")
			}
		})
	}
}

func TestProfileControlsExposureAcrossModes(t *testing.T) {
	profileTools, err := ExpandToolProfile(config.ToolProfile{Tools: []string{"GetSource"}})
	if err != nil {
		t.Fatalf("ExpandToolProfile: %v", err)
	}
	for _, mode := range []string{"focused", "expert", "hyperfocused"} {
		t.Run(mode, func(t *testing.T) {
			server := NewServer(&Config{
				BaseURL:      "https://example.invalid",
				Mode:         mode,
				ProfileName:  "test-profile",
				ProfileTools: profileTools,
				ToolsConfig:  map[string]bool{"CreateTable": true}, // ignored while a profile is active
			})
			got := toolNameSet(server.RegisteredTools())
			for _, name := range []string{"GetSource", "GetConnectionInfo", "GetFeatures", "GetAbapHelp"} {
				if !got[name] {
					t.Errorf("profile tool %s was not registered", name)
				}
			}
			for _, name := range []string{"CreateTable", "SAP", "GetProgram", "WriteSource"} {
				if got[name] {
					t.Errorf("out-of-profile tool %s was registered", name)
				}
			}
		})
	}
}

func TestProfileIncludesSAPOnlyWhenSelected(t *testing.T) {
	profileTools, err := ExpandToolProfile(config.ToolProfile{Tools: []string{"GetSource", "SAP"}})
	if err != nil {
		t.Fatalf("ExpandToolProfile: %v", err)
	}
	server := NewServer(&Config{
		BaseURL:      "https://example.invalid",
		Mode:         "hyperfocused",
		ProfileName:  "router",
		ProfileTools: profileTools,
	})
	got := toolNameSet(server.RegisteredTools())
	if !got["SAP"] || !got["GetSource"] {
		t.Fatalf("profile-selected tools missing: %v", server.RegisteredTools())
	}
	if got["CreateTable"] {
		t.Fatal("SAP opt-in exposed unrelated dedicated tools")
	}
}

func TestDisabledGroupsCanNarrowAProfile(t *testing.T) {
	profileTools, err := ExpandToolProfile(config.ToolProfile{Tools: []string{"ListTransports", "GetTransport", "CreateTransport"}})
	if err != nil {
		t.Fatalf("ExpandToolProfile: %v", err)
	}
	server := NewServer(&Config{
		BaseURL:        "https://example.invalid",
		Mode:           "expert",
		DisabledGroups: "C",
		ProfileName:    "transport",
		ProfileTools:   profileTools,
	})
	got := toolNameSet(server.RegisteredTools())
	for _, name := range []string{"ListTransports", "GetTransport", "CreateTransport"} {
		if got[name] {
			t.Errorf("disabled CTS tool %s remains registered", name)
		}
	}
	if !got["GetFeatures"] {
		t.Fatal("utility baseline should remain exposed")
	}
}

func TestProfileToolCatalogMatchesExpertRegistry(t *testing.T) {
	server := NewServer(&Config{BaseURL: "https://example.invalid", Mode: "expert"})
	registered := toolNameSet(server.RegisteredTools())
	known := knownProfileTools()
	for name := range registered {
		if !known[name] {
			t.Errorf("registered tool %q is missing from the profile catalog", name)
		}
	}
	for name := range known {
		if !registered[name] {
			t.Errorf("profile catalog references unregistered tool %q", name)
		}
	}
}

func TestPublishedProfileExampleExpands(t *testing.T) {
	profiles, err := config.LoadToolProfilesFromPaths("", filepath.Join("..", "..", "docs", "mcp-profiles.example.toml"))
	if err != nil {
		t.Fatalf("load example profiles: %v", err)
	}
	for name, profile := range profiles.Profiles {
		if _, err := ExpandToolProfile(profile); err != nil {
			t.Errorf("example profile %q: %v", name, err)
		}
	}
	name, profile, found, err := profiles.Resolve("")
	if err != nil || !found || name != "code-scout" {
		t.Fatalf("example default = %q, found=%v err=%v; want code-scout", name, found, err)
	}
	tools, err := ExpandToolProfile(profile)
	if err != nil || !tools["GetSource"] || tools["RunUnitTests"] || tools["RunATCCheck"] {
		t.Fatalf("code-scout tools=%v err=%v; tests/ATC should be excluded", tools, err)
	}
}

func toolNameSet(names []string) map[string]bool {
	set := make(map[string]bool, len(names))
	for _, name := range names {
		set[name] = true
	}
	return set
}
