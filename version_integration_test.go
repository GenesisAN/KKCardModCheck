//go:build integration

package main

import (
	"testing"
	"time"
)

func TestVersionCheckAgainstGitHub(t *testing.T) {
	deadline := time.Now().Add(20 * time.Second)
	for {
		latestVersion, outdated, checking, failed := versionSnapshot()
		if !checking {
			if failed {
				t.Fatal("GitHub version check failed")
			}
			if latestVersion == "" {
				t.Fatal("GitHub version check returned an empty latest version")
			}
			if outdated {
				t.Fatalf("local version %s was incorrectly marked outdated; latest is %s", version, latestVersion)
			}
			t.Logf("local=%s latest=%s outdated=%v", version, latestVersion, outdated)
			return
		}
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for GitHub version check")
		}
		time.Sleep(100 * time.Millisecond)
	}
}
