# vsp MCP Tool Exposure

This context names the MCP tool sets that vsp exposes to clients, so configuration and discussion stay precise about tool visibility.

## Language

**Tool profile**:
A named, complete set of MCP tools exposed to a client, composed from individual tools and descriptive groups, with default utility tools unless excluded.
_Avoid_: Tool preset, mode overlay

**Tool group**:
A descriptive category that expands to a specific set of MCP tools within a profile.
_Avoid_: Disable code

**Utility tool**:
A general self-inspection or help tool available across profiles unless a profile excludes it.
