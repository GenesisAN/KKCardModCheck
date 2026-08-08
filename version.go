package main

import (
	"KKCardModCheck/config"
	"fmt"
	"strings"
	"sync"

	"github.com/tcnksm/go-latest"
)

var (
	version   = "dev"
	buildDate = "unknown"
	gitHash   = "unknown"
	mod       = "debug"

	versionMu       sync.RWMutex
	version_res     *latest.CheckResponse
	version_err     error
	versionChecking bool
)

var githubTag = &latest.GithubTag{Owner: "GenesisAN", Repository: "KKCardModCheck"}

func init() {
	fmt.Println("Version:", version)
	fmt.Println("Build Date:", buildDate)
	fmt.Println("Git Hash:", gitHash)
	fmt.Println("Module Mode:", mod)
	config.Load()
	if !isDevelopmentVersion(version) {
		startVersionCheck()
	}
}

func isDevelopmentVersion(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "dev", "development", "unknown":
		return true
	default:
		return false
	}
}

func startVersionCheck() {
	versionMu.Lock()
	if isDevelopmentVersion(version) {
		version_res = nil
		version_err = nil
		versionChecking = false
		versionMu.Unlock()
		return
	}
	if versionChecking {
		versionMu.Unlock()
		return
	}
	versionChecking = true
	version_err = nil
	versionMu.Unlock()

	go func() {
		res, err := latest.Check(githubTag, version)
		versionMu.Lock()
		version_res = res
		version_err = err
		versionChecking = false
		versionMu.Unlock()
	}()
}

func versionSnapshot() (latestVersion string, outdated, checking, failed bool) {
	versionMu.RLock()
	defer versionMu.RUnlock()
	checking = versionChecking
	failed = version_err != nil
	if version_res != nil {
		latestVersion = version_res.Current
		outdated = version_res.Outdated
	}
	return
}
