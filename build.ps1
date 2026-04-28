param(
    [string]$EnvFile = ".env",
    [string]$Output = "go_selfupdate.exe"
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

if ($envValues.ContainsKey("SELFUPDATE_REPO")) {
    $repo = $envValues["SELFUPDATE_REPO"]
}
if ($envValues.ContainsKey("APP_VERSION")) {
    $appVersion = $envValues["APP_VERSION"]
}

$ldflags = @(
    "-X", "main.version=$appVersion",
    "-X", "main.selfUpdateRepo=$repo"
)

$outputDir = Split-Path -Path $Output -Parent
if ($outputDir -and -not (Test-Path -Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir | Out-Null
}

Write-Host "Building with:"
Write-Host "  APP_VERSION=$appVersion"
Write-Host "  SELFUPDATE_REPO=$repo"
Write-Host "  OUTPUT=$Output"

go build -ldflags ($ldflags -join " ") -o $Output .
if ($LASTEXITCODE -ne 0) {
    throw "go build failed with exit code $LASTEXITCODE"
}

Write-Host "Build completed successfully."
