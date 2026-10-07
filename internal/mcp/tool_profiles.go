package mcp

import (
	"fmt"

	"github.com/oisee/vibing-steampunk/pkg/config"
)

var profileUtilityTools = []string{"GetConnectionInfo", "GetFeatures", "GetAbapHelp"}

func isProfileUtilityTool(name string) bool {
	for _, utility := range profileUtilityTools {
		if name == utility {
			return true
		}
	}
	return false
}

// profileToolGroups is the descriptive, profile-facing catalog. It is
// intentionally separate from toolGroups, whose short codes are the legacy
// --disabled-groups interface. Groups may overlap; profile expansion unions
// their members.
func profileToolGroups() map[string][]string {
	return map[string][]string{
		"source_read": {
			"GetSource", "GetProgram", "GetClass", "GetInterface", "GetFunction", "GetFunctionGroup", "GetInclude", "GetClassInclude", "GetClassComponents",
		},
		"source_manage": {
			"WriteSource", "UpdateSource", "EditSource", "WriteProgram", "WriteClass", "UpdateClassInclude",
			"CreateAndActivateProgram", "CreateClassWithTests", "CreateTestInclude",
		},
		"source_transform": {
			"PrettyPrint",
		},
		"search": {
			"SearchObject", "GrepObject", "GrepObjects", "GrepPackage", "GrepPackages",
		},
		"code_intel": {
			"FindDefinition", "FindReferences", "CodeCompletion", "GetTypeHierarchy", "GetContext", "GetClassInfo", "PrettyPrint",
			"GetCallGraph", "GetCallersOf", "GetCalleesOf", "GetObjectStructure", "AnalyzeCallGraph", "CompareCallGraphs", "GraphStats",
			"GetCDSDependencies", "GetCDSImpactAnalysis", "GetCDSElementInfo", "CompareSource", "CheckBoundaries",
		},
		"static_analysis": {
			"AnalyzeABAPCode", "SyntaxCheck", "CheckBoundaries", "AnalyzeCallGraph", "GraphStats",
		},
		"testing": {
			"RunUnitTests", "RunATCCheck", "GetATCCustomizing", "GetCodeCoverage", "GetCheckRunResults",
		},
		"ddic_read": {
			"GetTable", "GetStructure", "GetTypeInfo", "GetCDSDependencies", "GetCDSImpactAnalysis", "GetCDSElementInfo",
		},
		"ddic_manage": {
			"CreateTable", "CreateDomain", "CreateDataElement", "CreateStructure",
		},
		"data_read": {
			"GetTableContents", "RunQuery",
		},
		"transport_read": {
			"ListTransports", "GetTransport", "GetTransportInfo", "GetUserTransports",
		},
		"transport_manage": {
			"CreateTransport", "ReleaseTransport", "DeleteTransport",
		},
		"documentation_read": {
			"GetMessages", "GetTextElements", "GetTextPool", "GetObjectTextsInLanguage", "GetDataElementLabels",
			"GetMessageClassTexts", "CompareLanguages", "GetAbapHelp",
		},
		"documentation_manage": {
			"SetTextElements", "WriteMessageClassTexts",
		},
		"system_read": {
			"GetSystemInfo", "GetInstalledComponents", "GetTransaction", "GetPackage", "GetFunctionGroup", "GetAPIReleaseState",
			"GetInactiveObjects", "GetPrettyPrinterSettings", "GetFeatures", "GetConnectionInfo",
		},
		"system_manage": {
			"SetPrettyPrinterSettings",
		},
		"object_manage": {
			"LockObject", "UnlockObject", "CreateObject", "CreatePackage", "DeleteObject", "CloneObject", "RenameObject", "MoveObject",
			"RecoverFailedCreate", "Activate", "ActivateMultiple", "ActivatePackage", "PublishServiceBinding", "UnpublishServiceBinding",
		},
		"debug": {
			"SetBreakpoint", "GetBreakpoints", "DeleteBreakpoint", "DebuggerListen", "DebuggerAttach", "DebuggerDetach", "DebuggerStep",
			"DebuggerGetStack", "DebuggerGetVariables", "AMDPDebuggerStart", "AMDPDebuggerResume", "AMDPDebuggerStop", "AMDPDebuggerStep",
			"AMDPGetVariables", "AMDPSetBreakpoint", "AMDPGetBreakpoints", "ListDumps", "GetDump", "ListTraces", "GetTrace",
			"GetSQLTraceState", "ListSQLTraces", "TraceExecution",
		},
		"report_execution": {
			"RunReport", "RunReportAsync", "GetAsyncResult", "GetVariants",
		},
		"rfc_execution": {
			"CallRFC", "ExecuteABAP",
		},
		"file_transfer": {
			"ImportFromFile", "ExportToFile", "DeployFromFile", "SaveToFile", "DeployZip",
		},
		"git": {
			"GitTypes", "GitExport",
		},
		"install": {
			"InstallZADTVSP", "ListDependencies",
		},
		"ui5_read": {
			"UI5ListApps", "UI5GetApp", "UI5GetFileContent",
		},
		"ui5_manage": {
			"UI5CreateApp", "UI5DeleteApp", "UI5DeleteFile", "UI5UploadFile",
		},
		"version_history": {
			"GetRevisions", "GetRevisionSource", "CompareVersions",
		},
		"iam_manage": {
			"CreateBusinessCatalog", "CreateIAMApp", "AssignIAMAppToCatalog",
		},
		"i18n_read": {
			"GetObjectTextsInLanguage", "GetDataElementLabels", "GetMessageClassTexts", "GetTextPool", "CompareLanguages",
		},
		"i18n_manage": {
			"WriteMessageClassTexts",
		},
		"utility": profileUtilityTools,
	}
}

// ExpandToolProfile validates profile references and returns the complete
// selected tool set. The small utility baseline is implicit; SAP is not.
func ExpandToolProfile(profile config.ToolProfile) (map[string]bool, error) {
	groups := profileToolGroups()
	known := make(map[string]bool)
	for _, members := range groups {
		for _, name := range members {
			known[name] = true
		}
	}
	known["SAP"] = true
	for _, name := range profileUtilityTools {
		known[name] = true
	}

	selected := make(map[string]bool)
	for _, name := range profileUtilityTools {
		selected[name] = true
	}
	for _, name := range profile.Tools {
		if !known[name] {
			return nil, fmt.Errorf("unknown MCP tool %q in profile", name)
		}
		selected[name] = true
	}
	for _, group := range profile.Groups {
		members, ok := groups[group]
		if !ok {
			return nil, fmt.Errorf("unknown MCP tool group %q in profile", group)
		}
		for _, name := range members {
			selected[name] = true
		}
	}
	for _, name := range profile.Exclude {
		if !known[name] {
			return nil, fmt.Errorf("unknown MCP tool %q in profile exclude list", name)
		}
		delete(selected, name)
	}
	return selected, nil
}

func knownProfileTools() map[string]bool {
	known := make(map[string]bool)
	for _, members := range profileToolGroups() {
		for _, name := range members {
			known[name] = true
		}
	}
	known["SAP"] = true
	for _, name := range profileUtilityTools {
		known[name] = true
	}
	return known
}
