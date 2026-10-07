package main

import (
	"testing"

	"github.com/oisee/vibing-steampunk/pkg/config"
)

func TestResolveMCPToolProfileSelection(t *testing.T) {
	profiles := &config.ToolProfiles{
		Profiles: map[string]config.ToolProfile{
			"global":  {Groups: []string{"source_read"}},
			"project": {Groups: []string{"transport_read"}},
			"server":  {Tools: []string{"GetTable"}},
		},
		GlobalDefault:  "global",
		ProjectDefault: "project",
	}
	for _, tc := range []struct {
		name     string
		explicit string
		wantName string
		wantTool string
	}{
		{name: "explicit server profile wins", explicit: "server", wantName: "server", wantTool: "GetTable"},
		{name: "project default wins global", wantName: "project", wantTool: "ListTransports"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name, tools, err := resolveMCPToolProfileFrom(tc.explicit, profiles)
			if err != nil {
				t.Fatalf("resolveMCPToolProfileFrom: %v", err)
			}
			if name != tc.wantName || !tools[tc.wantTool] {
				t.Fatalf("resolved %q with tools %v; want %q including %s", name, tools, tc.wantName, tc.wantTool)
			}
			if !tools["GetFeatures"] {
				t.Fatal("utility baseline missing from resolved profile")
			}
		})
	}
}

func TestResolveMCPToolProfileWithoutDefaultUsesLegacyMode(t *testing.T) {
	name, tools, err := resolveMCPToolProfileFrom("", &config.ToolProfiles{Profiles: map[string]config.ToolProfile{}})
	if err != nil || name != "" || tools != nil {
		t.Fatalf("resolved name=%q tools=%v err=%v; want no profile", name, tools, err)
	}
}

func TestResolveMCPToolProfileRejectsInvalidSelectedProfile(t *testing.T) {
	profiles := &config.ToolProfiles{Profiles: map[string]config.ToolProfile{
		"broken": {Groups: []string{"unknown_group"}},
	}}
	if _, _, err := resolveMCPToolProfileFrom("broken", profiles); err == nil {
		t.Fatal("resolveMCPToolProfileFrom accepted an invalid group")
	}
}

func TestResolveMCPToolProfileRejectsInvalidUnusedProfiles(t *testing.T) {
	profiles := &config.ToolProfiles{Profiles: map[string]config.ToolProfile{
		"active":   {Groups: []string{"source_read"}},
		"mistyped": {Groups: []string{"unknown_group"}},
	}, ProjectDefault: "active"}
	if _, _, err := resolveMCPToolProfileFrom("", profiles); err == nil {
		t.Fatal("resolveMCPToolProfileFrom accepted an invalid unused profile")
	}
}
