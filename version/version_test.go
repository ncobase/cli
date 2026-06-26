package version

import (
	"runtime"
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
