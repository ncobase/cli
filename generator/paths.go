package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ncobase/cli/utils"
)

// getPackagePath returns the package path based on options
func getPackagePath(opts *Options) string {
	if opts.Standalone {
		return opts.ModuleName
	}
	switch opts.Type {
	case "custom":
		if opts.CustomDir == "" {
			return fmt.Sprintf("%s/%s", opts.ModuleName, opts.Name)
		}
		return fmt.Sprintf("%s/%s/%s", opts.ModuleName, opts.CustomDir, opts.Name)
	case "direct":
		return fmt.Sprintf("%s/%s", opts.ModuleName, opts.Name)
	default:
		return fmt.Sprintf("%s/%s/%s", opts.ModuleName, opts.Type, opts.Name)
	}
}

// resolveOutputPath determines the output path
func resolveOutputPath(opts *Options) (string, error) {
	if opts.OutputPath == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %v", err)
		}
		return cwd, nil
	}
	return opts.OutputPath, nil
}

// resolveModuleName determines the module name
func resolveModuleName(opts *Options, outputPath string) string {
	if opts.ModuleName != "" {
		return opts.ModuleName
	}

	if opts.Standalone {
		return opts.Name
	}

	// Try to detect from go.mod
	goModPath := filepath.Join(outputPath, "go.mod")
	if utils.FileExists(goModPath) {
		if content, err := os.ReadFile(goModPath); err == nil {
			for _, line := range strings.Split(string(content), "\n") {
				if strings.HasPrefix(line, "module ") {
					return strings.TrimSpace(strings.TrimPrefix(line, "module "))
				}
			}
		}
	}

	// Use directory name as fallback
	dirs := strings.Split(outputPath, string(os.PathSeparator))
	return dirs[len(dirs)-1]
}

// getBasePath returns the base path for generation
func getBasePath(opts *Options, outputPath string) string {
	if opts.Standalone {
		if opts.Type == "direct" {
			return filepath.Join(outputPath, opts.Name)
		} else if opts.Type == "custom" {
			return filepath.Join(outputPath, opts.CustomDir, opts.Name)
		}
		return filepath.Join(outputPath, opts.Type, opts.Name)
	}

	switch opts.Type {
	case "core":
		return filepath.Join(outputPath, "core", opts.Name)
	case "business":
		return filepath.Join(outputPath, "business", opts.Name)
	case "plugin":
		return filepath.Join(outputPath, "plugin", opts.Name)
	case "direct":
		return filepath.Join(outputPath, opts.Name)
	case "custom":
		return filepath.Join(outputPath, opts.CustomDir, opts.Name)
	default:
		return filepath.Join(outputPath, opts.Name)
	}
}
