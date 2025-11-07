<#
.SYNOPSIS
  Cross-platform PowerShell build helper for this repository.

.DESCRIPTION
  This script centralizes build metadata emission and compilation. It provides three
  modes (use switches):
    -EmitMetadata : only write build/ldflags.txt, build/version.txt, build/git_hash.txt, build/build_date.txt
    -Build        : perform a normal build, outputting a versioned binary
    -VSCBuild     : build into a stable path for VS Code debugging (workspaceBasename.exe)

  Use -Release to enable release flags (strip symbols).
#>

[CmdletBinding()]
param(
    [Alias('m')][switch]$EmitMetadata,
    [Alias('b')][switch]$Build,
    [Alias('v')][switch]$VSCBuild,
    [Alias('r')][switch]$Release,
    [Alias('g','noGui')][switch]$DisableWindowsGui
)

function Get-GitInfo {
    $info = [ordered]@{ Version = 'dev'; GitHash = 'unknown' }
    if (Get-Command git -ErrorAction SilentlyContinue) {
        try {
            $gitHash = git rev-parse --short HEAD 2>$null
            if ($gitHash) { $info.GitHash = $gitHash.Trim() }
            $descr = git describe --tags --always 2>$null
            if ($descr) { $info.Version = $descr.Trim() }
        } catch {
            # keep defaults
        }
    }
    return $info
}

function Sanitize-VersionForFileName($v) {
    if ([string]::IsNullOrWhiteSpace($v)) { return 'dev' }
    $v2 = $v -replace '^v','' -replace '[^A-Za-z0-9._-]','-'
    if ([string]::IsNullOrWhiteSpace($v2)) { return 'dev' }
    return $v2
}

# Build a short flags string for compact prompts (e.g. -m -b -r -g)
$PSBound = $PSBoundParameters.Keys

# If the user did not explicitly pass -DisableWindowsGui, choose defaults per mode:
# - In VSCBuild mode, default DisableWindowsGui = $true (i.e., disable GUI by default)
# - In Build mode, default DisableWindowsGui = $false (i.e., keep GUI by default)
if (-not $PSBound.Contains('DisableWindowsGui')) {
    if ($VSCBuild) { $DisableWindowsGui = $true }
    elseif ($Build) { $DisableWindowsGui = $false }
}

# Build a short flags string for compact prompts (e.g. -m -b -r -g)
$shortFlags = @()
if ($EmitMetadata) { $shortFlags += '-m' }
if ($Build) { $shortFlags += '-b' }
if ($VSCBuild) { $shortFlags += '-v' }
if ($Release) { $shortFlags += '-r' }
if ($DisableWindowsGui) { $shortFlags += '-g' }
if ($shortFlags.Count -eq 0) { $shortFlags = @('-') }

Write-Host "Starting build.ps1  [flags: $($shortFlags -join ' ')]"

$repoRoot = Get-Location
$repoName = Split-Path -Leaf $repoRoot
# Reliable Windows detection for both Windows PowerShell and PowerShell Core
$isWindowsRuntime = ($env:OS -eq 'Windows_NT')

$git = Get-GitInfo
$version = $git.Version
$gitHash = $git.GitHash
$buildDate = (Get-Date).ToUniversalTime().ToString('yyyy-MM-ddTHH:mm:ssZ')

# Normalize version string for metadata/ldflags: remove leading 'v' and unsafe chars
$version = Sanitize-VersionForFileName $version

if ($Release) { $stripFlags = '-s -w' } else { $stripFlags = '' }

$modeFlag = if ($Release) { '-X main.mod=release' } else { '-X main.mod=debug' }

$windowsGuiFlag = if ($DisableWindowsGui) { '' } else { '-H=windowsgui' }
$ldflags = "$stripFlags $windowsGuiFlag -extldflags=-static -X main.version=$version -X main.buildDate=$buildDate -X main.gitHash=$gitHash $modeFlag -X main.APP_VERSION=$version"

New-Item -ItemType Directory -Force -Path build | Out-Null
Set-Content -Path build/ldflags.txt -Value $ldflags -NoNewline
Set-Content -Path build/flags.txt -Value ($shortFlags -join ' ') -NoNewline

Write-Host "Wrote metadata to build/ (version=$version git=$gitHash date=$buildDate)"
Write-Host "  short flags: $($shortFlags -join ' ')"
Write-Host "  ldflags: $ldflags"

if ($EmitMetadata) {
    exit 0
}

function Do-GoBuild($outPath) {
    Write-Host "Running go mod tidy"
    & go mod tidy
    if ($LASTEXITCODE -ne 0) { throw "go mod tidy failed" }

    $gcflags = if ($Release) { '' } else { 'all=-N -l' }

    if ([string]::IsNullOrEmpty($gcflags)) {
        Write-Host "go build -ldflags '$ldflags' -o '$outPath' ."
        & go build -ldflags $ldflags -o $outPath .
    } else {
        Write-Host "go build -gcflags '$gcflags' -ldflags '$ldflags' -o '$outPath' ."
        & go build -gcflags $gcflags -ldflags $ldflags -o $outPath .
    }
    if ($LASTEXITCODE -ne 0) { throw "go build failed" }
    Write-Host "Built: $outPath"
}

if ($Build) {
    # Build only the stable filename (e.g. ImGUIKKCardTools.exe).
    # Do NOT create an extra versioned binary (repoName-version.exe).
    # Determine runtime and stable name, then build directly to it.
    $isWindowsRuntime = ($env:OS -eq 'Windows_NT')
    $stable = if ($isWindowsRuntime) { "$repoName.exe" } else { $repoName }
    $outPath = Join-Path $repoRoot $stable
    Do-GoBuild $outPath
    Write-Host "Built stable filename: $stable"
    exit 0
}

if ($VSCBuild) {
    # Ensure metadata exists; if not, emit it first
    if (-not (Test-Path -Path build/ldflags.txt)) {
        Write-Host 'Metadata missing, emitting first...'
        # write previously computed metadata (already done above), so continue
    }
        # Use same reliable Windows detection as above
        $stable = if ($isWindowsRuntime) { "$repoName.exe" } else { $repoName }
    $outPath = Join-Path $repoRoot $stable
    Do-GoBuild $outPath
    exit 0
}

Write-Host "No action requested. Use -EmitMetadata, -Build or -VSCBuild."
exit 2
