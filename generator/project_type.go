package generator

import (
	"fmt"
	"strings"
)

const (
	// ProjectTypeService generates a lightweight standalone backend service.
	ProjectTypeService = "service"
	// ProjectTypeModular generates a product backend skeleton with module groups.
	ProjectTypeModular = "modular"
)

func normalizeProjectType(value string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", ProjectTypeService, "standalone", "app":
		return ProjectTypeService, nil
	case ProjectTypeModular, "product":
		return ProjectTypeModular, nil
	default:
		return "", fmt.Errorf("unsupported init type %q; supported types are service and modular", value)
	}
}
