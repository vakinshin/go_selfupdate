param(
    [string]$EnvFile = ".env",
    [string]$OutputDir = "dist",
    [string]$BinaryName = "go_selfupdate",
    [string[]]$Targets = @("windows/amd64")
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Read-DotEnv {
    param([string]$Path)

    $values = @{}
    if (-not (Test-Path -Path $Path)) {
        throw "Env file not found: $Path"
    }

    foreach ($line in [System.IO.File]::ReadAllLines($Path)) {
        $trimmed = $line.Trim()
        if ($trimmed -eq "" -or $trimmed.StartsWith("#")) {
            continue
        }

        $idx = $trimmed.IndexOf("=")
        if ($idx -lt 1) {
            continue
        }

        $key = $trimmed.Substring(0, $idx).Trim()
        $value = $trimmed.Substring($idx + 1).Trim()
        $values[$key] = $value
    }

    return $values
}

$envValues = Read-DotEnv -Path $EnvFile
$repo = ""
$appVersion = "0.1.0"
$githubToken = ""

if ($envValues.ContainsKey("SELFUPDATE_REPO")) {
    $repo = $envValues["SELFUPDATE_REPO"]
}
if ($envValues.ContainsKey("APP_VERSION")) {
    $appVersion = $envValues["APP_VERSION"]
}
if ($envValues.ContainsKey("GITHUB_TOKEN")) {
    $githubToken = $envValues["GITHUB_TOKEN"]
}

if ($repo -eq "") {
    throw "SELFUPDATE_REPO is required in $EnvFile"
}

$ldflags = @(
    "-X", "main.version=$appVersion",
    "-X", "main.selfUpdateRepo=$repo",
    "-X", "main.githubToken=$githubToken"
)

if (-not (Test-Path -Path $OutputDir)) {
    New-Item -ItemType Directory -Path $OutputDir | Out-Null
}

Write-Host "Building with:"
Write-Host "  APP_VERSION=$appVersion"
Write-Host "  SELFUPDATE_REPO=$repo"
Write-Host "  GITHUB_TOKEN_SET=$($githubToken -ne '')"
Write-Host "  OUTPUT_DIR=$OutputDir"
Write-Host "  TARGETS=$($Targets -join ', ')"

$ldflagsString = $ldflags -join " "
$builtAssets = @()

foreach ($target in $Targets) {
    $parts = $target.Split("/")
    if ($parts.Length -ne 2) {
        throw "Invalid target '$target'. Use format os/arch, for example windows/amd64"
    }

    $goos = $parts[0].Trim().ToLowerInvariant()
    $goarch = $parts[1].Trim().ToLowerInvariant()
    if ($goos -eq "" -or $goarch -eq "") {
        throw "Invalid target '$target'. os and arch are required."
    }

    $assetName = "${BinaryName}_${goos}_${goarch}"
    if ($goos -eq "windows") {
        $assetName += ".exe"
    }

    $outputPath = Join-Path $OutputDir $assetName
    Write-Host ""
    Write-Host "Building asset: $assetName"

    $env:GOOS = $goos
    $env:GOARCH = $goarch
    go build -ldflags $ldflagsString -o $outputPath .
    if ($LASTEXITCODE -ne 0) {
        throw "go build failed for target '$target' with exit code $LASTEXITCODE"
    }

    $builtAssets += $outputPath
}

Remove-Item Env:GOOS -ErrorAction SilentlyContinue
Remove-Item Env:GOARCH -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "Build completed successfully."
Write-Host "Assets:"
foreach ($asset in $builtAssets) {
    Write-Host "  $asset"
}
