//go:build integration

package main

import (
	"testing"

	"github.com/tcnksm/go-latest"
)

func TestDevelopmentBuildSkipsAutomaticVersionCheck(t *testing.T) {
	if !isDevelopmentVersion(version) {
		t.Fatalf("expected source build version to be development, got %q", version)
	}
	latestVersion, outdated, checking, failed := versionSnapshot()
	if latestVersion != "" || outdated || checking || failed {
		t.Fatalf("development build should not automatically check updates: latest=%q outdated=%v checking=%v failed=%v", latestVersion, outdated, checking, failed)
	}
}

func TestReleaseVersionCheckAgainstGitHub(t *testing.T) {
	response, err := latest.Check(githubTag, "1.0.1")
	if err != nil {
		t.Fatalf("GitHub version check failed: %v", err)
	}
	if response.Current == "" {
		t.Fatal("GitHub version check returned an empty latest version")
	}
	if response.Outdated {
		t.Fatalf("release 1.0.1 was incorrectly marked outdated; latest is %s", response.Current)
	}
	t.Logf("release=1.0.1 latest=%s outdated=%v", response.Current, response.Outdated)
}
