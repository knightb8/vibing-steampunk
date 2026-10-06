# MCP Tool Profiles: Quick Guide

Profiles control which MCP tools vsp advertises to one server process. Edit the
TOML file, then restart the MCP server for its advertised tool list to update.
Profiles only control visibility; existing read-only, operation, package, and
transport safety settings still apply.

## Create or update a profile

The example file contains starter profiles:

```sh
mkdir -p ~/.vsp
cp docs/mcp-profiles.example.toml ~/.vsp/profiles.toml
```

Use `~/.vsp/profiles.toml` for global profiles or `.vsp/profiles.toml` in the
project directory for project profiles. Both files use the same TOML format.
Edit either file to add a profile or change its `groups`, `tools`, `exclude`,
or description.

```toml
default_profile = "code-scout"

[profiles.transport-expert]
description = "Inspect and manage classic CTS requests."
groups = ["transport_read", "transport_manage"]

[profiles.transport-auditor]
description = "Read transports without the default feature probe."
groups = ["transport_read"]
exclude = ["GetFeatures"]
```

## Select a profile

Selection priority is:

1. `--profile NAME` on the vsp MCP server command.
2. `default_profile` in the project file.
3. `default_profile` in the global file.
4. The existing `--mode` behavior if no profile is selected.

For example, select a different profile for one server process with
`vsp --profile transport-expert`. A project default overrides the global
default. Project profiles are added to the global profile catalog; a
same-named project profile replaces the global definition as a whole.

When a profile is active, it replaces mode-based exposure and the legacy
`.vsp.json` `tools` map. `--disabled-groups` can still narrow its tool set.
Without an active profile, existing modes and `.vsp.json` tool visibility work
as before.

## Compose the exposed tool set

- `groups` and exact tool names in `tools` are combined into a set.
- `exclude` removes exact tool names after groups and tools are expanded.
- Groups may overlap; duplicate tools appear only once.
- Unknown profile fields, group names, tool names, or exclusions stop startup.
- `GetConnectionInfo`, `GetFeatures`, and `GetAbapHelp` are included by default;
  use `exclude` to omit one. The broad `SAP` router is opt-in, e.g.
  `tools = ["SAP"]`.

Built-in group names:

| Area | Groups |
|---|---|
| Source and code | `source_read`, `source_manage`, `source_transform`, `search`, `code_intel`, `static_analysis`, `testing` |
| DDIC and data | `ddic_read`, `ddic_manage`, `data_read` |
| Transports | `transport_read`, `transport_manage` |
| Documentation and language | `documentation_read`, `documentation_manage`, `i18n_read`, `i18n_manage` |
| System and objects | `system_read`, `system_manage`, `object_manage`, `version_history` |
| Runtime and integration | `debug`, `report_execution`, `rfc_execution`, `file_transfer`, `git`, `install`, `ui5_read`, `ui5_manage`, `iam_manage` |
| Utilities | `utility` |

For the exact members of the role groups and complete profile examples, see
[`mcp-profiles.example.toml`](mcp-profiles.example.toml). In particular,
`transport_manage` includes the irreversible `ReleaseTransport` tool.
