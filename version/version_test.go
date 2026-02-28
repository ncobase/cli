package version

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"
)

func withBuildMetadata(version, branch, revision, builtAt string, fn func()) {
	previousVersion := Version
	previousBranch := Branch
	previousRevision := Revision
	previousBuiltAt := BuiltAt
	defer func() {
		Version = previousVersion
		Branch = previousBranch
		Revision = previousRevision
		BuiltAt = previousBuiltAt
	}()

	Version = version
	Branch = branch
	Revision = revision
	BuiltAt = builtAt
	fn()
}

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()

	originalStdout := os.Stdout
	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create stdout pipe: %v", err)
	}

	os.Stdout = writer
	fn()
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close stdout writer: %v", err)
	}
	os.Stdout = originalStdout

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, reader); err != nil {
		t.Fatalf("failed to copy stdout content: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("failed to close stdout reader: %v", err)
	}
	return buf.String()
}

func TestInfoString(t *testing.T) {
	info := Info{
		Version:   "v1.2.3",
		Branch:    "main",
		Revision:  "abc1234",
		BuiltAt:   "2026-01-01T00:00:00Z",
		GoVersion: "go1.23.0",
	}

	output := info.String()
	expected := []string{
		"Version: v1.2.3",
		"Branch: main",
		"Revision: abc1234",
		"Built At: 2026-01-01T00:00:00Z",
		"Go Version: go1.23.0",
	}

	for _, item := range expected {
		if !strings.Contains(output, item) {
			t.Fatalf("expected output to contain %q, got %q", item, output)
		}
	}
}

func TestInfoJSON(t *testing.T) {
	info := Info{
		Version:   "v1.2.3",
		Branch:    "main",
		Revision:  "abc1234",
		BuiltAt:   "2026-01-01T00:00:00Z",
		GoVersion: "go1.23.0",
	}

	raw := info.JSON()

	var parsed Info
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		t.Fatalf("failed to parse json output: %v; output=%q", err, raw)
	}

	if parsed != info {
		t.Fatalf("json round-trip mismatch, expected %#v, got %#v", info, parsed)
	}
}

func TestGetVersionInfo_UsesBuildMetadata(t *testing.T) {
	withBuildMetadata("v9.9.9", "release", "123abcd", "2026-02-28T00:00:00Z", func() {
		got := GetVersionInfo()

		if got.Version != "v9.9.9" {
			t.Fatalf("unexpected version: %q", got.Version)
		}
		if got.Branch != "release" {
			t.Fatalf("unexpected branch: %q", got.Branch)
		}
		if got.Revision != "123abcd" {
			t.Fatalf("unexpected revision: %q", got.Revision)
		}
		if got.BuiltAt != "2026-02-28T00:00:00Z" {
			t.Fatalf("unexpected build time: %q", got.BuiltAt)
		}
		if got.GoVersion != runtime.Version() {
			t.Fatalf("unexpected go version: %q", got.GoVersion)
		}
	})
}

func TestGetVersionInfo_FallbackBuildTime(t *testing.T) {
	withBuildMetadata("v9.9.9", "release", "123abcd", "unknown", func() {
		start := time.Now().Add(-time.Second)
		got := GetVersionInfo()
		end := time.Now().Add(time.Second)

		parsed, err := time.Parse(time.RFC3339, got.BuiltAt)
		if err != nil {
			t.Fatalf("expected RFC3339 build time, got %q: %v", got.BuiltAt, err)
		}
		if parsed.Before(start) || parsed.After(end) {
			t.Fatalf("build time not in expected range: %s", got.BuiltAt)
		}
	})
}

func TestPrint(t *testing.T) {
	withBuildMetadata("v1.0.0", "main", "abc1234", "2026-02-28T00:00:00Z", func() {
		output := captureOutput(t, func() {
			Print()
		})

		if !strings.Contains(output, "Version: v1.0.0") {
			t.Fatalf("expected printed output to include version, got %q", output)
		}
		if !strings.Contains(output, "Branch: main") {
			t.Fatalf("expected printed output to include branch, got %q", output)
		}
	})
}

func TestFlags_NoVersion(t *testing.T) {
	previous := showVersion
	defer func() {
		showVersion = previous
	}()

	showVersion = false
	Flags()
}

func TestFlags_ShowVersionExits(t *testing.T) {
	const envKey = "NCO_VERSION_FLAG_EXIT"

	if os.Getenv(envKey) == "1" {
		withBuildMetadata("v1.0.0", "main", "abc1234", "2026-02-28T00:00:00Z", func() {
			showVersion = true
			Flags()
		})
		t.Fatal("expected Flags to terminate process with os.Exit(0)")
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestFlags_ShowVersionExits$")
	cmd.Env = append(os.Environ(), envKey+"=1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected child process to exit with code 0, got error: %v; output=%s", err, string(output))
	}
	if !strings.Contains(string(output), "Version: v1.0.0") {
		t.Fatalf("expected child process output to include version information, got %q", string(output))
	}
}

func TestGitChecksWhenCommandUnavailable(t *testing.T) {
	t.Setenv("PATH", "")

	if isGitAvailable() {
		t.Fatal("expected git to be unavailable when PATH is empty")
	}
	if isGitRepository() {
		t.Fatal("expected repository check to be false when git is unavailable")
	}
}

func TestGetRuntimeGitInfo_FallbackOnCommandFailure(t *testing.T) {
	withBuildMetadata("0.0.0", "unknown", "unknown", "unknown", func() {
		t.Setenv("PATH", "")

		branch, revision, ver := getRuntimeGitInfo()
		if branch != "unknown" {
			t.Fatalf("expected fallback branch value, got %q", branch)
		}
		if revision != "unknown" {
			t.Fatalf("expected fallback revision value, got %q", revision)
		}
		if ver != "0.0.0" {
			t.Fatalf("expected fallback version value, got %q", ver)
		}
	})
}

func TestGetVersionInfo_DefaultMetadataWithoutGit(t *testing.T) {
	withBuildMetadata("0.0.0", "unknown", "unknown", "unknown", func() {
		t.Setenv("PATH", "")

		got := GetVersionInfo()
		if got.Version != "0.0.0" {
			t.Fatalf("unexpected version value: %q", got.Version)
		}
		if got.Branch != "unknown" {
			t.Fatalf("unexpected branch value: %q", got.Branch)
		}
		if got.Revision != "unknown" {
			t.Fatalf("unexpected revision value: %q", got.Revision)
		}
		if _, err := time.Parse(time.RFC3339, got.BuiltAt); err != nil {
			t.Fatalf("expected RFC3339 build time, got %q: %v", got.BuiltAt, err)
		}
	})
}
