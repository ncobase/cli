package utils

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"text/template"
)

// EnsureDir creates directory if not exists
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// WriteTemplateFile writes content to file
func WriteTemplateFile(path, content string, data any) error {
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}

	// Parse template
	tmpl, err := template.New("file").Parse(content)
	if err != nil {
		return err
	}

	// Execute template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0644)
}

// GetPlatformExt returns platform extension
func GetPlatformExt() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

// ValidateName validates Go package name
func ValidateName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z][a-zA-Z0-9_-]*$`, name)
	return matched
}

// FileExists checks file existence
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// PathExists checks path existence
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
