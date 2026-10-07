# ADR-007: Named MCP Tool Profiles

**Date:** 2026-10-05
**Status:** PROPOSED

## Summary

Add named MCP tool profiles so each server can advertise a focused, role-specific tool set. Profiles live in dedicated TOML files at `~/.vsp/profiles.toml` and `.vsp/profiles.toml`; project definitions augment global definitions and replace a same-named global profile as a whole. This keeps profile discovery independent of `.vsp.json`, whose current first-found behavior does not merge project and home configuration.

## Decision

1. A selected profile replaces mode-based tool exposure and the legacy `.vsp.json` `tools` map. Without a selected profile, existing mode and tool-map behavior remains. Profile selection resolves in this order: `--profile`, project `default_profile`, global `default_profile`, then the existing mode fallback. A configured but unresolved profile is a startup error.
2. A profile composes exact tool names and built-in descriptive groups. Tool and group expansions are unioned, overlapping groups are deduplicated, and exact tool exclusions are applied last. Unknown tool or group names fail startup rather than being ignored.
3. The built-in catalog covers registered tool domains, with read/manage groups split where useful. Groups are descriptive and separate from legacy `--disabled-groups` codes. gCTS is excluded initially because its registration function is not currently called; the multi-character `GC` legacy-code parsing issue remains a separate follow-up.
4. Every tool is filtered by the selected profile except the utility baseline: `GetConnectionInfo`, `GetFeatures`, and `GetAbapHelp` are included unless explicitly excluded. The universal `SAP` action router is opt-in and is not part of that baseline.
5. Profiles control MCP visibility only. Existing read-only, operation, transport, and package safety checks remain authoritative.
6. The initial `transport_read` group contains `ListTransports` and `GetTransport`; `transport_manage` contains `CreateTransport`, `ReleaseTransport`, and `DeleteTransport`. `ReleaseTransport` remains in the manage group. `code-scout` excludes test/ATC execution; `code-documenter` is read-only.
7. `docs/mcp-profiles.example.toml` demonstrates `transport-expert`, `code-scout`, and `code-documenter`. The global example default is `code-scout`; project defaults may override it. The transport example combines both CTS groups.

## Consequences

- A project can reuse global profiles and selectively override one by name without copying the entire global profile catalog.
- The tool list stays small for specialist agents, while safety controls remain configured independently.
- Adding or renaming MCP tools requires keeping the built-in group catalog and its validation in sync with actual registration.
