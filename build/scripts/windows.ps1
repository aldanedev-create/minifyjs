# Windows-specific build wrapper. Called by CI on windows runners.
# Local users normally do not need this; `python build_core.py` from
# the build\ directory works on any platform.

$ErrorActionPreference = "Stop"
Set-Location (Join-Path $PSScriptRoot "..")

$Version = if ($env:VERSION) { $env:VERSION } else {
    try { git describe --tags --always --dirty } catch { "dev" }
}

Write-Host "building MinifyJS for Windows (version $Version)"

python build_core.py --platform win_amd64 --out dist --version $Version
python scripts\verify_binary.py dist\win_amd64\minifyjs.exe