package mcp

import "github.com/mark3labs/mcp-go/mcp"

// registerDDICTools exposes DDIC form creators as distinct MCP tools.
func (s *Server) registerDDICTools(shouldRegister func(string) bool) {
	if shouldRegister("CreateDomain") {
		s.mcpServer.AddTool(mcp.NewTool("CreateDomain",
			mcp.WithDescription("Create and activate a DDIC domain. Requires the normal create/package safety gates."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Domain name (max 30 characters)")),
			mcp.WithString("description", mcp.Required(), mcp.Description("Domain description")),
			mcp.WithString("data_type", mcp.Required(), mcp.Description("ABAP type, such as CHAR, NUMC, DEC, INT4, or DATS")),
			mcp.WithNumber("length", mcp.Description("Length; required for variable-length types")),
			mcp.WithNumber("decimals", mcp.Description("Decimal places for numeric types")),
			mcp.WithNumber("output_length", mcp.Description("Output length (defaults to length)")),
			mcp.WithBoolean("lowercase", mcp.Description("Allow lowercase characters")),
			mcp.WithBoolean("signed", mcp.Description("Allow signed values")),
			mcp.WithString("conversion_exit", mcp.Description("Conversion exit name, without CONVERSION_EXIT_ prefix")),
			mcp.WithString("value_table", mcp.Description("Optional check table name")),
			mcp.WithArray("fixed_values",
				mcp.Description(`Optional fixed values: [{"low":"A","high":"","text":"Alpha"}]`),
				mcp.Items(map[string]any{
					"type": "object",
					"properties": map[string]any{
						"low":  map[string]any{"type": "string", "description": "Inclusive lower value"},
						"high": map[string]any{"type": "string", "description": "Optional upper value"},
						"text": map[string]any{"type": "string", "description": "Value description (max 60 characters)"},
					},
					"required": []string{"low"},
				}),
			),
			mcp.WithString("package", mcp.Description("Target package (defaults to $TMP)")),
			mcp.WithString("transport", mcp.Description("Transport request for a transportable package")),
			mcp.WithString("language", mcp.Description("Master language (defaults to the session language)")),
		), s.handleCreateDomain)
	}

	if shouldRegister("CreateDataElement") {
		s.mcpServer.AddTool(mcp.NewTool("CreateDataElement",
			mcp.WithDescription("Create and activate a DDIC data element using either a domain or a predefined type. Requires the normal create/package safety gates."),
			mcp.WithString("name", mcp.Required(), mcp.Description("Data element name (max 30 characters)")),
			mcp.WithString("description", mcp.Required(), mcp.Description("Data element description")),
			mcp.WithString("domain", mcp.Description("Existing domain name; mutually exclusive with data_type")),
			mcp.WithString("data_type", mcp.Description("Predefined ABAP type; mutually exclusive with domain")),
			mcp.WithNumber("length", mcp.Description("Length for a predefined type")),
			mcp.WithNumber("decimals", mcp.Description("Decimal places for numeric predefined types")),
			mcp.WithString("short_label", mcp.Description("Short field label (max 10 characters)")),
			mcp.WithString("medium_label", mcp.Description("Medium field label (max 20 characters)")),
			mcp.WithString("long_label", mcp.Description("Long field label (max 40 characters)")),
			mcp.WithString("heading", mcp.Description("Heading label (max 55 characters)")),
			mcp.WithString("search_help", mcp.Description("Optional search help name")),
			mcp.WithString("parameter_id", mcp.Description("Optional SET/GET parameter ID")),
			mcp.WithBoolean("change_document", mcp.Description("Enable change document recording")),
			mcp.WithString("package", mcp.Description("Target package (defaults to $TMP)")),
			mcp.WithString("transport", mcp.Description("Transport request for a transportable package")),
			mcp.WithString("language", mcp.Description("Master language (defaults to the session language)")),
		), s.handleCreateDataElement)
	}

	if shouldRegister("CreateStructure") {
		s.mcpServer.AddTool(mcp.NewTool("CreateStructure",
			mcp.WithDescription("Create and activate a DDIC structure or append structure from DDL source. The structure name can be inferred from the source. Requires the normal create/package safety gates."),
			mcp.WithString("name", mcp.Description("Structure name; optional when declared in source")),
			mcp.WithString("description", mcp.Required(), mcp.Description("Structure description")),
			mcp.WithString("source", mcp.Required(), mcp.Description("DDL source: define structure <name> { ... } or extend type <base> with <append> { ... }")),
			mcp.WithString("package", mcp.Description("Target package (defaults to $TMP)")),
			mcp.WithString("transport", mcp.Description("Transport request for a transportable package")),
			mcp.WithString("language", mcp.Description("Master language (defaults to the session language)")),
		), s.handleCreateStructure)
	}
}
