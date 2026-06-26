package generator

import (
	"fmt"
	"strings"
)

// DatabasePlan describes the primary data access selection for generated code.
type DatabasePlan struct {
	ORM    string `json:"orm"`
	Driver string `json:"driver"`
}

// IntegrationPlan describes optional ncore integration modules.
type IntegrationPlan struct {
	Cache     []string `json:"cache,omitempty"`
	Search    []string `json:"search,omitempty"`
	Messaging []string `json:"messaging,omitempty"`
	Storage   []string `json:"storage,omitempty"`
	Services  []string `json:"services,omitempty"`
}

// ModuleRequirement describes a Go module requirement in the generated go.mod.
type ModuleRequirement struct {
	Module  string `json:"module"`
	Version string `json:"version"`
}

// Operation describes commands required after templates are written.
type Operation struct {
	Name       string   `json:"name"`
	Command    string   `json:"command"`
	Args       []string `json:"args"`
	WorkingDir string   `json:"working_dir"`
	Outputs    []string `json:"outputs,omitempty"`
}

// Plan describes the exact generation target and the files that will be written.
type Plan struct {
	Name               string              `json:"name"`
	Type               string              `json:"type"`
	CustomDir          string              `json:"custom_dir,omitempty"`
	OutputPath         string              `json:"output_path"`
	BasePath           string              `json:"base_path"`
	ModuleName         string              `json:"module_name"`
	PackagePath        string              `json:"package_path"`
	Description        string              `json:"description"`
	Standalone         bool                `json:"standalone"`
	WithCmd            bool                `json:"with_cmd"`
	WithTest           bool                `json:"with_test"`
	WithGRPC           bool                `json:"with_grpc"`
	WithTracing        bool                `json:"with_tracing"`
	Database           DatabasePlan        `json:"database"`
	Integrations       IntegrationPlan     `json:"integrations"`
	Directories        []string            `json:"directories"`
	Files              []string            `json:"files"`
	ModuleRequirements []ModuleRequirement `json:"module_requirements,omitempty"`
	Operations         []Operation         `json:"operations,omitempty"`
	Conflicts          []string            `json:"conflicts,omitempty"`
}

// Result is returned for both dry-run and applied generation.
type Result struct {
	DryRun  bool   `json:"dry_run"`
	Applied bool   `json:"applied"`
	Message string `json:"message"`
	Plan    *Plan  `json:"plan"`
}

// Text returns a concise human-readable result summary.
func (r *Result) Text() string {
	if r == nil || r.Plan == nil {
		return ""
	}

	var b strings.Builder
	if r.Message != "" {
		b.WriteString(r.Message)
		b.WriteString("\n")
	}

	fmt.Fprintf(&b, "Target: %s\n", r.Plan.BasePath)
	fmt.Fprintf(&b, "Module: %s\n", r.Plan.ModuleName)
	fmt.Fprintf(&b, "Package: %s\n", r.Plan.PackagePath)
	fmt.Fprintf(&b, "Database: %s/%s\n", r.Plan.Database.ORM, r.Plan.Database.Driver)
	fmt.Fprintf(&b, "Directories: %d\n", len(r.Plan.Directories))
	fmt.Fprintf(&b, "Files: %d\n", len(r.Plan.Files))

	if len(r.Plan.Operations) > 0 {
		fmt.Fprintf(&b, "Operations: %d\n", len(r.Plan.Operations))
	}
	if len(r.Plan.Conflicts) > 0 {
		b.WriteString("Conflicts:\n")
		for _, conflict := range r.Plan.Conflicts {
			fmt.Fprintf(&b, "  - %s\n", conflict)
		}
	}

	return strings.TrimRight(b.String(), "\n")
}
