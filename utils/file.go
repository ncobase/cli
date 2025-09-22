package utils

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

// EnsureDir creates directory if not exists
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// WriteTemplateFile writes content to file
func WriteTemplateFile(path, content string, data interface{}) error {
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
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
