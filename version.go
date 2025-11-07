package main

import (
	"KKCardModCheck/config"
	"fmt"

	"github.com/tcnksm/go-latest"
)

// Additional build metadata injected via -ldflags. Defaults provided for local/dev builds.
var (
	version     = "0.1.10"
	buildDate   = "unknown"
	gitHash     = "unknown"
	mod         = "debug"
	version_res *latest.CheckResponse
	version_err error
)

var githubTag = &latest.GithubTag{
	Owner:      "GenesisAN",
	Repository: "KKCardModCheck",
}

func init() {
	// 如果通过 ldflags 注入了 version，则把它赋给 APP_VERSION 供界面使用
	fmt.Println("Version:", version)
	fmt.Println("Build Date:", buildDate)
	fmt.Println("Git Hash:", gitHash)
	fmt.Println("Module Mode:", mod)
	// reference build metadata variables so they are not reported as unused
	_ = buildDate
	_ = gitHash
	_ = mod

	// 修复：避免使用 := 导致局部变量遮蔽包级别的 version_res
	go func() {
		res, err := latest.Check(githubTag, version)
		version_res = res
		version_err = err
	}()

	config.Load()
}
